package service

import (
	"example.com/chat_app/user_service/dto"
	"example.com/chat_app/user_service/repository"
)

func MapUserEntityToDto(user *repository.UserEntity) *dto.UserDto {
	return &dto.UserDto{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}
