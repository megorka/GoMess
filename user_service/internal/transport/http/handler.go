package router

import (
	"encoding/json"
	"github.com/megorka/goproject/user_service/internal/models"
	"github.com/megorka/goproject/user_service/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateFriend(w http.ResponseWriter, r *http.Request) {
	var req *models.Friends

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FriendID == 0 || req.UserID == 0 {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateFriend(r.Context(), req.UserID, req.FriendID); err != nil {
		http.Error(w, "Failed to create friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Friend successfuly created"})
}
