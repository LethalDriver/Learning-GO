package structs

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

type WsIncomingMessage struct {
	Type MessageType `json:"type"`
	Data any         `json:"data"`
}

type Message struct {
	Id            string         `json:"id"`
	Content       string         `json:"content"`
	EmbeddedMedia *EmbeddedMedia `json:"embeddedMedia"`
	SentBy        UserDetails    `json:"sentBy"`
	SentAt        time.Time      `json:"sentAt"`
	SeenBy        []UserDetails  `json:"seenBy"`
}

type EmbeddedMedia struct {
	ContentType string `json:"contentType"`
	Url         string `json:"url"`
}

type UserDetails struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type SeenMessage struct {
	MessageId string      `json:"messageId"`
	SeenBy    UserDetails `json:"seenBy"`
}

type DeleteMessage struct {
	MessageId string `json:"messageId"`
	SentBy    string `json:"sentBy"`
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

func DetermineMessageType(messageBytes []byte) (MessageType, error) {
	var msg WsIncomingMessage
	err := json.Unmarshal(messageBytes, &msg)
	if err != nil {
		return -1, err
	}
	switch msg.Type {
	case TypeTextMessage:
		return TypeTextMessage, nil
	case TypeSeenMessage:
		return TypeSeenMessage, nil
	case TypeDeleteMessage:
		return TypeDeleteMessage, nil
	default:
		return -1, errors.New("unknown message type")
	}
}
