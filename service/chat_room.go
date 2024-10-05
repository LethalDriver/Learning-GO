package service

import (
	"context"
	"log"
	"time"

	"example.com/myproject/structs"
)

type ChatRoom struct {
	Id          string
	Members     map[*Connection]bool
	Broadcast   chan structs.Message
	Register    chan *Connection
	Unregister  chan *Connection
}

func NewChatRoom(roomId string) *ChatRoom {
	return &ChatRoom{
		Id:          roomId,
		Members:     make(map[*Connection]bool),
		Broadcast:   make(chan structs.Message),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
	}
}

func (r *ChatRoom) Run(service *RoomService) {
	for {
		select {
		case conn := <-r.Register:
			log.Printf("Registering connection to room %s, address: %p", r.Id, conn)
			r.Members[conn] = true
		case conn := <-r.Unregister:
			log.Printf("Unregistering connection from room %s", r.Id)
			if _, ok := r.Members[conn]; ok {
				delete(r.Members, conn)
				close(conn.send)
			}
		case message := <-r.Broadcast:
			log.Printf("Broadcasting message to room %s: %s", r.Id, string(message.Content))
			ctx := context.Background()
			message.SentAt = time.Now()

			id, err := service.SaveMessage(ctx, r.Id, &message)
			if err != nil {
				log.Printf("Error saving message %q in room %s", string(message.Content), r.Id)
				break
			}
			message.Id = id
			for conn := range r.Members {
				log.Printf("Sending message to write pump of connection %p", conn)
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(r.Members, conn)
				}
			}
		}
	}
}