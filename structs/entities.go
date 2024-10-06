package structs

import (
	"time"
)

type Entity interface {
	GetId() string
}

type ChatRoomEntity struct {
	Id       string          `bson:"id" json:"id"`
	Messages []MessageEntity `bson:"messages" json:"messages"`
}

type MessageEntity struct {
	Id         string    `bson:"id" json:"id"`
	Content    string    `bson:"content" json:"content"`
	ChatRoomId string    `bson:"chatRoomId" json:"chatRoomId"`
	SentBy     string    `bson:"sentBy" json:"sentBy"`
	SentAt     time.Time `bson:"sentAt" json:"sentAt"`
	SeenBy     []string  `bson:"seenBy" json:"seenBy"`
}


func (message MessageEntity) GetId() string {
	return message.Id
}
func (chat ChatRoomEntity) GetId() string {
	return chat.Id
}
