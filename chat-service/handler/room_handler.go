package handler

import (
	"fmt"
	"net/http"

	"example.com/chat_app/chat_service/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(rs *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: rs}
}

func (rh *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")
	//Repository import, should not be there, there should be a separate package for structs used across all layers
	room, err := rh.roomService.CreateRoom(ctx, userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		RoomId string `json:"roomId"`
	}{
		RoomId: room.Id,
	}
	writeJsonResponse(w, response)
	w.WriteHeader(http.StatusCreated)
}
