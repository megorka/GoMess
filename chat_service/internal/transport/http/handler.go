package router

import (
	"encoding/json"
	"fmt"
	"github.com/megorka/goproject/chat_service/internal/models"
	"github.com/megorka/goproject/chat_service/internal/service"
	"github.com/megorka/goproject/chat_service/pkg/logger"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	from := r.Context().Value("user_id").(string)
	id, err := strconv.Atoi(from)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req *models.Message

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ToID == 0 || req.Content == "" {
		fmt.Println(req.ToID, req.Content)
		logger.GetLoggerFromCtx(r.Context()).Info(r.Context(), "All fields required")
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if err := h.service.SendMessage(r.Context(), id, req.ToID, req.Content); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message successfuly send"})

}

func (h *Handler) GetUnreadMessages(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(string)
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := h.service.GetUnreadMessages(r.Context(), id)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) MarkMessagesAsRead(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(string)
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.MarkMessagesAsRead(r.Context(), id); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Messages marked as read"})
}
