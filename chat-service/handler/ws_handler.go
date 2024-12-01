package handler

import (
	"log"
	"net/http"

	"example.com/chat_app/chat_service/service"
	"github.com/gorilla/websocket"
)

// WebsocketHandler handles WebSocket connections.
type WebsocketHandler struct {
	upgrader    websocket.Upgrader
	chatService *service.ChatService
}

// NewWebsocketHandler creates a new WebsocketHandler with the provided ChatService.
// It initializes the WebSocket upgrader with default values.
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

// HandleWebSocketUpgradeRequest handles WebSocket upgrade requests.
// It validates the connection and upgrades the HTTP connection to a WebSocket connection.
// If the connection request points to a non-existent room, it returns a 404 Not Found error.
// If the user is not a member of a room, it returns a 401 Unauthorized error.
func (wsh *WebsocketHandler) HandleWebSocketUpgradeRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")
	userId := r.Header.Get("X-User-Id")

	if err := wsh.chatService.ValidateConnection(ctx, roomId, userId); err != nil {
		log.Println("Failed to validate connection:", err)
		switch err {
		case service.ErrRoomNotFound:
			http.Error(w, "Room not found", http.StatusNotFound)
		case service.ErrInsufficientPermissions:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
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
