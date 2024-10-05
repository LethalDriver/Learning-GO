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
    send chan structs.Message
    room *ChatRoom
}

func HandleConnection(ws *websocket.Conn, s *RoomService, roomId string, user structs.UserDetails) error {
    ctx := context.Background()

    log.Println("Handling connection")
    conn := &Connection{
        ws:   ws,
        user: user,
        send: make(chan structs.Message, 256),
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
		var msg structs.Message
		err = json.Unmarshal(messageBytes, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			break
		}
		msg.SentBy = c.user
        c.room.Broadcast <- msg
    }
    log.Println("Exiting readPump")
}

func (c *Connection) writePump() {
    defer func() {
        log.Println("Sending close message and closing WebSocket connection in writePump")
        c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
        c.ws.Close()
    }()

    for message := range c.send {
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling message:", err)
			break
		}
        if err := c.ws.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
            log.Printf("Error writing message in writePump: %v", err)
            break
        }
        log.Printf("Wrote message to connection: %q, address: %p", string(messageBytes), c)
    }
    log.Println("Exiting writePump")
}
