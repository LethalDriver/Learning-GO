package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// MediaServiceClient is an http client wrapper for communication with media service.
type MediaServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new MediaServiceClient.
// It reads the base URL for the media service from the MEDIA_SERVICE_URL environment variable.
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

// UploadMedia uploads media to the media service.
// It sends a POST request to the media service with the media type and media bytes.
// It returns the blob ID of the uploaded media.
func (c *MediaServiceClient) UploadMedia(ctx context.Context, mediaType string, mediaBytes []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getMediaURL(mediaType, ""), bytes.NewReader(mediaBytes))
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

// DownloadMedia downloads media from the media service.
// It sends a GET request to the media service with the blob ID and media type.
// It returns the media bytes.
func (c *MediaServiceClient) DownloadMedia(ctx context.Context, blobId, mediaType string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getMediaURL(mediaType, blobId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download image request: %v", err)
	}

	log.Printf("Sending request to %s", req.URL.String())

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

// getMediaURL returns the URL for the media service endpoint.
func (c *MediaServiceClient) getMediaURL(mediaType, blobId string) string {
	if blobId == "" {
		return fmt.Sprintf("%s/%s", c.BaseURL, mediaType)
	}
	return fmt.Sprintf("%s/%s/%s", c.BaseURL, mediaType, blobId)
}
