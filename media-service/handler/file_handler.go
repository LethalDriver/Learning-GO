package handler

import (
	"io"
	"media_service/service"
	"net/http"
)

// FileHandler handles file upload and download requests.
type FileHandler struct {
	service *service.AzureBlobStorageService
}

// NewFileHandler creates a new FileHandler with the provided AzureBlobStorageService.
func NewFileHandler(s *service.AzureBlobStorageService) *FileHandler {
	return &FileHandler{service: s}
}

// HandleMediaUpload handles file upload requests.
// It reads the file from the request body and uploads it to the specified media type container.
func (h *FileHandler) HandleMediaUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaType := r.PathValue("mediaType")

	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	blobId, err := h.service.UploadFile(ctx, mediaType, fileBytes)
	if err != nil {
		http.Error(w, "Unable to upload file to storage", http.StatusInternalServerError)
		return
	}

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

// HandleMediaDownload handles file download requests.
// It retrieves the file from the specified media type container and writes it to the response.
func (h *FileHandler) HandleMediaDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaType := r.PathValue("mediaType")
	blobId := r.PathValue("blobId")

	fileBytes, err := h.service.DownloadFile(ctx, mediaType, blobId)
	if err != nil {
		http.Error(w, "Unable to download file from storage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
