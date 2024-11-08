package service

import (
	"context"

	"example.com/chat_app/chat_service/structs"
)

type MediaRepository interface {
	GetFile(ctx context.Context, id string) (*structs.MediaFile, error)
	DeleteFile(ctx context.Context, id string) error
	SaveFile(ctx context.Context, file *structs.MediaFile) error
}

type MediaService struct {
	repo MediaRepository
}

func NewMediaService(repo MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

func (s *MediaService) GetMedia(ctx context.Context, id string) (*structs.MediaFile, error) {
	return s.repo.GetFile(ctx, id)
}
