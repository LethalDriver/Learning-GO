package service

import (
	"context"

	"example.com/chat_app/chat_service/structs"
)

// ChatRoomRepository provides methods to interact with chat rooms.
type ChatRoomRepository interface {
	CreateRoom(ctx context.Context) (*structs.ChatRoomEntity, error)
	GetRoom(ctx context.Context, id string) (*structs.ChatRoomEntity, error)
	DeleteRoom(ctx context.Context, id string) error
	AddMessageToRoom(ctx context.Context, roomId string, message *structs.Message) error
	InsertSeenBy(ctx context.Context, roomId string, messageId string, userId string) error
	DeleteMessage(ctx context.Context, roomId string, messageId string) error
	InsertUserIntoRoom(ctx context.Context, roomId string, user structs.UserPermissions) error
	DeleteUserFromRoom(ctx context.Context, roomId string, userId string) error
	GetUsersPermissions(ctx context.Context, roomId string, userId string) (*structs.UserPermissions, error)
	ChangeUserRole(ctx context.Context, roomId string, userId string, role structs.Role) error
	GetUnseenMessages(ctx context.Context, roomId, userId string) ([]structs.Message, error)
}
