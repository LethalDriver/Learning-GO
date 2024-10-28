package client

import (
	"fmt"
	"net/http"
	"os"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() (*Client, error) {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("USER_SERVICE_URL environment variable not set")
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *Client) PostPermissionUpdate() error {
	req, err := http.NewRequest("POST", c.BaseURL+"/permission-update", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
