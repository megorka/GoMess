package service

import (
	"context"
	"fmt"
	websocket2 "github.com/gorilla/websocket"
	"github.com/megorka/goproject/chat_service/internal/repository"
	websocket "github.com/megorka/goproject/chat_service/internal/transport/websocket"
)

type Service struct {
	repo    *repository.Repository
	connMgr *websocket.ConnectionManager
}

func NewService(repo *repository.Repository, connMgr *websocket.ConnectionManager) *Service {
	return &Service{repo: repo, connMgr: connMgr}
}

func (s *Service) SendMessage(ctx context.Context, fromID, toID int, message string) error {
	err := s.repo.SendMessage(ctx, fromID, toID, message)
	if err != nil {
		return fmt.Errorf("SendMessage: %w", err)
	}

	conn, exists := s.connMgr.GetConnection(toID)
	if !exists {
		return nil
	}

	err = conn.WriteMessage(websocket2.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("SendMessage: %w", err)
	}

	err = s.repo.UpdateMessageStatus(ctx, fromID, toID, "delivered")
	if err != nil {
		return fmt.Errorf("UpdateMessageStatus: %w", err)
	}

	return nil
}

func (s *Service) GetMessages(ctx context.Context, fromID, toID int) ([]string, error) {
	return s.repo.GetMessages(ctx, fromID, toID)
}

func (s *Service) AddConnection(ctx context.Context, userID int, conn *websocket2.Conn) {
	s.connMgr.AddConnection(userID, conn)
}

func (s *Service) RemoveConnection(ctx context.Context, userID int) {
	s.connMgr.RemoveConnection(userID)
}

func (s *Service) MarkMessagesAsRead(ctx context.Context, userID int) error {
	err := s.repo.UpdateMessageStatusForUser(ctx, userID, "read")
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}
	return nil
}

func (s *Service) GetUnreadMessages(ctx context.Context, userID int) ([]string, error) {
	return s.repo.GetUnreadMessages(ctx, userID)
}
