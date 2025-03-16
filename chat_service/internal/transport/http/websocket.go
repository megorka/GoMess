package router

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/megorka/goproject/chat_service/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), err.Error())
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "failed to upgrade to WebSocket", zap.Error(err))
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	h.service.AddConnection(r.Context(), id, conn)
	defer h.service.RemoveConnection(r.Context(), id)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "failed to read message", zap.Error(err))
			break
		}

		logger.GetLoggerFromCtx(r.Context()).Info(r.Context(), "Message received", zap.String("message", string(message)))

		var req struct {
			ToID int `json:"to_id"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "failed to unmarshal message", zap.Error(err))
			continue
		}

		if err := h.service.SendMessage(r.Context(), id, req.ToID, string(message)); err != nil {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "failed to send message", zap.Error(err))
			continue
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			logger.GetLoggerFromCtx(r.Context()).Error(r.Context(), "failed to write message", zap.Error(err))
			break
		}
	}
}
