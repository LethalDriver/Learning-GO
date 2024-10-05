package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"example.com/myproject/structs"
	"github.com/gorilla/websocket"
)

type Connection struct {
    ws   *websocket.Conn
    user structs.UserDetails
    sendMessage chan structs.Message
	sendSeenUpdate chan structs.SeenUpdate
    room *ChatRoom
}

func HandleConnection(ws *websocket.Conn, s *RoomService, roomId string, user structs.UserDetails) error {
    ctx := context.Background()

    log.Println("Handling connection")
    conn := &Connection{
        ws:   ws,
        user: user,
        sendMessage: make(chan structs.Message, 256),
		sendSeenUpdate: make(chan structs.SeenUpdate, 256),
    }

    room, err := s.GetOrCreateRoom(ctx, roomId, conn)
    if err != nil {
        return fmt.Errorf("error creating room: %v", err)
    }
    conn.room = room

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        conn.writePump()
    }()

    go func() {
        defer wg.Done()
        conn.readPump()
    }()

    wg.Wait()
    room.Unregister <- conn
    log.Printf("Connection unregistered from room ID: %s", roomId)

    return nil
}

func (c *Connection) readPump() {
    defer func() {
        log.Println("Closing WebSocket connection in readPump")
        c.ws.Close() // Close the WebSocket connection
    }()

    for {
        _, messageBytes, err := c.ws.ReadMessage()
        if err != nil {
            log.Printf("Error reading message in readPump: %v", err)
            break
        }
        log.Printf("Read message from connection: %q, address: %p", string(messageBytes), c)
		messageType, err := structs.DetermineDataType(messageBytes)
        if err != nil {
            log.Printf("Error determining message type: %v", err)
            break
        }

        switch messageType {
        case structs.StatusUpdate:
            var seenUpdate structs.SeenUpdate
            err = json.Unmarshal(messageBytes, &seenUpdate)
            if err != nil {
                log.Printf("Error unmarshalling SeenUpdate: %v", err)
                break
            }
			seenUpdate.SeenBy = c.user
			log.Printf("Received SeenUpdate: %+v", seenUpdate)
			c.room.StatusUpdates <- seenUpdate
        case structs.MessageWithContent:
            var msg structs.Message
            err = json.Unmarshal(messageBytes, &msg)
            if err != nil {
                log.Printf("Error unmarshalling Message: %v", err)
                break
            }
            msg.SentBy = c.user
            c.room.Broadcast <- msg
            log.Printf("Received Message: %+v", msg)
        default:
            log.Printf("Unknown message type: %q", string(messageBytes))
        }
    }
    log.Println("Exiting readPump")
}

func (c *Connection) writePump() {
    defer func() {
        log.Println("Sending close message and closing WebSocket connection in writePump")
        c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
        c.ws.Close()
    }()

    for {
        select {
        case message, ok := <-c.sendMessage:
            if !ok {
                log.Println("sendMessage channel closed")
                return
            }
            messageBytes, err := json.Marshal(message)
            if err != nil {
                log.Println("Error marshalling message:", err)
                return
            }
            if err := c.ws.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
                log.Printf("Error writing message in writePump: %v", err)
                return
            }
            log.Printf("Wrote message to connection: %q, address: %p", string(messageBytes), c)

        case seenUpdate, ok := <-c.sendSeenUpdate:
            if !ok {
                log.Println("seenUpdate channel closed")
                return
            }
            seenUpdateBytes, err := json.Marshal(seenUpdate)
            if err != nil {
                log.Println("Error marshalling seen update:", err)
                return
            }
            if err := c.ws.WriteMessage(websocket.TextMessage, seenUpdateBytes); err != nil {
                log.Printf("Error writing seen update in writePump: %v", err)
                return
            }
            log.Printf("Wrote seen update to connection: %q, address: %p", string(seenUpdateBytes), c)
        }
    }
}
