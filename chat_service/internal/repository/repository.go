package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SendMessage(ctx context.Context, fromID, toID int, message string) error {
	query := `INSERT INTO direct_messages (sender_id, receiver_id, content, status) VALUES ($1, $2, $3, 'sent')`

	_, err := r.db.Exec(ctx, query, fromID, toID, message)
	if err != nil {
		return fmt.Errorf("SendMessage: %w", err)
	}
	return nil
}

func (r *Repository) GetMessages(ctx context.Context, fromID, toID int) ([]string, error) {
	query := `SELECT content FROM direct_messages WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) ORDER BY id`
	rows, err := r.db.Query(ctx, query, fromID, toID)
	if err != nil {
		return nil, fmt.Errorf("GetMessages: %w", err)
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, fmt.Errorf("GetMessages: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *Repository) UpdateMessageStatus(ctx context.Context, fromID, toID int, status string) error {
	query := `UPDATE direct_messages SET status = $1 WHERE (sender_id = $2 AND receiver_id = $3) OR (sender_id = $3 AND receiver_id = $2)`
	_, err := r.db.Exec(ctx, query, status, fromID, toID)
	if err != nil {
		return fmt.Errorf("UpdateMessageStatus: %w", err)
	}
	return nil
}

func (r *Repository) UpdateMessageStatusForUser(ctx context.Context, toID int, status string) error {
	query := "UPDATE direct_messages SET status = $1 WHERE receiver_id = $2 AND status != $3"
	_, err := r.db.Exec(ctx, query, status, toID, "read")
	if err != nil {
		return fmt.Errorf("failed to update message status for user: %w", err)
	}
	return nil
}

func (r *Repository) GetUnreadMessages(ctx context.Context, userID int) ([]string, error) {
	query := "SELECT content FROM direct_messages WHERE receiver_id = $1 AND status != 'read'"

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread messages: %w", err)
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, fmt.Errorf("failed to get unread messages: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}
