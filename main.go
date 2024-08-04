package main

import (
	"log"
	"net/http"
	"sync"

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
	r.HandleFunc("/ws/{room_id}", func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Error upgrading connection")
			return
		}
		handleConnection(ws, m, req)
	})
	log.Default().Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":8080", r))
}


func handleConnection(ws *websocket.Conn, m *RoomManager, req *http.Request) {
	vars := mux.Vars(req)
    roomId := vars["room_id"]

	room := m.GetOrCreateRoom(roomId)

	conn := &Connection{
		ws: ws,
		send: make(chan []byte, 256),
		room: room,
	}

	room.Register <- conn

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

	wg.Wait()

	room.Unregister <- conn
}






