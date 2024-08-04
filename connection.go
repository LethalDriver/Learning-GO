package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ws *websocket.Conn 
	send chan []byte
	room *ChatRoom
}

func (c *Connection) writePump() {
    for message := range c.send {
        if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Printf("Error writing message: %v", err)
            return
        }
    }
    c.ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Connection) readPump() {
	defer func() {
		c.room.Unregister <- c
		c.ws.Close()
	}()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message %v", err)
			break
		}
		c.room.Broadcast <- message
	}

}
