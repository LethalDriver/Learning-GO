package handler

import (
	"log"
	"net/http"

	"example.com/chat_app/chat_service/service"
	"github.com/gorilla/websocket"
)

type ContextKey string

const (
	UserIdKey   ContextKey = "userId"
	UsernameKey ContextKey = "username"
)

type WebsocketHandler struct {
	upgrader    websocket.Upgrader
	chatService *service.ChatService
}

func NewWebsocketHandler(rs *service.ChatService) *WebsocketHandler {
	return &WebsocketHandler{
		chatService: rs,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (wsh *WebsocketHandler) HandleWebSocketUpgradeRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")

	userId := r.Header.Get("X-User-Id")
	username := r.Header.Get("X-Username")

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}
	log.Println("WebSocket connection upgraded successfully")

	wsh.chatService.ConnectToRoom(ctx, roomId, userId, username, conn)
}
