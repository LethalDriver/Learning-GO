package service

import (
	"encoding/json"
	"log"
	"sync"

	"example.com/myproject/structs"
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws             *websocket.Conn
	user           structs.UserDetails
	sendMessage    chan structs.Message
	sendSeenUpdate chan structs.SeenMessage
	room           *ChatRoom
}

func handleConnection(ws *websocket.Conn, room *ChatRoom, user structs.UserDetails) error {
	clientIP := ws.RemoteAddr().String()
	log.Printf("Handling connection from %s", clientIP)

	conn := &Connection{
		ws:             ws,
		user:           user,
		sendMessage:    make(chan structs.Message, 256),
		sendSeenUpdate: make(chan structs.SeenMessage, 256),
		room:           room,
	}

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

	room.Register <- conn
	wg.Wait()
	room.Unregister <- conn
	log.Printf("Connection %s unregistered from room ID: %s", clientIP, room.Id)

	return nil
}

func (c *Connection) readPump() {
	defer c.closeWebSocket("Closing WebSocket connection in readPump")

	for {
		_, messageBytes, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message in readPump: %v", err)
			break
		}
		log.Printf("Read message from connection: %q, address: %p", string(messageBytes), c)

		messageType, err := structs.DetermineMessageType(messageBytes)
		if err != nil {
			log.Printf("Error determining message type: %v", err)
			break
		}

		switch messageType {
		case structs.TypeSeenMessage:
			var seenUpdate structs.SeenMessage
			if err := c.unmarshalMessage(messageBytes, &seenUpdate); err != nil {
				break
			}
			seenUpdate.SeenBy = c.user
			log.Printf("Received SeenUpdate: %+v", seenUpdate)
			c.room.StatusUpdates <- seenUpdate

		case structs.TypeTextMessage:
			var msg structs.Message
			if err := c.unmarshalMessage(messageBytes, &msg); err != nil {
				break
			}
			msg.SentBy = c.user
			c.room.Broadcast <- msg
			log.Printf("Received Message: %+v", msg)

		case structs.TypeDeleteMessage:
			//TODO

		default:
			log.Printf("Unknown message type: %q", string(messageBytes))
		}
	}
	log.Println("Exiting readPump")
}

func (c *Connection) writePump() {
	defer c.closeWebSocket("Sending close message and closing WebSocket connection in writePump")

	for {
		select {
		case message, ok := <-c.sendMessage:
			if !ok {
				log.Println("sendMessage channel closed")
				return
			}
			if err := c.writeMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case seenUpdate, ok := <-c.sendSeenUpdate:
			if !ok {
				log.Println("seenUpdate channel closed")
				return
			}
			if err := c.writeMessage(websocket.TextMessage, seenUpdate); err != nil {
				return
			}
		}
	}
}

func (c *Connection) closeWebSocket(logMessage string) {
	log.Println(logMessage)
	c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.ws.Close()
}

func (c *Connection) unmarshalMessage(messageBytes []byte, v interface{}) error {
	err := json.Unmarshal(messageBytes, v)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
	}
	return err
}

func (c *Connection) writeMessage(messageType int, data interface{}) error {
	messageBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}
	if err := c.ws.WriteMessage(messageType, messageBytes); err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}
	log.Printf("Wrote message to connection: %q, address: %p", string(messageBytes), c)
	return nil
}
