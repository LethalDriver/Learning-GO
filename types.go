package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ws *websocket.Conn 
	send chan []byte
	room *ChatRoom
}

