package service

import (
	"context"
	"log"

	"example.com/myproject/structs"
)

type ChatRoom struct {
	Id            string
	Members       map[*Connection]bool
	Broadcast     chan structs.Message
	StatusUpdates chan structs.SeenUpdate
	Register      chan *Connection
	Unregister    chan *Connection
}

func NewChatRoom(roomId string) *ChatRoom {
	return &ChatRoom{
		Id:            roomId,
		Members:       make(map[*Connection]bool),
		Broadcast:     make(chan structs.Message),
		StatusUpdates: make(chan structs.SeenUpdate),
		Register:      make(chan *Connection),
		Unregister:    make(chan *Connection),
	}
}

func (r *ChatRoom) Run(service *RoomService) {
	ctx := context.Background()
	for {
		select {
		case conn := <-r.Register:
			log.Printf("Registering connection to room %s, address: %p", r.Id, conn)
			r.Members[conn] = true
		case conn := <-r.Unregister:
			log.Printf("Unregistering connection from room %s", r.Id)
			if _, ok := r.Members[conn]; ok {
				delete(r.Members, conn)
				close(conn.sendMessage)
			}
		case message := <-r.Broadcast:
			log.Printf("Broadcasting message to room %s: %s", r.Id, string(message.Content))
			message, err := service.ProcessMessage(ctx, r.Id, &message)
			if err != nil {
				log.Printf("Error saving message %q in room %s", string(message.Content), r.Id)
				break
			}
			for conn := range r.Members {
				select {
				case conn.sendMessage <- message:
				default:
					close(conn.sendMessage)
					delete(r.Members, conn)
				}
			}
		case seenUpdate := <-r.StatusUpdates:
			log.Printf("Broadcasting seen update to room %s: %s", r.Id, seenUpdate.MessageId)

			err := service.ProcessSeenUpdate(ctx, r.Id, &seenUpdate)
			if err != nil {
				log.Printf("Error saving seen update for message %s in room %s", seenUpdate.MessageId, r.Id)
				break
			}

			for conn := range r.Members {
				select {
				case conn.sendSeenUpdate <- seenUpdate:
				default:
					close(conn.sendSeenUpdate)
					delete(r.Members, conn)
				}
			}
		}

	}
}
