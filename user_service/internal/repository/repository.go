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

func (r *Repository) CreateFriends(ctx context.Context, userID, friendID int) error {
	query := `INSERT INTO friends (friend1_id, friend2_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, userID, friendID)
	if err != nil {
		return fmt.Errorf("CreateFriends: %w", err)
	}
	return nil
}
