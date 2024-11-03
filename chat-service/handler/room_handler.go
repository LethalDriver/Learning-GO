package handler

import (
	"net/http"

	"example.com/chat_app/chat_service/service"
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

func (rh *RoomHandler) AddUsersToRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	newUsersIds := r.URL.Query()["userId"]
	errsInsert, errPermission := rh.roomService.AddUsersToRoom(ctx, roomId, newUsersIds, addingUserId)
	if errPermission != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response := struct {
		Errors []error `json:"errors"`
	}{
		Errors: errsInsert,
	}
	writeJsonResponse(w, response)
	w.WriteHeader(http.StatusOK)
}
