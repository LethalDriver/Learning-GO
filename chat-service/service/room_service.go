package service

import (
	"context"
	"errors"

	"example.com/chat_app/chat_service/repository"
)

var ErrInsufficientPermissions = errors.New("insufficient permissions")

type RoomService struct {
	repo repository.ChatRoomRepository
}

func NewRoomService(repo repository.ChatRoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(ctx context.Context, roomId string) (*repository.ChatRoomEntity, error) {
	return s.repo.CreateRoom(ctx, roomId)
}

func (s *RoomService) GetRoom(ctx context.Context, roomId string) (*repository.ChatRoomEntity, error) {
	return s.repo.GetRoom(ctx, roomId)
}

func (s *RoomService) AddUserToRoom(ctx context.Context, roomId string, userId string) error {
	userPermissions := repository.UserPermissions{
		UserId: userId,
		Role:   repository.Member,
	}
	return s.repo.InsertUserIntoRoom(ctx, roomId, userPermissions)
}

func (s *RoomService) RemoveUserFromRoom(ctx context.Context, roomId, requestingUserId, removedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, requestingUserId); err != nil {
		return err
	}
	return s.repo.DeleteUserFromRoom(ctx, roomId, removedUserId)
}

func (s *RoomService) PromoteUserToAdmin(ctx context.Context, roomId string, promotingUserId, promotedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, promotingUserId); err != nil {
		return err
	}
	return s.repo.PromoteUserToAdmin(ctx, roomId, promotedUserId)
}

func (s *RoomService) MakeUserAdmin(ctx context.Context, roomId string, userId string) error {
	userPermissions := repository.UserPermissions{
		UserId: userId,
		Role:   repository.Admin,
	}
	return s.repo.InsertUserIntoRoom(ctx, roomId, userPermissions)
}

func (s *RoomService) AddAdminToRoom(ctx context.Context, roomId string, userId string) error {
	userPermission := repository.UserPermissions{
		UserId: userId,
		Role:   repository.Admin,
	}
	return s.repo.InsertUserIntoRoom(ctx, roomId, userPermission)
}
func (s *RoomService) validateAdminPrivileges(ctx context.Context, roomId, userId string) error {
	userPermissions, err := s.repo.GetUsersPermissions(ctx, roomId, userId)
	if err != nil {
		if err == repository.ErrEntityNotFound {
			return ErrInsufficientPermissions
		}
		return err
	}
	if userPermissions.Role != repository.Admin {
		return ErrInsufficientPermissions
	}
	return nil
}
