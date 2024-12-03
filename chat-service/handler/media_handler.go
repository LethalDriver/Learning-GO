package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"example.com/chat_app/chat_service/service"
)

// MediaHandler handles media upload and download requests.
type MediaHandler struct {
	mediaService *service.MediaService
}

// NewMediaHandler creates a new MediaHandler with the provided MediaService.
func NewMediaHandler(ms *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: ms}
}

// UploadMedia handles media upload requests.
// It reads the file from the multipart form data and uploads it to the specified room and media type.
// It reads the binary file data from the "file" field in the form data.
// It returns the metadata of the uploaded media.
func (mh *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	if roomId == "" {
		http.Error(w, "Missing roomId query parameter", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // Limit of 10 MB
	if err != nil {
		http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read image file", http.StatusInternalServerError)
	}

	mediaType := r.FormValue("mediaType")

	// Pass the file to the service
	fileDocument, err := mh.mediaService.CreateMediaResource(ctx, roomId, mediaType, userId, fileBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, fileDocument)
	w.WriteHeader(http.StatusCreated)
}

// GetMediaMetadata returns the metadata of the specified media as a JSON response.
func (mh *MediaHandler) GetMediaMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaId := r.PathValue("mediaId")

	fileMetadata, err := mh.mediaService.GetMediaMetadata(ctx, mediaId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileMetadata); err != nil {
		http.Error(w, "Failed to encode metadata", http.StatusInternalServerError)
		return
	}
}

// GetMediaFile retrieves the binary image data from the media service and returns it with the appropriate content type.
func (mh *MediaHandler) GetMediaFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mediaId := r.PathValue("mediaId")

	fileBytes, err := mh.mediaService.GetMediaBinary(ctx, mediaId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(fileBytes) < 512 {
		fileBytes = append(fileBytes, make([]byte, 512-len(fileBytes))...)
	}
	contentType := http.DetectContentType(fileBytes[:512])

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(fileBytes); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}
}
