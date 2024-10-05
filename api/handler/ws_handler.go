package handler

import (
	"log"
	"net/http"

	"example.com/myproject/service"
	"example.com/myproject/structs"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
    upgrader websocket.Upgrader
    rs *service.RoomService
	us *service.UserService
}

func NewWebsocketHandler(rs *service.RoomService, us *service.UserService) *WebsocketHandler {
    return &WebsocketHandler{
        rs: rs,
		us: us,
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
    log.Printf("Upgrading HTTP connection to WebSocket for room ID: %s", roomId)

    userId, err := service.GetUserIdFromContext(r.Context())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
	user, err := wsh.us.GetUserById(ctx, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	userDetails := structs.UserDetails{
		Id:       user.Id,
		Username: user.Username,
	}

    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := wsh.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade connection to WebSocket:", err)
        return
    }
    log.Println("WebSocket connection upgraded successfully")

    // Handle the WebSocket connection
    go func() {
        if err := service.HandleConnection(conn, wsh.rs, roomId, userDetails); err != nil {
            log.Println("Error handling WebSocket connection:", err)
            conn.Close()
        }
    }()
}

