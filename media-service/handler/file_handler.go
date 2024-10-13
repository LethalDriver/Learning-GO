package handler

import (
	"media_service/service"
	"net/http"
)

type FileHandler struct {
	service service.FileService
}

func (h *FileHandler) HandleMediaUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")
	roomId := r.Header.Get("X-Room-Id")

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	createdFile, err := h.service.CreateFile(ctx, file, header, userId, roomId)
	if err != nil {
		http.Error(w, "Unable to upload file", http.StatusInternalServerError)
		return
	}
	writeResponse(w, createdFile)
	w.WriteHeader(http.StatusCreated)
}
