package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"media_service/repository"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

var ErrPermissionDenied = errors.New("you don't have permissions")

type FileService struct {
	repo    repository.FileRepository
	storage StorageService
}

func (s *FileService) CreateFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, sentBy string, roomId string) (*repository.MediaFile, error) {
	id := uuid.New().String()
	mediaType := determineMediaType(header.Filename)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q content: %w", header.Filename, err)
	}

	url, err := s.storage.UploadFile(fileContent)
	if err != nil {
		return nil, fmt.Errorf("unable to save file to storage: %w", err)
	}

	mediaFile := &repository.MediaFile{
		Id:   id,
		Type: mediaType,
		Url:  url,
		Metadata: &repository.FileMetadata{
			CreatedAt: time.Now(),
			CreatedBy: sentBy,
			Size:      header.Size,
			RoomIds:   []string{roomId},
		},
	}

	// Save the media file metadata to the repository
	err = s.repo.SaveFile(ctx, mediaFile, mediaType)
	if err != nil {
		return nil, fmt.Errorf("unable to save media file metadata: %w", err)
	}

	mediaFile.Url, err = constructLocalUrl(id, mediaType)
	if err != nil {
		return nil, fmt.Errorf("unable to construct local URL: %w", err)
	}

	return mediaFile, nil

}

func (s *FileService) DeleteFile(ctx context.Context, id string, userId string, mediaType repository.MediaType) error {
	file, err := s.repo.GetFile(ctx, id, mediaType)
	if err != nil {
		return fmt.Errorf("failed fetching file %q from database: %w", id, err)
	}
	if file.Metadata.CreatedBy != userId {
		return fmt.Errorf("as %q: %w to delete this file", userId, err)
	}
	return s.repo.DeleteFile(ctx, userId, mediaType)
}

func constructLocalUrl(id string, mediaType repository.MediaType) (string, error) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return "", fmt.Errorf("BASE_URL environment variable not set")
	}
	return fmt.Sprintf("%s/%s/%s", baseUrl, mediaType.String(), id), nil
}

func determineMediaType(filename string) repository.MediaType {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return repository.Image
	case ".mp4", ".avi", ".mov":
		return repository.Video
	case ".mp3", ".wav":
		return repository.Audio
	default:
		return repository.Other
	}
}
