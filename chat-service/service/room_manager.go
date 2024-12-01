package service

import (
	"log"
	"sync"
)

// RoomManager is an interface for managing chat rooms.
type RoomManager interface {
	ManageRoom(roomId string) *ChatRoom
}

// InMemoryRoomManager is an implementation of RoomManager that stores chat rooms in memory.
type InMemoryRoomManager struct {
	rooms map[string]*ChatRoom
	lock  sync.Mutex
}

// NewRoomManager creates a new instance of InMemoryRoomManager.
func NewRoomManager() *InMemoryRoomManager {
	return &InMemoryRoomManager{
		rooms: make(map[string]*ChatRoom),
	}
}

// ManageRoom creates a new chat room or returns an existing one.
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
