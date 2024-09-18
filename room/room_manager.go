package room

import (
	"fmt"
	"log"
	"sync"

	"example.com/myproject/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomManager interface {
	GetOrCreateRoom(roomId string, conn *Connection) (*ChatRoom, error)
}

type InMemoryRoomManager struct {
	rooms map[string]*ChatRoom
	lock sync.Mutex
	repo repository.ChatRoomRepository
}

func NewRoomManager(repo repository.ChatRoomRepository) *InMemoryRoomManager {
	return &InMemoryRoomManager{
		rooms: make(map[string]*ChatRoom),
		repo: repo,
	}
}

func (m *InMemoryRoomManager) GetOrCreateRoom(roomId string, conn *Connection) (*ChatRoom, error) {
    m.lock.Lock()
    defer m.lock.Unlock()

    roomEntity, err := m.repo.GetRoom(roomId)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Printf("Room: %s doesn't exist, creating new room", roomId)
            if _, err := m.repo.CreateRoom(roomId); err != nil {
                log.Printf("Failed to create room: %v", err)
                return nil, fmt.Errorf("failed to create room: %w", err)
            }
        } else {
            log.Printf("Error checking for existence of room: %v", err)
            return nil, fmt.Errorf("error checking for existence of room: %w", err)
        }
    }

    if room, exists := m.rooms[roomId]; exists {
        log.Printf("Room: %s exists, registering connection", roomId)
		room.Register <- conn
        for _, message := range roomEntity.Messages {
            log.Printf("Sending existing message to connection: %s", message.Content)
            conn.send <- []byte(message.Content)
        }
        return room, nil
    }

    // Create a new room and register the connection
    log.Printf("Creating new room with roomId: %s", roomId)
    newRoom := NewChatRoom(roomId, m.repo)
    m.rooms[roomId] = newRoom

    // Register new connection to the room and pump messages existing in the repository to the broadcast channel of the room
	go newRoom.Run(m.repo)
    newRoom.Register <- conn
	go func() {
		for _, message := range roomEntity.Messages {
			log.Printf("Broadcasting existing message to room: %s", message.Content)
			newRoom.Broadcast <- []byte(message.Content)
		}
	}()
    
    log.Printf("Room: %s created and running", roomId)

    return newRoom, nil
}


