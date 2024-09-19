package room

import "time"

type MessageType int

const (
	TextMessage MessageType = iota
	ImageMessage
)

type UserDetails struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type SeenUpdate struct {
	MessageId string `json:"messageId"`
	SeenBy UserDetails `json:"seenBy"`
}

type Message struct {
	Id     string      `json:"id"`
	Type MessageType `json:"messageType"`
	Content string `json:"content"`
	SentBy UserDetails `json:"sentBy"`
	SentAt time.Time `json:"sentAt"`
	SeenBy []UserDetails `json:"seenBy"`
}
