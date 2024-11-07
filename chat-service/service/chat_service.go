package service

import (
	"context"
	"errors"
	"log"
	"time"

	"example.com/chat_app/chat_service/structs"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var RoomNotFound = errors.New("room not found")

type ChatService struct {
	roomRepo    ChatRoomRepository
	roomManager RoomManager
}

func NewChatService(roomRepo ChatRoomRepository, roomManager RoomManager) *ChatService {
	return &ChatService{
		roomRepo:    roomRepo,
		roomManager: roomManager,
	}
}

func (s *ChatService) ConnectToRoom(ctx context.Context, roomId, userId string, ws *websocket.Conn) {
	memoryRoom := s.roomManager.ManageRoom(roomId)
	go memoryRoom.Run(s)

	userDetails := structs.UserDetails{
		Id: userId,
	}
	go handleConnection(ws, memoryRoom, userDetails)

	log.Printf("Room: %s running", roomId)
}

func (s *ChatService) ValidateConnection(ctx context.Context, roomId, userId string) error {
	dbRoom, err := s.roomRepo.GetRoom(ctx, roomId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return RoomNotFound
		}
		return err
	}
	if !checkIfUserBelongsToRoom(dbRoom, userId) {
		return ErrInsufficientPermissions
	}
	return nil
}

func checkIfUserBelongsToRoom(room *structs.ChatRoomEntity, userId string) bool {
	for _, user := range room.Users {
		if user.UserId == userId {
			return true
		}
	}
	return false
}

func (s *ChatService) pumpExistingMessages(conn *Connection, messages []structs.Message) {
	for _, message := range messages {
		conn.sendMessage <- message
	}
}

func (s *ChatService) processAndSaveMessage(ctx context.Context, roomId string, message *structs.Message) (structs.Message, error) {
	message.Id = uuid.New().String()
	message.SentAt = time.Now()
	message.ChatRoomId = roomId
	message.SeenBy = []structs.UserDetails{message.SentBy}
	return *message, s.roomRepo.AddMessageToRoom(ctx, roomId, message)
}
