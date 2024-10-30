package service

import (
	"context"
	"log"
	"time"

	"example.com/chat_app/chat_service/exception"
	"example.com/chat_app/chat_service/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	roomRepo    repository.ChatRoomRepository
	roomManager RoomManager
}

func NewChatService(roomRepo *repository.MongoChatRoomRepository, roomManager RoomManager) *ChatService {
	return &ChatService{
		roomRepo:    roomRepo,
		roomManager: roomManager,
	}
}

func (s *ChatService) ConnectToRoom(ctx context.Context, roomId string, user repository.UserDetails, ws *websocket.Conn) error {
	_, err := s.roomRepo.CreateRoom(ctx)
	if err != nil {
		if err != exception.ErrRoomExists {
			return err
		}
	}
	room := s.roomManager.ManageRoom(roomId)
	go room.Run(s)

	go handleConnection(ws, room, user)

	log.Printf("Room: %s running", roomId)
	return nil
}

func (s *ChatService) pumpExistingMessages(conn *Connection, messages []repository.Message) {
	for _, message := range messages {
		conn.sendMessage <- message
	}
}

func (s *ChatService) processAndSaveMessage(ctx context.Context, roomId string, message *repository.Message) (repository.Message, error) {
	message.Id = uuid.New().String()
	message.SentAt = time.Now()
	message.ChatRoomId = roomId
	message.SeenBy = []repository.UserDetails{message.SentBy}
	return *message, s.roomRepo.AddMessageToRoom(ctx, roomId, message)
}
