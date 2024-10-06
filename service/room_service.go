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
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomService struct {
	roomRepo    repository.ChatRoomRepository
	userRepo    repository.UserRepository
	roomManager RoomManager
}

func NewRoomService(roomRepo *repository.MongoChatRoomRepository, userRepo *repository.MongoUserRepository, roomManager RoomManager) *RoomService {
	return &RoomService{
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		roomManager: roomManager,
	}
}

func (s *RoomService) GetOrCreateRoom(ctx context.Context, roomId string, conn *Connection) (*ChatRoom, error) {
	roomEntity, err := s.getOrCreateRoomEntity(ctx, roomId)
	if err != nil {
		return nil, err
	}

	room := s.roomManager.ManageRoom(roomId)
	go room.Run(s)

	log.Printf("Registering connection for user %s", conn.user.Id)
	room.Register <- conn

	if roomEntity != nil {
		go s.mapAndPumpMessages(ctx, conn, roomEntity.Messages)
	}

	log.Printf("Room: %s running", roomId)
	return room, nil
}

func (s *RoomService) getOrCreateRoomEntity(ctx context.Context, roomId string) (*structs.ChatRoomEntity, error) {
	roomEntity, err := s.roomRepo.GetRoom(ctx, roomId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Room: %s doesn't exist in database, creating new room", roomId)
			if roomEntity, err = s.roomRepo.CreateRoom(ctx, roomId); err != nil {
				return nil, fmt.Errorf("failed to create room: %w", err)
			}
			return roomEntity, nil
		}
		return nil, fmt.Errorf("error checking for existence of room in database: %w", err)
	}
	return roomEntity, nil
}

func (s *RoomService) mapAndPumpMessages(ctx context.Context, conn *Connection, messageEntities []structs.MessageEntity) {
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

func (s *RoomService) ProcessAndSaveMsg(ctx context.Context, roomId string, message *structs.Message) (structs.Message, error) {
	message.Id = uuid.New().String()
	message.SentAt = time.Now()
	message.SeenBy = []structs.UserDetails{}
	entity := mappers.MapMessageToEntity(message, roomId)
	return *message, s.roomRepo.AddMessageToRoom(ctx, roomId, entity)
}
