package router

import (
	"context"
	"encoding/json"
	"github.com/megorka/goproject/authorization/internal/models"
	"github.com/megorka/goproject/authorization/internal/oauth"
	"github.com/megorka/goproject/authorization/internal/service"
	"github.com/megorka/goproject/authorization/pkg/auth"
	"github.com/megorka/goproject/authorization/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req *models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if req.Email == "" || hashedPassword == "" || req.Name == "" || req.LastName == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := h.service.CreateUser(ctx, req.Name, req.LastName, req.Email, hashedPassword); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	logger := logger.GetLoggerFromCtx(r.Context())
	logger.Info(r.Context(), "user created", zap.String("email", req.Email))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User successfuly created"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		http.Error(w, "Invalid name or password", http.StatusUnauthorized)
		return
	}

	logger := logger.GetLoggerFromCtx(r.Context())
	logger.Info(r.Context(), "user logged in", zap.String("email", req.Email))

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth.GoogleOauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := h.service.HandleGoogleCallback(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Work"))
}
