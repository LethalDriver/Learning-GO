package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)


var (upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}
repo *RedisMessageRepository
)
func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var msg Message
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Println(err)
			return
		}
		repo.Create(&msg)
		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}

func handleMessage(message []byte) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println(err)
		return
	}
	repo.Create(&msg)
}

func main() {
	r := mux.NewRouter()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	repo = NewRedisMessageRepository(rdb)
	port := 8080
	r.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", r)
	log.Default().Println("Server started on port", port)
}

func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var room Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repo.CreateRoom(&room)
	w.WriteHeader(http.StatusCreated)
}

