package handler

import (
	"log"
	"net/http"

	"example.com/myproject/service"
)

type UserHandler struct {
	s *service.UserService
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		s: s,
	}
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var regReq service.RegistrationRequest
	err := parseRequest(r, &regReq)
	if err != nil {
		log.Printf("Failed registering user %q: %v", regReq.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.s.RegisterUser(ctx, regReq)
	if err != nil {
		if err == service.ErrUserExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		log.Printf("Failed registering user %q: %v", regReq.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	resp := &TokenResponse{
		AccessToken: token,
	}
	err = writeResponse(w, resp)
	if err != nil {
		log.Printf("Failed registering user %q: %v", regReq.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var logReq service.LoginRequest
	err := parseRequest(r, &logReq)
	if err != nil {
		log.Printf("Failed parsing login request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.s.LoginUser(ctx, logReq)
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

	resp := &TokenResponse{
		AccessToken: token,
	}
	err = writeResponse(w, resp)
	if err != nil {
		log.Printf("Failed writing login response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
