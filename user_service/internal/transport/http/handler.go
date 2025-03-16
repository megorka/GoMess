package router

import (
	"encoding/json"
	"fmt"
	"github.com/megorka/goproject/user_service/internal/service"
	"github.com/megorka/goproject/user_service/pkg/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateFriend(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FriendID int `json:"friend_id"`
	}

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

	if req.FriendID == 0 {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "All fields required")
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateFriend(r.Context(), id, req.FriendID); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to create friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Friend successfuly created"})
}

func (h *Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FriendID int `json:"friend_id"`
	}
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

	if req.FriendID == 0 {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "All fields required")
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFriend(r.Context(), id, req.FriendID); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to delete friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Friend successfuly deleted"})
}

func (h *Handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	friends, err := h.service.GetFriends(r.Context(), id)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to get friends", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}

func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to upload avatar", http.StatusInternalServerError)
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", time.Now().Unix()) + fileExt
	filePath := "https://localhost:8080/images/" + filename

	out, err := os.Create("../../public/single/" + filename)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to upload avatar", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to upload avatar", http.StatusInternalServerError)
		return
	}
	if err := h.service.UploadAvatar(r.Context(), id, filePath); err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to upload avatar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"avatarURL": filePath})
}

func (h *Handler) GetPostsOnProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	posts, err := h.service.GetPostsOnProfile(r.Context(), id)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Failed to get posts on profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
