package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/chat_app/user_service/service"
	"example.com/chat_app/user_service/structs"
)

type UserHandler struct {
	s *service.UserService
}

type LoginResponse struct {
	User        structs.UserDto `json:"user"`
	AccessToken string          `json:"accessToken"`
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		s: s,
	}
}

// HandleMe returns the user information for the authenticated user. The user ID is read from the request headers
// appended by the API Gateway from the JWT.
func (h *UserHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.Header.Get("X-User-Id")
	userDto, err := h.s.GetUserDto(ctx, userId)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = writeJsonResponse(w, userDto)
	if err != nil {
		log.Printf("Failed writing user response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleRegister registers a new user and returns the user DTO and JWT token.
func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var regReq service.RegistrationRequest
	err := parseRequest(r, &regReq)
	if err != nil {
		log.Printf("Failed registering user %q: %v", regReq.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userDto, token, err := h.s.RegisterUser(ctx, regReq)
	if err != nil {
		if err == service.ErrUserExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		log.Printf("Failed registering user %q: %v", regReq.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	resp := &LoginResponse{
		AccessToken: token,
		User:        *userDto,
	}
	err = writeJsonResponse(w, resp)
	if err != nil {
		log.Printf("Failed registering user %q: %v", regReq.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleLogin logs in a user and returns the user DTO and JWT token.
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var logReq service.LoginRequest
	err := parseRequest(r, &logReq)
	if err != nil {
		log.Printf("Failed parsing login request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dto, token, err := h.s.LoginUser(ctx, logReq)
	if err != nil {
		switch err {
		case service.ErrNoUser:
			http.Error(w, "User not found", http.StatusNotFound)
		case service.ErrWrongPassword:
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := &LoginResponse{
		User:        *dto,
		AccessToken: token,
	}
	err = writeJsonResponse(w, resp)
	if err != nil {
		log.Printf("Failed writing login response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// parseRequest reads the request body and unmarshals it into the given struct.
func parseRequest(r *http.Request, reqStruct any) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed reading body: %w", err)
	}

	err = json.Unmarshal(bodyBytes, reqStruct) // Unmarshal into the pointer
	if err != nil {
		return fmt.Errorf("failed parsing body: %w", err)
	}

	return nil
}

// writeJsonResponse writes the given struct as JSON to the response and sets the content type.
func writeJsonResponse(w http.ResponseWriter, respStruct any) error {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(respStruct)
	if err != nil {
		return fmt.Errorf("failed marshaling response: %v", err)
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed writing to response: %v", err)
	}
	return nil
}
