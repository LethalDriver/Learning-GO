package structs

import "time"

type MediaType int

const (
	Audio MediaType = iota
	Image
	Video
	Unviewable
)

type MediaFile struct {
	Id       string       `bson:"id" json:"id"`
	Type     MediaType    `bson:"type" json:"type"`
	Url      string       `bson:"url" json:"url"`
	Metadata FileMetadata `bson:"metadata" json:"metadata"`
}

type FileMetadata struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	Size      int64     `bson:"size" json:"size"`
}
