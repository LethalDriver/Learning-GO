package handler

import (
	"net/http"

	"example.com/chat_app/chat_service/service"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(cs *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: cs}
}

func (ch *ChatHandler) GetMessagesSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")
	userId := r.Header.Get("X-User-Id")

	messagesSummary, err := ch.chatService.GetMessagesSummary(ctx, roomId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, messagesSummary)
	w.WriteHeader(http.StatusOK)
}
