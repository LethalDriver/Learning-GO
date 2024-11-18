package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type MediaServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() (*MediaServiceClient, error) {
	baseURL := os.Getenv("MEDIA_SERVICE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("MEDIA_SERVICE_URL environment variable not set")
	}

	return &MediaServiceClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *MediaServiceClient) UploadMedia(ctx context.Context, roomId string, mediaType string, mediaBytes []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getMediaURL(roomId, mediaType, ""), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create upload image request: %v", err)
	}

	log.Printf("Sending request to %s", req.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send upload image request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code from upload image : %d", resp.StatusCode)
	}

	payload := struct {
		BlobId string `json:"blobId"`
	}{
		BlobId: "",
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read upload image response: %w", err)
	}

	if err := json.Unmarshal(respBytes, &payload); err != nil {
		return "", fmt.Errorf("failed to unmarshal upload image response: %w", err)
	}

	return payload.BlobId, nil
}

func (c *MediaServiceClient) DownloadMedia(ctx context.Context, roomId, blobId, mediaType string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getMediaURL(roomId, mediaType, blobId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download image request: %v", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send dowload image request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from download image request: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *MediaServiceClient) getMediaURL(roomId, mediaType, blobId string) string {
	if blobId == "" {
		return fmt.Sprintf("%s/%s/%s", c.BaseURL, roomId, mediaType)
	}
	return fmt.Sprintf("%s/%s/%s/%s", c.BaseURL, roomId, mediaType, blobId)
}
