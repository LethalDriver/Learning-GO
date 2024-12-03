package service

import (
	"context"
	"log"
	"time"

	"example.com/chat_app/chat_service/structs"
	"github.com/google/uuid"
)

// MediaRepository provides methods to interact with the media file storage.
type MediaRepository interface {
	GetFile(ctx context.Context, id string) (*structs.MediaFile, error)
	DeleteFile(ctx context.Context, id string) error
	SaveFile(ctx context.Context, file *structs.MediaFile) error
}

// Client is an interface for interacting with the media storage service.
type Client interface {
	UploadMedia(ctx context.Context, mediaType string, mediaBytes []byte) (string, error)
	DownloadMedia(ctx context.Context, blobId, mediaType string) ([]byte, error)
}

// MediaService provides methods to manage media files.
type MediaService struct {
	repo   MediaRepository
	client Client
}

// NewMediaService creates a new instance of MediaService.
// It takes a MediaRepository and a Client as dependencies.
func NewMediaService(repo MediaRepository, client Client) *MediaService {
	return &MediaService{
		repo:   repo,
		client: client,
	}
}

// CreateMediaResource creates a new media resource.
// It uploads the media to the media service and saves the metadata in the repository.
func (s *MediaService) CreateMediaResource(ctx context.Context, roomId, mediaTypeStr, userId string, mediaBytes []byte) (*structs.MediaFile, error) {
	log.Printf("Uploading media of type %s\n to room %s", mediaTypeStr, roomId)
	blobId, err := s.client.UploadMedia(ctx, mediaTypeStr, mediaBytes)
	if err != nil {
		return nil, err
	}
	mediaType, err := structs.ParseMediaType(mediaTypeStr)
	if err != nil {
		return nil, err
	}
	file := &structs.MediaFile{
		Id:        uuid.New().String(),
		RoomId:    roomId,
		Type:      mediaType,
		CreatedAt: time.Now(),
		BlobId:    blobId,
		CreatedBy: userId,
	}
	err = s.repo.SaveFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetMediaMetadata retrieves a media file by its ID.
// It downloads the media from the media service and returns the metadata and binary data.
func (s *MediaService) GetMediaMetadata(ctx context.Context, id string) (*structs.MediaFile, error) {
	fileMetadata, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return nil, err
	}
	return fileMetadata, nil
}

// GetMediaBinary retrieves the binary data of a media file by its ID.
// It downloads the media from the media service and returns the binary data.
func (s *MediaService) GetMediaBinary(ctx context.Context, id string) ([]byte, error) {
	fileMetadata, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return nil, err
	}
	imageBytes, err := s.client.DownloadMedia(ctx, fileMetadata.BlobId, fileMetadata.Type.String())
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}
