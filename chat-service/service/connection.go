package service

import (
	"encoding/json"
	"log"
	"sync"

	"example.com/chat_app/chat_service/structs"
	"github.com/gorilla/websocket"
)

// Connection represents a WebSocket connection to a chat room.
type Connection struct {
	ws          *websocket.Conn
	user        structs.UserDetails
	sendMessage chan structs.Message
	sendSeen    chan structs.SeenMessage
	sendDelete  chan structs.DeleteMessage
	room        *ChatRoom
}

// handleConnection handles a new WebSocket connection to a chat room.
func handleConnection(ws *websocket.Conn, room *ChatRoom, user structs.UserDetails) error {
	clientIP := ws.RemoteAddr().String()
	log.Printf("Handling connection from %s", clientIP)

	conn := &Connection{
		ws:          ws,
		user:        user,
		sendMessage: make(chan structs.Message, 256),
		sendSeen:    make(chan structs.SeenMessage, 256),
		sendDelete:  make(chan structs.DeleteMessage, 256),
		room:        room,
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

// readPump reads messages from the WebSocket connection.
func (c *Connection) readPump() {
	defer c.closeWebSocket("Closing WebSocket connection in readPump")

	for {
		_, messageBytes, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message in readPump: %v", err)
			break
		}
		log.Printf("Read message from connection: %q, address: %p", string(messageBytes), c)

		var incomingMessage structs.WsMessage
		if err := c.unmarshalMessage(messageBytes, &incomingMessage); err != nil {
			log.Printf("Error unmarshalling incoming message: %v", err)
			break
		}
		data := incomingMessage.Data
		switch incomingMessage.Type {
		case structs.TypeTextMessage:
			var msg structs.Message
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("Error unmarshalling message data: %v", err)
				break
			}
			msg.SentBy = c.user.Id
			c.room.Text <- msg
			log.Printf("Received structs.Message: %+v", msg)

		case structs.TypeSeenMessage:
			var seenMessage structs.SeenMessage
			if err := json.Unmarshal(data, &seenMessage); err != nil {
				log.Printf("Error unmarshalling seen message: %v", err)
				break
			}
			seenMessage.SeenBy = c.user
			log.Printf("Received SeenUpdate: %+v", seenMessage)
			c.room.Seen <- seenMessage

		case structs.TypeDeleteMessage:
			var deleteMessage structs.DeleteMessage
			if err := json.Unmarshal(data, &deleteMessage); err != nil {
				log.Printf("Error unmarshalling delete message: %v", err)
				break
			}
			deleteMessage.SentBy = c.user
			log.Printf("Received DeleteMessage: %+v", deleteMessage)
			c.room.Delete <- deleteMessage

		default:
			log.Printf("Unknown message type: %s", incomingMessage.Type)
		}
	}
	log.Println("Exiting readPump")
}

// writePump writes messages to the WebSocket connection.
func (c *Connection) writePump() {
	defer c.closeWebSocket("Sending close message and closing WebSocket connection in writePump")

	for {
		select {
		case message, ok := <-c.sendMessage:
			if !ok {
				log.Println("sendMessage channel closed")
				return
			}
			if err := c.writeMessage(structs.TypeTextMessage, message); err != nil {
				log.Printf("Error writing message: %v to websocket connection: %s", err, c.ws.RemoteAddr().String())
				return
			}

		case seenMessage, ok := <-c.sendSeen:
			if !ok {
				log.Println("seenUpdate channel closed")
				return
			}
			if err := c.writeMessage(structs.TypeSeenMessage, seenMessage); err != nil {
				log.Printf("Error writing seen message: %v to websocket connection: %s", err, c.ws.RemoteAddr().String())
				return
			}

		case deleteMessage, ok := <-c.sendDelete:
			if !ok {
				log.Println("deleteMessage channel closed")
				return
			}
			if err := c.writeMessage(structs.TypeDeleteMessage, deleteMessage); err != nil {
				log.Printf("Error writing delete message: %v to websocket connection: %s", err, c.ws.RemoteAddr().String())
				return
			}
		}
	}
}

// closeWebSocket closes the WebSocket connection with a log message.
func (c *Connection) closeWebSocket(logMessage string) {
	log.Println(logMessage)
	c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.ws.Close()
}

// unmarshalMessage unmarshals a byte slice into a given interface.
func (c *Connection) unmarshalMessage(messageBytes []byte, v any) error {
	err := json.Unmarshal(messageBytes, v)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
	}
	return err
}

// writeMessage marshals data and writes it to the WebSocket connection.
func (c *Connection) writeMessage(messageType structs.MessageType, data any) error {
	contentBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}
	message := structs.WsMessage{
		Type: messageType,
		Data: contentBytes,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if err := c.ws.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}
	log.Printf("Wrote message to connection: %q, address: %p", string(contentBytes), c)
	return nil
}
