package service

import (
	"bytes"
	"context"
	"fmt"
	"media_service/repository"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type MediaType int

const (
	Image MediaType = iota
	Video
	Audio
	Other
)

func (m MediaType) String() string {
	return [...]string{"images", "videos", "audios", "others"}[m]
}

type MediaStorageService interface {
	DownloadFile(ctx context.Context, mediaType repository.MediaType, mediaId string) ([]byte, error)
	UploadFile(ctx context.Context, mediaType repository.MediaType, blobId string, data []byte) error
}

type AzureBlobStorageService struct {
	serviceClient *azblob.Client
}

func NewAzureBlobStorageService(accountName, accountKey string) (*AzureBlobStorageService, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	serviceClient, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create service client: %w", err)
	}

	return &AzureBlobStorageService{
		serviceClient: serviceClient,
	}, nil
}

func (s *AzureBlobStorageService) DownloadFile(ctx context.Context, mediaType repository.MediaType, mediaId string) ([]byte, error) {
	get, err := s.serviceClient.DownloadStream(ctx, mediaType.String()+"s", mediaId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob: %w", err)
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from retry reader: %w", err)
	}

	err = retryReader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close retry reader: %w", err)
	}

	return downloadedData.Bytes(), nil
}

func (s *AzureBlobStorageService) UploadFile(ctx context.Context, mediaType repository.MediaType, blobId string, data []byte) error {
	_, err := s.serviceClient.UploadStream(ctx, mediaType.String()+"s", blobId, bytes.NewReader(data), nil)
	if err != nil {
		return fmt.Errorf("failed to upload blob: %w", err)
	}

	return nil
}
