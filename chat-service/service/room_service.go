package service

import (
	"context"
	"errors"

	"example.com/chat_app/chat_service/dto"
	"example.com/chat_app/chat_service/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrInsufficientPermissions = errors.New("insufficient permissions")

type RoomService struct {
	repo repository.ChatRoomRepository
}

func NewRoomService(repo repository.ChatRoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) GetRoomDto(ctx context.Context, roomId string, userId string) (*dto.RoomDto, error) {
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

func (s *RoomService) DeleteRoom(ctx context.Context, roomId string, userId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, userId); err != nil {
		return err
	}
	return s.repo.DeleteRoom(ctx, roomId)
}

func (s *RoomService) AddUserToRoom(ctx context.Context, roomId string, newUserId string, addingUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, addingUserId); err != nil {
		return err
	}
	userPermissions := repository.UserPermissions{
		UserId: newUserId,
		Role:   repository.Member,
	}
	return s.repo.InsertUserIntoRoom(ctx, roomId, userPermissions)
}

func (s *RoomService) AddUsersToRoom(ctx context.Context, roomId string, newUsers []string, addingUserId string) ([]error, error) {
	var dbInsertErrors []error
	if err := s.validateAdminPrivileges(ctx, roomId, addingUserId); err != nil {
		return nil, ErrInsufficientPermissions
	}
	for _, userId := range newUsers {
		permission := repository.UserPermissions{
			UserId: userId,
			Role:   repository.Member,
		}
		err := s.repo.InsertUserIntoRoom(ctx, roomId, permission)
		if err != nil {
			dbInsertErrors = append(dbInsertErrors, err)
		}
	}
	return dbInsertErrors, nil
}

func (s *RoomService) RemoveUserFromRoom(ctx context.Context, roomId, requestingUserId, removedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, requestingUserId); err != nil {
		return err
	}
	return s.repo.DeleteUserFromRoom(ctx, roomId, removedUserId)
}

func (s *RoomService) LeaveRoom(ctx context.Context, roomId, userId string) error {
	return s.repo.DeleteUserFromRoom(ctx, roomId, userId)
}

func (s *RoomService) PromoteUser(ctx context.Context, roomId string, promotingUserId, promotedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, promotingUserId); err != nil {
		return err
	}
	return s.repo.ChangeUserRole(ctx, roomId, promotedUserId, repository.Admin)
}

func (s *RoomService) DemoteUser(ctx context.Context, roomId, demotingUserId, demotedUserId string) error {
	if err := s.validateAdminPrivileges(ctx, roomId, demotingUserId); err != nil {
		return err
	}
	return s.repo.ChangeUserRole(ctx, roomId, demotedUserId, repository.Member)
}

func (s *RoomService) CreateRoom(ctx context.Context, userId string) (*repository.ChatRoomEntity, error) {
	room, err := s.repo.CreateRoom(ctx)
	if err != nil {
		return nil, err
	}
	err = s.AddAdminToRoom(ctx, room.Id, userId)
	if err != nil {
		return nil, err
	}
	return room, nil
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
		if err == mongo.ErrNoDocuments {
			return ErrInsufficientPermissions
		}
		return err
	}
	if userPermissions.Role != repository.Admin {
		return ErrInsufficientPermissions
	}
	return nil
}
