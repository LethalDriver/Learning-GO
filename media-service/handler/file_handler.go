package handler

import (
	"encoding/json"
	"fmt"
	"media_service/repository"
	"media_service/service"
	"net/http"
)

type FileHandler struct {
	service service.FileService
}

func (h *FileHandler) HandleMediaUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")
	roomId := r.PathValue("roomId")

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

func (h *FileHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaTypeStr := r.PathValue("mediaType")
	roomId := r.PathValue("roomId")
	fileId := r.PathValue("fileId")

	mediaType, err := repository.ParseMediaType(mediaTypeStr)
	if err != nil {
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}

	fileMetadata, fileData, err := h.service.GetFile(ctx, fileId, roomId, mediaType)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; mediaId=%s", fileMetadata.MediaId))

	// Create a JSON response with metadata and binary data
	response := struct {
		Metadata *repository.MediaFile `json:"metadata"`
		File     []byte                `json:"file"`
	}{
		Metadata: fileMetadata,
		File:     fileData,
	}

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
