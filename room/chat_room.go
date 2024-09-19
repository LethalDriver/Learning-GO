package room

import (
	"context"
	"log"

	"example.com/myproject/repository"
)

type ChatRoom struct {
	Id         string
	Members    map[*Connection]bool
	Broadcast  chan []byte
	Register   chan *Connection
	Unregister chan *Connection
	repo repository.ChatRoomRepository
}

func NewChatRoom(roomId string, repo repository.ChatRoomRepository) *ChatRoom {
	return &ChatRoom{
        Id:         roomId,
        Members:    make(map[*Connection]bool),
        Broadcast:  make(chan []byte),
        Register:   make(chan *Connection),
        Unregister: make(chan *Connection),
		repo: repo,
    }
}

func (r *ChatRoom) Run(repo repository.ChatRoomRepository) {
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
            log.Printf("Broadcasting message to room %s: %s", r.Id, string(message))
            ctx := context.Background()
			err := repo.AddMessageToRoom(ctx, r.Id, string(message))
			if err != nil {
				log.Printf("Error saving message %q in room %s", string(message), r.Id)
				break
			}
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
