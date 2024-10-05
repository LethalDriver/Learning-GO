package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

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

func (mt MessageType) String() string {
	switch mt {
	case TextMessage:
		return "TextMessage"
	case ImageMessage:
		return "ImageMessage"
	default:
		return "Unknown"
	}
}

func MessageTypeFromString(s string) (MessageType, error) {
	switch s {
	case "TextMessage":
		return TextMessage, nil
	case "ImageMessage":
		return ImageMessage, nil
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
        *mt = TextMessage
    case "ImageMessage":
        *mt = ImageMessage
    default:
        return errors.New("invalid MessageType")
    }

    return nil
}
