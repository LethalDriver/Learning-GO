package main

import (
	"github.com/google/uuid"
)

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