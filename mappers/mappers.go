package mappers

import (
	"context"
	"fmt"
	"log"

	"example.com/myproject/repository"
	"example.com/myproject/structs"
	"example.com/myproject/utils"
)

func MapMessageToEntity(message *structs.Message, chatRoomId string) *structs.MessageEntity {
	return &structs.MessageEntity{
		Id:         message.Id,
		Content:    message.Content,
		ChatRoomId: chatRoomId,
		Type:       message.Type.String(),
		SentBy:     message.SentBy.Id,
		SentAt:     message.SentAt,
		SeenBy:     utils.MapSlice(message.SeenBy, func(user structs.UserDetails) string { return user.Id }),
	}
}

func MapEntityToMessage(ctx context.Context, entity *structs.MessageEntity, userRepo repository.UserRepository) (*structs.Message, error) {
	msgType, err := structs.MessageTypeFromString(entity.Type)
	if err != nil {
		log.Printf("Error converting message type: %v", err)
	}

	sentByEntity, err := userRepo.GetById(ctx, entity.SentBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	sentByDetails := structs.UserDetails{
		Id:       sentByEntity.Id,
		Username: sentByEntity.Username,
	}

	seenBy := make([]structs.UserDetails, 0)
	for _, userId := range entity.SeenBy {
		user, err := userRepo.GetById(ctx, userId)
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
		SeenBy:  seenBy,
	}, nil
}
