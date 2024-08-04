package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main() {
	r := mux.NewRouter()
	m := NewRoomManager()
	port := 8080
	http.HandleFunc("/ws/{room_id}", func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		handleConnection(ws, m, req)
	})
	http.ListenAndServe(":8080", r)
	log.Default().Println("Server started on port", port)
}


func handleConnection(ws *websocket.Conn, m *RoomManager, req *http.Request) {
	vars := mux.Vars(req)
    roomId := vars["room_id"]

	room := m.GetOrCreateRoom(roomId)

	connection := &Connection{
		ws: ws,
		send: make(chan []byte, 256),
		room: room,
	}

	room.Register <- connection
}






