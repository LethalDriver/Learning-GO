package structs

import (
	"time"
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
    Id         string    `bson:"id" json:"id"`
    Content    string    `bson:"content" json:"content"`
    ChatRoomId string    `bson:"chatRoomId" json:"chatRoomId"`
    Type       string    `bson:"messageType" json:"messageType"`
    SentBy     string    `bson:"sentBy" json:"sentBy"`
    SentAt     time.Time `bson:"sentAt" json:"sentAt"`
    SeenBy     []string  `bson:"seenBy" json:"seenBy"`
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
