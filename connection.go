package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
    ws   *websocket.Conn
    send chan []byte
	room *ChatRoom
}

func (c *Connection) readPump() {
	defer func() {
		c.ws.Close() // Close the WebSocket connection
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		c.room.Broadcast <- message
	}
}

func (c *Connection) writePump() {
	defer func() {
		c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.ws.Close()
	}()

	for message := range c.send {
		if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}


func NewConnection(ws *websocket.Conn) *Connection {
    return &Connection{
        ws:   ws,
        send: make(chan []byte, 256),
    }
}
