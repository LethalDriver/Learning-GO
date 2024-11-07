package service

import (
	"example.com/chat_app/user_service/structs"
)

func MapUserEntityToDto(user *structs.UserEntity) *structs.UserDto {
	return &structs.UserDto{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}
