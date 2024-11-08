package handler

import (
	"io"
	"media_service/service"
	"net/http"
)

type FileHandler struct {
	service *service.AzureBlobStorageService
}

func NewFileHandler(s *service.AzureBlobStorageService) *FileHandler {
	return &FileHandler{service: s}
}

func (h *FileHandler) HandleMediaUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaType := r.URL.Query().Get("mediaType")

	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Stream the file directly to Azure Blob Storage
	blobId, err := h.service.UploadFile(ctx, mediaType, fileBytes)
	if err != nil {
		http.Error(w, "Unable to upload file to storage", http.StatusInternalServerError)
		return
	}

	// Create a response with the blob ID
	response := struct {
		BlobId string `json:"blobId"`
	}{
		BlobId: blobId,
	}

	if err := writeJsonResponse(w, response, http.StatusCreated); err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *FileHandler) HandleMediaDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaType := r.URL.Query().Get("mediaType")
	blobId := r.URL.Query().Get("blobId")

	fileBytes, err := h.service.DownloadFile(ctx, mediaType, blobId)
	if err != nil {
		http.Error(w, "Unable to download file from storage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
