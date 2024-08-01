package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
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
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var msg Message
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Println(err)
			return
		}
		err = repo.Create(&msg)
		if err != nil {
			log.Println(err)
		}
		messageType := websocket.TextMessage
		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		} else {
			log.Printf("Message content: %s", msg.Content)
		}
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	repo = NewRedisMessageRepository(rdb)
	port := 8080
	log.Default().Println("Server started on port", port)
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

