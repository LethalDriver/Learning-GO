package service

import (
	"context"
	"errors"

	"example.com/chat_app/chat_service/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoomRepository interface {
	CreateRoom(ctx context.Context, name string) (*structs.ChatRoomEntity, error)
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

// ErrInsufficientPermissions is an error indicating that the user does not have sufficient permissions.
var ErrInsufficientPermissions = errors.New("insufficient permissions")

// RoomService provides methods to manage chat rooms and handle user permissions.
type RoomService struct {
	repo ChatRoomRepository
}

// NewRoomService creates a new instance of RoomService.
func NewRoomService(repo ChatRoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

// GetRoomDto retrieves a chat room DTO if the user belongs to the room.
func (s *RoomService) GetRoomDto(ctx context.Context, roomId string, userId string) (*structs.RoomDto, error) {
	room, err := s.repo.GetRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}
	if !checkIfUserBelongsToRoom(room, userId) {
		return nil, ErrInsufficientPermissions
	}
	roomDto := MapRoomEntityToDto(room)
	return roomDto, nil
}

// CreateRoom creates a new chat room and adds the creating user as an admin.
func (s *RoomService) CreateRoom(ctx context.Context, userId, name string) (*structs.ChatRoomEntity, error) {
	room, err := s.repo.CreateRoom(ctx, name)
	if err != nil {
		return nil, err
	}
	err = s.AddAdminToRoom(ctx, room.Id, userId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// DeleteRoom deletes a chat room if the user has admin privileges.
func (s *RoomService) DeleteRoom(ctx context.Context, roomId string, userId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, userId); err != nil {
		return err
	}
	return s.repo.DeleteRoom(ctx, roomId)
}

// AddUserToRoom adds a user to a chat room if the requesting user has admin privileges.
func (s *RoomService) AddUserToRoom(ctx context.Context, roomId string, newUserId string, addingUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, addingUserId); err != nil {
		return err
	}
	userPermissions := structs.UserPermissions{
		UserId: newUserId,
		Role:   structs.Member,
	}
	return s.repo.InsertUserIntoRoom(ctx, roomId, userPermissions)
}

// AddUsersToRoom adds multiple users to a chat room if the requesting user has admin privileges.
func (s *RoomService) AddUsersToRoom(ctx context.Context, roomId string, newUsers []string, addingUserId string) ([]error, error) {
	var dbInsertErrors []error
	if err := s.validateAdminPrivileges(ctx, roomId, addingUserId); err != nil {
		return nil, ErrInsufficientPermissions
	}
	for _, userId := range newUsers {
		permission := structs.UserPermissions{
			UserId: userId,
			Role:   structs.Member,
		}
		err := s.repo.InsertUserIntoRoom(ctx, roomId, permission)
		if err != nil {
			dbInsertErrors = append(dbInsertErrors, err)
		}
	}
	return dbInsertErrors, nil
}

// RemoveUserFromRoom removes a user from a chat room if the requesting user has admin privileges.
func (s *RoomService) RemoveUserFromRoom(ctx context.Context, roomId, requestingUserId, removedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, requestingUserId); err != nil {
		return err
	}
	return s.repo.DeleteUserFromRoom(ctx, roomId, removedUserId)
}

// LeaveRoom allows a user to leave a chat room.
func (s *RoomService) LeaveRoom(ctx context.Context, roomId, userId string) error {
	return s.repo.DeleteUserFromRoom(ctx, roomId, userId)
}

// PromoteUser promotes a user to admin in a chat room if the requesting user has admin privileges.
func (s *RoomService) PromoteUser(ctx context.Context, roomId string, promotingUserId, promotedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, promotingUserId); err != nil {
		return err
	}
	return s.repo.ChangeUserRole(ctx, roomId, promotedUserId, structs.Admin)
}

// DemoteUser demotes a user to member in a chat room if the requesting user has admin privileges.
func (s *RoomService) DemoteUser(ctx context.Context, roomId, demotingUserId, demotedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, demotingUserId); err != nil {
		return err
	}
	return s.repo.ChangeUserRole(ctx, roomId, demotedUserId, structs.Member)
}

// AddAdminToRoom adds an admin to a chat room.
func (s *RoomService) AddAdminToRoom(ctx context.Context, roomId string, userId string) error {
	userPermission := structs.UserPermissions{
		UserId: userId,
		Role:   structs.Admin,
	}
	return s.repo.InsertUserIntoRoom(ctx, roomId, userPermission)
}

// validateAdminPrivileges checks if a user has admin privileges in a chat room.
func (s *RoomService) validateAdminPrivileges(ctx context.Context, roomId, userId string) error {
	userPermissions, err := s.repo.GetUsersPermissions(ctx, roomId, userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrInsufficientPermissions
		}
		return err
	}
	if userPermissions.Role != structs.Admin {
		return ErrInsufficientPermissions
	}
	return nil
}
