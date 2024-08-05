package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main() {
	mongoClientOption := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), mongoClientOption)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongo!")

	repo := NewMongoChatRoomRepository(client, "chatdb", "chatrooms")


	r := mux.NewRouter()
	m := NewRoomManager()
	port := 8080
	r.HandleFunc("/ws/{room_id}", func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Error upgrading connection")
			return
		}
		handleConnection(ws, m, req, repo)
	})
	log.Default().Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":8080", r))
}


func handleConnection(ws *websocket.Conn, m *RoomManager, req *http.Request, repo ChatRoomRepository) {
	vars := mux.Vars(req)
    roomId := vars["room_id"]

	conn := &Connection{
		ws: ws,
		send: make(chan []byte, 256),
	}

	room := m.GetOrCreateRoom(roomId, repo, conn)

	conn.room = room

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






