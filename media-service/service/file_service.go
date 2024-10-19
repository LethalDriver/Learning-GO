package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"media_service/repository"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/google/uuid"
)

var ErrPermissionDenied = errors.New("you don't have permissions")
var ErrFileNotInRoom = errors.New("does not belong to room")

type FileService struct {
	repo    repository.FileRepository
	storage MediaStorageService
}

func (s *FileService) CreateFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, sentBy string, roomId string) (*repository.MediaFile, error) {
	id := uuid.New().String()
	mediaId := uuid.New().String()
	mediaType := determineMediaType(header.Filename)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q content: %w", header.Filename, err)
	}

	err = s.storage.UploadFile(ctx, mediaType, mediaId, fileContent)
	if err != nil {
		return nil, fmt.Errorf("unable to save file to storage: %w", err)
	}

	mediaFile := &repository.MediaFile{
		Id:      id,
		Type:    mediaType,
		MediaId: mediaId,
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

	return mediaFile, nil

}

func (s *FileService) DeleteFile(ctx context.Context, id string, roomId string, userId string, mediaType repository.MediaType) error {
	file, err := s.repo.GetFile(ctx, id, mediaType)
	if err != nil {
		return fmt.Errorf("failed fetching file %q from database: %w", id, err)
	}
	if file.Metadata.CreatedBy != userId {
		return fmt.Errorf("as %q: %w to delete this file", userId, ErrPermissionDenied)
	}
	if err := checkIfFileInRoom(file, roomId); err != nil {
		return err
	}
	return s.repo.DeleteFile(ctx, userId, mediaType)
}

func (s *FileService) GetFile(ctx context.Context, id string, roomId string, mediaType repository.MediaType, w http.ResponseWriter) error {
	// Fetch the file metadata
	file, err := s.repo.GetFile(ctx, id, mediaType)
	if err != nil {
		return fmt.Errorf("failed fetching file %q from database: %w", id, err)
	}
	if err := checkIfFileInRoom(file, roomId); err != nil {
		return err
	}

	imageData, err := s.storage.DownloadFile(ctx, mediaType, file.MediaId)
	if err != nil {
		return fmt.Errorf("failed downloading image %q: %w", id, err)
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; mediaId=%s", file.MediaId))

	// Create a JSON response with metadata and binary data
	response := struct {
		Metadata *repository.MediaFile `json:"metadata"`
		Image    []byte                `json:"image"`
	}{
		Metadata: file,
		Image:    imageData,
	}

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("failed encoding response: %w", err)
	}

	return nil
}

func checkIfFileInRoom(file *repository.MediaFile, roomId string) error {
	if !slices.Contains(file.Metadata.RoomIds, roomId) {
		return fmt.Errorf("file %q %w: %q", file.Id, ErrFileNotInRoom, roomId)
	}
	return nil
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
