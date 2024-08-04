package main

type ChatRoom struct {
	ID         string
	Members    map[*Connection]bool
	Broadcast  chan []byte
	Register   chan *Connection
	Unregister chan *Connection
}

func (room *ChatRoom) Run() {
	for {
		select {
		case conn := <-room.Register:
			room.Members[conn] = true
		case conn := <-room.Unregister:
			if _, ok := room.Members[conn]; ok {
				delete(room.Members, conn)
				close(conn.send)
			}
		case message := <-room.Broadcast:
			for conn := range room.Members {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(room.Members, conn)
				}
			}
		}
	}
}
