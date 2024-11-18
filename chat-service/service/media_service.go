package service

import (
	"context"
	"time"

	"example.com/chat_app/chat_service/structs"
	"github.com/google/uuid"
)

type MediaRepository interface {
	GetFile(ctx context.Context, id string) (*structs.MediaFile, error)
	DeleteFile(ctx context.Context, id string) error
	SaveFile(ctx context.Context, file *structs.MediaFile) error
}

type Client interface {
	UploadMedia(ctx context.Context, mediaType string, mediaBytes []byte) (string, error)
	DownloadMedia(ctx context.Context, blobId, mediaType string) ([]byte, error)
}

type MediaService struct {
	repo   MediaRepository
	client Client
}

func NewMediaService(repo MediaRepository, client Client) *MediaService {
	return &MediaService{
		repo:   repo,
		client: client,
	}
}

func (s *MediaService) CreateMediaResource(ctx context.Context, roomId, mediaTypeStr, userId string, mediaBytes []byte) (*structs.MediaFile, error) {
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

func (s *MediaService) GetMedia(ctx context.Context, id, mediaTypeStr, roomId string) (*structs.MediaFile, []byte, error) {
	imageBytes, err := s.client.DownloadMedia(ctx, id, mediaTypeStr)
	if err != nil {
		return nil, nil, err
	}
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return file, imageBytes, nil
}
