package handler

import (
	"log"
	"net/http"

	"example.com/myproject/service"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	upgrader    websocket.Upgrader
	chatService *service.ChatService
	userService *service.UserService
}

func NewWebsocketHandler(rs *service.ChatService, us *service.UserService) *WebsocketHandler {
	return &WebsocketHandler{
		chatService: rs,
		userService: us,
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

	userId, err := service.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}
	log.Println("WebSocket connection upgraded successfully")

	wsh.chatService.ConnectToRoom(ctx, roomId, userId, conn)
}
