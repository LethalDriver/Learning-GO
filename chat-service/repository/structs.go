package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type MessageType int

const (
	TypeTextMessage MessageType = iota
	TypeSeenMessage
	TypeDeleteMessage
)

type Role int

const (
	Member Role = iota
	Admin
)

type WsIncomingMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type UserPermissions struct {
	UserId string `json:"userId"`
	Role   Role   `json:"permission"`
}

type ChatRoomEntity struct {
	Id       string            `bson:"id" json:"id"`
	Messages []Message         `bson:"messages" json:"messages"`
	Users    []UserPermissions `bson:"users" json:"users"`
}

type Message struct {
	Id            string         `bson:"id" json:"id"`
	Content       string         `bson:"content" json:"content"`
	EmbeddedMedia *EmbeddedMedia `bson:"embeddedMedia" json:"embeddedMedia"`
	ChatRoomId    string         `bson:"chatRoomId" json:"chatRoomId"`
	SentBy        UserDetails    `bson:"sentBy" json:"sentBy"`
	SentAt        time.Time      `bson:"sentAt" json:"sentAt"`
	SeenBy        []UserDetails  `bson:"seenBy" json:"seenBy"`
}

type EmbeddedMedia struct {
	ContentType string `bson:"contentType" json:"contentType"`
	Url         string `bosn:"url" json:"url"`
}

type UserDetails struct {
	Id string `bson:"id" json:"id"`
}

type SeenMessage struct {
	MessageId string      `json:"messageId"`
	SeenBy    UserDetails `json:"seenBy"`
}

type DeleteMessage struct {
	MessageId string      `json:"messageId"`
	SentBy    UserDetails `json:"sentBy"`
}

func (mt MessageType) String() string {
	switch mt {
	case TypeTextMessage:
		return "TextMessage"
	case TypeSeenMessage:
		return "SeenMessage"
	case TypeDeleteMessage:
		return "DeleteMessage"
	default:
		return "Unknown"
	}
}

func MessageTypeFromString(s string) (MessageType, error) {
	switch s {
	case "TextMessage":
		return TypeTextMessage, nil
	case "SeenMessage":
		return TypeSeenMessage, nil
	case "DeleteMessage":
		return TypeDeleteMessage, nil
	default:
		return -1, fmt.Errorf("unknown message type: %s", s)
	}
}

func (mt MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

func (mt *MessageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "TextMessage":
		*mt = TypeTextMessage
	case "SeenMessage":
		*mt = TypeSeenMessage
	case "DeleteMessage":
		*mt = TypeDeleteMessage
	default:
		return errors.New("invalid MessageType")
	}

	return nil
}

func (r Role) String() string {
	switch r {
	case Member:
		return "Member"
	case Admin:
		return "Admin"
	default:
		return "Unknown"
	}
}

func RoleFromString(s string) (Role, error) {
	switch s {
	case "Member":
		return Member, nil
	case "Admin":
		return Admin, nil
	default:
		return -1, fmt.Errorf("unknown role: %s", s)
	}
}

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "Member":
		*r = Member
	case "Admin":
		*r = Admin
	default:
		return errors.New("invalid Role")
	}

	return nil
}
