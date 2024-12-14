package handler

import (
	"net/http"

	"example.com/chat_app/chat_service/service"
)

// RoomHandler handles room-related requests.
type RoomHandler struct {
	roomService  *service.RoomService
	mediaService *service.MediaService
}

// NewRoomHandler creates a new RoomHandler with the provided RoomService.
func NewRoomHandler(rs *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: rs}
}

// CreateRoom handles room creation requests.
// It creates a new room with a UUID and returns the id.
func (rh *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")
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

// GetRoom returns the room information for the specified room ID.
// It returns the room DTO if the user is a member of the room.
// If the user is not a member of the room, it returns a 403 Forbidden error.
func (rh *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")
	userId := r.Header.Get("X-User-Id")
	room, err := rh.roomService.GetRoomDto(ctx, roomId, userId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "User doesn't belong to room", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, room)
	w.WriteHeader(http.StatusOK)
}

// DeleteRoom deletes the room with the specified room ID.
// If requesting user is not its admin it returns a 403 Forbidden error.
func (rh *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")
	userId := r.Header.Get("X-User-Id")
	err := rh.roomService.DeleteRoom(ctx, roomId, userId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "This action requires admin privileges", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddUsersToRoom adds the specified users to the room with the specified room ID.
// It returns a list of errors for each user that could not be added to the room.
// If the requesting user is not an admin of the room, it returns a 403 Forbidden error.
func (rh *RoomHandler) AddUsersToRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	newUsersIds := r.URL.Query()["userId"]
	errsInsert, errPermission := rh.roomService.AddUsersToRoom(ctx, roomId, newUsersIds, addingUserId)
	if errPermission != nil {
		http.Error(w, "This action requires admin privileges", http.StatusForbidden)
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

// PromoteUser promotes the specified user to admin in the room with the specified room ID.
// If the requesting user is not an admin of the room, it returns a 403 Forbidden error.
func (rh *RoomHandler) PromoteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	promotingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	newAdminId := r.PathValue("userId")
	err := rh.roomService.PromoteUser(ctx, roomId, newAdminId, promotingUserId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "This action requires admin privileges", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DemoteUser demotes the specified user from admin in the room with the specified room ID.
// If the requesting user is not an admin of the room, it returns a 403 Forbidden error.
func (rh *RoomHandler) DemoteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	demotingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	demotedUserId := r.PathValue("userId")
	err := rh.roomService.DemoteUser(ctx, roomId, demotedUserId, demotingUserId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "This action requires admin privileges", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// LeaveRoom removes the requesting user from the room with the specified room ID.
func (rh *RoomHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	err := rh.roomService.LeaveRoom(ctx, roomId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteUserFromRoom removes the specified user from the room with the specified room ID.
// If the requesting user is not an admin of the room, it returns a 403 Forbidden error.
func (rh *RoomHandler) DeleteUserFromRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestingUserId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")
	removedUserId := r.PathValue("userId")
	err := rh.roomService.RemoveUserFromRoom(ctx, roomId, requestingUserId, removedUserId)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			http.Error(w, "This action requires admin privileges", http.StatusForbidden)
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

