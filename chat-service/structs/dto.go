package structs

import "encoding/json"

type RoomDto struct {
	Id      string    `json:"id"`
	Members []UserDto `json:"members"`
}

type UserDto struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

type SeenMessage struct {
	MessageId string      `json:"messageId"`
	SeenBy    UserDetails `json:"seenBy"`
}

type DeleteMessage struct {
	MessageId string      `json:"messageId"`
	SentBy    UserDetails `json:"sentBy"`
}

type WsIncomingMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}
