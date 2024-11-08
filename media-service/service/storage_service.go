package service

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

type AzureBlobStorageService struct {
	serviceClient *azblob.Client
}

func NewAzureBlobStorageService() (*AzureBlobStorageService, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
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

func (s *AzureBlobStorageService) DownloadFile(ctx context.Context, containerName, mediaId string) ([]byte, error) {
	get, err := s.serviceClient.DownloadStream(ctx, containerName, mediaId, nil)
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

func (s *AzureBlobStorageService) UploadFile(ctx context.Context, containerName string, data []byte) (string, error) {
	blobId := uuid.NewString()
	_, err := s.serviceClient.UploadStream(ctx, containerName, blobId, bytes.NewReader(data), nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	return blobId, nil
}
