package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/myproject/repository"
	"example.com/myproject/structs"
	"example.com/myproject/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomService struct {
	roomRepo repository.ChatRoomRepository
	userRepo repository.UserRepository
	roomManager RoomManager
}

func NewRoomService(roomRepo *repository.MongoChatRoomRepository, userRepo *repository.MongoUserRepository, roomManager RoomManager) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
		userRepo: userRepo,
		roomManager: roomManager,
	}
}

func (s *RoomService) GetOrCreateRoom(ctx context.Context, roomId string, conn *Connection) (*ChatRoom, error) {
    roomEntity, err := s.GetRoom(ctx, roomId)
    if err != nil {
        // If room doesn't exist in db, create a room in db and proceed
        if err == mongo.ErrNoDocuments {
            log.Printf("Room: %s doesn't exist in database, creating new room", roomId)
            if _, err := s.CreateRoom(ctx, roomId); err != nil {
                return nil, fmt.Errorf("failed to create room: %w", err)
            }
        } else {
            return nil, fmt.Errorf("error checking for existence of room in database: %w", err)
        }
    }
	room := s.roomManager.ManageRoom(roomId)
	go room.Run(s)

    log.Printf("Registering connection for user %s", conn.user.Id)
    room.Register <- conn

    if roomEntity != nil {
		go func() {
			messages := make([]structs.Message, 0)
			for _, msgEntity := range roomEntity.Messages {
				message, err := s.MapEntityToMessage(ctx, &msgEntity)
				if err != nil {
					log.Printf("Error mapping message entity to message: %v", err)
					return
				}
				messages = append(messages, *message)
			}
			pumpExistingMessagesToNewConnection(conn, messages)
		}()
	}
    log.Printf("Room: %s running", roomId)

    return room, nil
}

func (s *RoomService) GetRoom(ctx context.Context, roomId string) (*structs.ChatRoomEntity, error) {
	return s.roomRepo.GetRoom(ctx, roomId)
}

func (s *RoomService) CreateRoom(ctx context.Context, roomId string) (*structs.ChatRoomEntity, error) {
	return s.roomRepo.CreateRoom(ctx, roomId)
}

func (s *RoomService) ProcessMessage(ctx context.Context, roomId string, message *structs.Message) (string, error) {
	message.SentAt = time.Now()
	message.Id = uuid.New().String()
	entity := MapMessageToEntity(message, roomId)
	return entity.Id, s.roomRepo.AddMessageToRoom(ctx, roomId, entity)
}

func MapMessageToEntity(message *structs.Message, chatRoomId string) *structs.MessageEntity {
	return &structs.MessageEntity{
		Id:      message.Id,
		Content: message.Content,
		ChatRoomId: chatRoomId,
		Type:    message.Type.String(),
		SentBy:  message.SentBy.Id,
		SentAt:  message.SentAt,
		SeenBy:  utils.MapSlice(message.SeenBy, func(user structs.UserDetails) string { return user.Id }),
	}
}

func (s *RoomService) MapEntityToMessage(ctx context.Context, entity *structs.MessageEntity) (*structs.Message, error) {
	msgType, err := structs.MessageTypeFromString(entity.Type)
	if err != nil {
		log.Printf("Error converting message type: %v", err)
	}

	sentByEntity, err := s.userRepo.GetById(ctx, entity.SentBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	sentByDetails := structs.UserDetails{
		Id:       sentByEntity.Id,
		Username: sentByEntity.Username,
	}

	seenBy := make([]structs.UserDetails, 0)
	for _, userId := range entity.SeenBy {
		user, err := s.userRepo.GetById(ctx, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get seen by: %w", err)
		}
		seenBy = append(seenBy, structs.UserDetails{
			Id:       user.Id,
			Username: user.Username,
		})
	}

	return &structs.Message{
		Id:      entity.Id,
		Type:    msgType,
		Content: entity.Content,
		SentBy:  sentByDetails,
		SentAt:  entity.SentAt,
		SeenBy: seenBy,
	}, nil
}

func pumpExistingMessagesToNewConnection(conn *Connection, messages []structs.Message) {
    for _, message := range messages {
        conn.send <- message
    }
}


