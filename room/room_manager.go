package room

import (
	"context"
	"fmt"
	"log"
	"sync"

	"example.com/myproject/entity"
	"example.com/myproject/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomManager interface {
	GetOrCreateRoom(ctx context.Context, roomId string, conn *Connection) (*ChatRoom, error)
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

func (m *InMemoryRoomManager) GetOrCreateRoom(ctx context.Context, roomId string, conn *Connection) (*ChatRoom, error) {
    m.lock.Lock()
    defer m.lock.Unlock()

    roomEntity, err := m.repo.GetRoom(ctx, roomId)
    if err != nil {
        // If room doesn't exist in db, create a room in db and proceed
        if err == mongo.ErrNoDocuments {
            log.Printf("Room: %s doesn't exist in database, creating new room", roomId)
            if _, err := m.repo.CreateRoom(ctx, roomId); err != nil {
                return nil, fmt.Errorf("failed to create room: %w", err)
            }
        } else {
            return nil, fmt.Errorf("error checking for existence of room: %w", err)
        }
    }

    var room *ChatRoom
    // If room doesn't exist in the in-memory chat room manager, create a new room and run it
    room, exists := m.rooms[roomId]
    if !exists {
        log.Printf("Creating new room with roomId: %s", roomId)
        room = NewChatRoom(roomId, m.repo)
        m.rooms[roomId] = room
        go room.Run(m.repo)
    }

    // Register new connection to the room 
    log.Printf("Registering connection for user %s", conn.userId)
    room.Register <- conn

    // Pump messages existing in the repo to the new connection
	go func() {
		pumpExistingMessagesToNewConnection(conn, roomEntity.Messages)
	}()
    log.Printf("Room: %s running", roomId)

    return room, nil
}

func pumpExistingMessagesToNewConnection(conn *Connection, messages []entity.MessageEntity) {
    for _, message := range messages {
        conn.send <- []byte(message.Content)
    }
}


