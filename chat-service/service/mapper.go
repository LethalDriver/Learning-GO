package service

import (
	"example.com/chat_app/chat_service/dto"
	"example.com/chat_app/chat_service/repository"
)

func MapRoomEntityToDto(room *repository.ChatRoomEntity) *dto.RoomDto {
	members := make([]dto.UserDto, 0)
	for _, member := range room.Users {
		members = append(members, *MapUserPermissionsToDto(&member))
	}
	return &dto.RoomDto{
		Id:      room.Id,
		Members: members,
	}
}

func MapUserPermissionsToDto(user *repository.UserPermissions) *dto.UserDto {
	return &dto.UserDto{
		Id:   user.UserId,
		Role: user.Role.String(),
	}
}
