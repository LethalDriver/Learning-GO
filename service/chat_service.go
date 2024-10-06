package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/myproject/mappers"
	"example.com/myproject/repository"
	"example.com/myproject/structs"
	"example.com/myproject/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	roomRepo    repository.ChatRoomRepository
	userRepo    repository.UserRepository
	roomManager RoomManager
}

func NewChatService(roomRepo *repository.MongoChatRoomRepository, userRepo *repository.MongoUserRepository, roomManager RoomManager) *ChatService {
	return &ChatService{
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		roomManager: roomManager,
	}
}

func (s *ChatService) ConnectToRoom(ctx context.Context, roomId string, userId string, ws *websocket.Conn) error {
	_, err := s.roomRepo.CreateRoom(ctx, roomId)
	if err != nil {
		if err != repository.ErrRoomExists {
			return err
		}
	}
	room := s.roomManager.ManageRoom(roomId)
	go room.Run(s)

	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}
	details := mappers.MapEntityToUserDetails(user)

	go handleConnection(ws, room, details)

	log.Printf("Room: %s running", roomId)
	return nil
}

func (s *ChatService) mapAndPumpMessages(ctx context.Context, conn *Connection, messageEntities []structs.MessageEntity) {
	messages := utils.MapSlice(messageEntities, func(entity structs.MessageEntity) structs.Message {
		message, err := mappers.MapEntityToMessage(ctx, &entity, s.userRepo)
		if err != nil {
			log.Printf("Error mapping message entity to message: %v", err)
		}
		return *message
	})
	for _, message := range messages {
		conn.sendMessage <- message
	}
}

func (s *ChatService) processAndSaveMessage(ctx context.Context, roomId string, message *structs.Message) (structs.Message, error) {
	message.Id = uuid.New().String()
	message.SentAt = time.Now()
	message.SeenBy = []structs.UserDetails{}
	entity := mappers.MapMessageToEntity(message, roomId)
	return *message, s.roomRepo.AddMessageToRoom(ctx, roomId, entity)
}
