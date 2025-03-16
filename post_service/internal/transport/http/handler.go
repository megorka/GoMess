package router

import (
	"encoding/json"
	"github.com/megorka/goproject/post_service/internal/models"
	"github.com/megorka/goproject/post_service/internal/service"
	"github.com/megorka/goproject/post_service/pkg/logger"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	id, err := strconv.Atoi(userId)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req *models.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Content == "" || req.Title == "" {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "All fields required")
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePost(r.Context(), id, req.Title, req.Content); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post successfuly created"})
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var req *models.Post
	userID := r.Context().Value("user_id").(string)
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" || req.Title == "" {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "All fields required")
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	usersID, err := h.service.GetPostById(r.Context(), req.ID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to update chat", http.StatusInternalServerError)
		return
	}

	if usersID.UserID != id {
		logger.GetLoggerFromCtx(r.Context()).Info(r.Context(), "You can't update this chat")
		http.Error(w, "You can't update this chat", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePost(r.Context(), id, req.ID, req.Title, req.Content); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to update chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post successfuly updated"})
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	id, err := strconv.Atoi(userID)
	var req struct {
		ID int `json:"id" db:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	usersID, err := h.service.GetPostById(r.Context(), req.ID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to update chat", http.StatusInternalServerError)
		return
	}

	if usersID.UserID != id {
		logger.GetLoggerFromCtx(r.Context()).Info(r.Context(), "You can't update this chat")
		http.Error(w, "You can't update this chat", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePost(r.Context(), req.ID); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to delete chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post successfuly deleted"})
}
