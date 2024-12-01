package structs

import (
	"time"
)

type MediaType int

const (
	Image MediaType = iota
	Video
	Audio
	Other
)

func (m MediaType) String() string {
	switch m {
	case Image:
		return "image"
	case Video:
		return "video"
	case Audio:
		return "audio"
	default:
		return "other"
	}
}

func ParseMediaType(s string) (MediaType, error) {
	switch s {
	case "image":
		return Image, nil
	case "video":
		return Video, nil
	case "audio":
		return Audio, nil
	default:
		return Other, nil
	}
}

type MediaFile struct {
	Id        string    `bson:"id" json:"id"`
	RoomId    string    `bson:"roomId" json:"roomId"`
	Type      MediaType `bson:"type" json:"type"`
	BlobId    string    `bson:"blobId" json:"blobId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	Size      int64     `bson:"size" json:"size"`
}
