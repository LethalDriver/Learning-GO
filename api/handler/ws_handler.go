package handler

import (
	"log"
	"net/http"

	"example.com/myproject/room"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
    upgrader websocket.Upgrader
    m room.RoomManager
}

func NewWebsocketHandler(m room.RoomManager) *WebsocketHandler {
    return &WebsocketHandler{
        m: m,
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
    roomId := r.PathValue("roomId")
    log.Printf("Upgrading HTTP connection to WebSocket for room ID: %s", roomId)

    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := wsh.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade connection to WebSocket:", err)
        return
    }
    log.Println("WebSocket connection upgraded successfully")

    // Handle the WebSocket connection
    go func() {
        if err := room.HandleConnection(conn, wsh.m, roomId); err != nil {
            log.Println("Error handling WebSocket connection:", err)
            conn.Close()
        }
    }()
}

