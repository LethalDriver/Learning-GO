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
	port := 8080
	http.ListenAndServe(":8080", r)
	log.Default().Println("Server started on port", port)
}





