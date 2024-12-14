package client

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"example.com/chat_app/chat_service/structs"
)

// AiAssistantClient is an http client wrapper for communication with media service.
type AiAssistantClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAiClient creates a new AiAssistantClient.
// It reads the base URL for the media service from the AI_ASSISTANT_URL environment variable.
func NewAiClient() (*AiAssistantClient, error) {
	baseURL := os.Getenv("AI_ASSISTANT_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_ASSISTANT_URL environment variable not set")
	}

	return &AiAssistantClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *AiAssistantClient) GetMessagesSummary(ctx context.Context, messages []structs.Message) (*structs.MessagesSummary, error) {
	return &structs.MessagesSummary{
		Summary: "This is a summary of the messages",
	}, nil
}
