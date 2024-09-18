package entity

import (
	"github.com/google/uuid"
)

type Entity interface {
	GetId() string
}

type UserEntity struct {
	Id string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type ChatRoomEntity struct {
	Id       string          `bson:"id" json:"id"`
	Messages []MessageEntity `bson:"messages" json:"messages"`
}


type MessageEntity struct {
	Id       string         `bson:"id" json:"id"`
	Content  string         `bson:"content" json:"content"`
	ChatRoomId string `bson:"chatRoomId" json:"chatRoomId"`
}

func NewMessageEntity(content string, chatRoomId string) *MessageEntity {
	return &MessageEntity{
		Id:       uuid.New().String(),
		Content:  content,
		ChatRoomId: chatRoomId,
	}
}

func NewUserEntity(username string, email string, password string) *UserEntity {
    return &UserEntity{
        Id:       uuid.New().String(),
        Username: username,
        Email:    email,
        Password: password,
    }
}

func (user UserEntity) GetId() string {
	return user.Id
}
func (message MessageEntity) GetId() string {
	return message.Id
}
func (chat ChatRoomEntity) GetId() string {
	return chat.Id
}
