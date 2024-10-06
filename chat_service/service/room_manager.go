package service

import (
	"log"
	"sync"
)

type RoomManager interface {
	ManageRoom(roomId string) *ChatRoom
}

type InMemoryRoomManager struct {
	rooms map[string]*ChatRoom
	lock sync.Mutex
}

func NewRoomManager() *InMemoryRoomManager {
	return &InMemoryRoomManager{
		rooms: make(map[string]*ChatRoom),
	}
}

func (m *InMemoryRoomManager) ManageRoom(roomId string) *ChatRoom {
	m.lock.Lock()
    defer m.lock.Unlock()
	var room *ChatRoom
    room, exists := m.rooms[roomId]
    if !exists {
        log.Printf("Creating new room with roomId: %s", roomId)
        room = NewChatRoom(roomId)
        m.rooms[roomId] = room
    }
	return room
}




