package main

import (
	"sync"
)

type RoomManager struct {
	rooms map[string]*ChatRoom
	lock sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*ChatRoom),
	}
}

func (manager *RoomManager) GetOrCreateRoom(roomId string) *ChatRoom {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	
	if room, exists := manager.rooms[roomId]; exists {
		return room
	}

	newRoom := &ChatRoom{
		ID: roomId,
		Members: make(map[*Connection]bool),
		Broadcast: make(chan []byte),
		Register: make(chan *Connection),
		Unregister: make(chan *Connection),
	}
	manager.rooms[roomId] = newRoom

	go newRoom.Run()

	return newRoom
}