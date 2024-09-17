package main

import "log"

type ChatRoom struct {
	Id         string
	Members    map[*Connection]bool
	Broadcast  chan []byte
	Register   chan *Connection
	Unregister chan *Connection
}

func (r *ChatRoom) Run(repo ChatRoomRepository) {
    log.Printf("Room %s is running", r.Id)
    for {
        select {
        case conn := <-r.Register:
            log.Printf("Registering connection to room %s", r.Id)
            r.Members[conn] = true
        case conn := <-r.Unregister:
            log.Printf("Unregistering connection from room %s", r.Id)
            if _, ok := r.Members[conn]; ok {
                delete(r.Members, conn)
                close(conn.send)
            }
        case message := <-r.Broadcast:
            log.Printf("Broadcasting message to room %s: %s", r.Id, string(message))
            for conn := range r.Members {
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
