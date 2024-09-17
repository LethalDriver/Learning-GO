package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
    upgrader websocket.Upgrader
    m RoomManager
}

func NewWebsocketHandler(m RoomManager) *WebsocketHandler {
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

func (wsh *WebsocketHandler) handleWebSocketUpgradeRequest(w http.ResponseWriter, r *http.Request) {
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
        if err := handleConnection(conn, wsh.m, roomId); err != nil {
            log.Println("Error handling WebSocket connection:", err)
            conn.Close()
        }
    }()
}

func handleConnection(ws *websocket.Conn, m RoomManager, roomId string) error {
    log.Println("Handling connection")
    conn := &Connection{
        ws:   ws,
        send: make(chan []byte, 256),
    }

    room, err := m.GetOrCreateRoom(roomId, conn)
    if err != nil {
        return fmt.Errorf("error creating room: %v", err)
    }
    conn.room = room

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        conn.writePump()
    }()

    go func() {
        defer wg.Done()
        conn.readPump()
    }()

    wg.Wait()
    room.Unregister <- conn
    log.Printf("Connection unregistered from room ID: %s", roomId)

    return nil
}