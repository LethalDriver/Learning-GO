package handler

import (
	"log"
	"net/http"

	"example.com/chat_app/chat_service/repository"
	"example.com/chat_app/chat_service/service"
	"example.com/chat_app/common"
	"github.com/gorilla/websocket"
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

	userId, ok := ctx.Value(common.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	username, ok := ctx.Value(common.UsernameKey).(string)
	if !ok {
		http.Error(w, "Username not found in context", http.StatusInternalServerError)
		return
	}

	user := repository.UserDetails{
		Id:       userId,
		Username: username,
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}
	log.Println("WebSocket connection upgraded successfully")

	wsh.chatService.ConnectToRoom(ctx, roomId, user, conn)
}
