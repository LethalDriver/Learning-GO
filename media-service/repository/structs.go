package repository

import (
	"io"
	"time"
)

type Downloadable interface {
	Download() (io.ReadCloser, error)
}

type MediaType int

const (
	Image MediaType = iota
	Video
	Audio
	Other
)

type MediaFile struct {
	Id       string        `bson:"id" json:"id"`
	Type     MediaType     `bson:"type" json:"type"`
	Url      string        `bson:"url" json:"url"`
	Metadata *FileMetadata `bson:"metadata" json:"metadata"`
	RoomId   string        `bson:"roomId" json:"roomId"`
}

type FileMetadata struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	Size      int64     `bson:"size" json:"size"`
}
