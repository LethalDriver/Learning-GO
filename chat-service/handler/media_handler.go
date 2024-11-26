package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"example.com/chat_app/chat_service/service"
)

type MediaHandler struct {
	mediaService *service.MediaService
}

func NewMediaHandler(ms *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: ms}
}

func (mh *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")
	userId := r.Header.Get("X-User-Id")
	mediaType := r.PathValue("mediaType")

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

	// Pass the file to the service
	fileDocument, err := mh.mediaService.CreateMediaResource(ctx, roomId, mediaType, userId, fileBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, fileDocument)
	w.WriteHeader(http.StatusCreated)
}

func (mh *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomId := r.PathValue("roomId")
	mediaId := r.PathValue("mediaId")
	mediaTypeStr := r.PathValue("mediaType")

	// Retrieve the media and metadata from the service
	fileMetadata, fileBytes, err := mh.mediaService.GetMedia(ctx, mediaId, mediaTypeStr, roomId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mw.FormDataContentType())

	metadataPart, err := mw.CreateFormField("metadata")
	if err != nil {
		http.Error(w, "Unable to create metadata part", http.StatusInternalServerError)
		return
	}
	metadataBytes, err := json.Marshal(fileMetadata)
	if err != nil {
		http.Error(w, "Unable to marshal metadata", http.StatusInternalServerError)
		return
	}
	_, err = metadataPart.Write(metadataBytes)
	if err != nil {
		http.Error(w, "Unable to write metadata part", http.StatusInternalServerError)
		return
	}

	filePart, err := mw.CreateFormFile("file", fileMetadata.BlobId)
	if err != nil {
		http.Error(w, "Unable to create file part", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(filePart, bytes.NewReader(fileBytes))
	if err != nil {
		http.Error(w, "Unable to write file part", http.StatusInternalServerError)
		return
	}

	err = mw.Close()
	if err != nil {
		http.Error(w, "Unable to close multipart writer", http.StatusInternalServerError)
		return
	}
}