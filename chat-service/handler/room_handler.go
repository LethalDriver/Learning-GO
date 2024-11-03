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

func (rh *RoomHandler) MakeUserAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	promotingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	newAdminId := r.PathValue("userId")
	err := rh.roomService.MakeUserAdmin(ctx, roomId, newAdminId, promotingUserId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rh *RoomHandler) DeleteUserFromRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	removedUserId := r.PathValue("userId")
	err := rh.roomService.RemoveUserFromRoom(ctx, roomId, requestingUserId, removedUserId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := struct {
		Message string `json:"message"`
	}{
		Message: "User removed from room",
	}
	writeJsonResponse(w, response)
	w.WriteHeader(http.StatusOK)
}
