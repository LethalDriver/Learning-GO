package service

import (
	"context"
	"fmt"
	"io"
	"media_service/repository"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	image   repository.ImageRepository
	storage StorageService
}

func (s *FileService) createFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, sentBy string, roomId string) (string, error) {
	id := uuid.New().String()
	mediaType := determineMediaType(header.Filename)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("unable to read file %q content: %w", header.Filename, err)
	}

	url, err := s.storage.UploadFile(fileContent)
	if err != nil {
		return "", fmt.Errorf("unable to save file to storage: %w", err)
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
	err = s.image.CreateImage(ctx, mediaFile)
	if err != nil {
		return "", fmt.Errorf("unable to save media file metadata: %w", err)
	}

	return mediaFile.Id, nil

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
