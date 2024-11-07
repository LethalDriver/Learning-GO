package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"time"

	"media_service/structs"

	"github.com/google/uuid"
)

var ErrPermissionDenied = errors.New("you don't have permissions")
var ErrFileNotInRoom = errors.New("does not belong to room")

type FileRepository interface {
	GetFile(ctx context.Context, id string, mediaType structs.MediaType) (*structs.MediaFile, error)
	DeleteFile(ctx context.Context, userId string, mediaType structs.MediaType) error
	SaveFile(ctx context.Context, file *structs.MediaFile, mediaType structs.MediaType) error
}

type FileService struct {
	repo    FileRepository
	storage MediaStorageService
}

func NewFileService(repo FileRepository, storage MediaStorageService) *FileService {
	return &FileService{repo: repo, storage: storage}
}

func (s *FileService) CreateFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, sentBy string, roomId string) (*structs.MediaFile, error) {
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

	mediaFile := &structs.MediaFile{
		Id:      id,
		Type:    mediaType,
		MediaId: mediaId,
		Metadata: &structs.FileMetadata{
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

func (s *FileService) DeleteFile(ctx context.Context, id string, roomId string, userId string, mediaType structs.MediaType) error {
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

func (s *FileService) GetFile(ctx context.Context, id string, roomId string, mediaType structs.MediaType) (*structs.MediaFile, []byte, error) {
	// Fetch the file metadata
	file, err := s.repo.GetFile(ctx, id, mediaType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed fetching file %q from database: %w", id, err)
	}
	if err := checkIfFileInRoom(file, roomId); err != nil {
		return nil, nil, err
	}

	imageData, err := s.storage.DownloadFile(ctx, file.Type, file.MediaId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed downloading image %q: %w", id, err)
	}

	return file, imageData, nil
}

func checkIfFileInRoom(file *structs.MediaFile, roomId string) error {
	if !slices.Contains(file.Metadata.RoomIds, roomId) {
		return fmt.Errorf("file %q %w: %q", file.Id, ErrFileNotInRoom, roomId)
	}
	return nil
}

func constructLocalUrl(id string, mediaType structs.MediaType) (string, error) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return "", fmt.Errorf("BASE_URL environment variable not set")
	}
	return fmt.Sprintf("%s/%s/%s", baseUrl, mediaType.String(), id), nil
}

func determineMediaType(filename string) structs.MediaType {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return structs.Image
	case ".mp4", ".avi", ".mov":
		return structs.Video
	case ".mp3", ".wav":
		return structs.Audio
	default:
		return structs.Other
	}
}
