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
        log.Println("Closing WebSocket connection in readPump")
        c.ws.Close() // Close the WebSocket connection
    }()

    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            log.Printf("Error reading message in readPump: %v", err)
            break
        }
        log.Printf("Read message from connection: %q, address: %p", string(message), c)
        c.room.Broadcast <- message
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
        if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Printf("Error writing message in writePump: %v", err)
            break
        }
        log.Printf("Wrote message to connection: %q, address: %p", string(message), c)
    }
    log.Println("Exiting writePump")
}

func (c *Connection) close() {
    log.Println("Closing send channel and WebSocket connection")
    close(c.send)
    c.ws.Close()
}