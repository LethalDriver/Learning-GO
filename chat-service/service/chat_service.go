package service

import (
	"context"
	"log"
	"time"

	"example.com/chat_app/chat_service/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
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

func (s *ChatService) ConnectToRoom(ctx context.Context, roomId, userId string, ws *websocket.Conn) error {
	dbRoom, err := s.roomRepo.GetRoom(ctx, roomId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}
	if !checkIfUserBelongsToRoom(ctx, dbRoom, userId) {
		return ErrInsufficientPermissions
	}
	userDetails := repository.UserDetails{
		Id: userId,
	}
	memoryRoom := s.roomManager.ManageRoom(dbRoom.Id)
	go memoryRoom.Run(s)

	go handleConnection(ws, memoryRoom, userDetails)

	log.Printf("Room: %s running", roomId)
	return nil
}

func checkIfUserBelongsToRoom(ctx context.Context, room *repository.ChatRoomEntity, userId string) bool {
	for _, user := range room.Users {
		if user.UserId == userId {
			return true
		}
	}
	return false
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
