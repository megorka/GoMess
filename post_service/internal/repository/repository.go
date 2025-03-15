package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/megorka/goproject/post_service/internal/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePost(ctx context.Context, userId int, title, content string) error {
	query := `INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, userId, title, content)
	if err != nil {
		return fmt.Errorf("CreatePost: %w", err)
	}
	return nil
}

func (r *Repository) UpdatePost(ctx context.Context, postId int, title, content string) error {
	query := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`

	_, err := r.db.Exec(ctx, query, title, content, postId)
	if err != nil {
		return fmt.Errorf("UpdatePost: %w", err)
	}
	return nil
}

func (r *Repository) DeletePost(ctx context.Context, postId int) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := r.db.Exec(ctx, query, postId)
	if err != nil {
		return fmt.Errorf("DeletePost: %w", err)
	}
	return nil
}

func (r *Repository) GetPostById(ctx context.Context, postId int) (models.Post, error) {
	query := `SELECT id, title, content, user_id FROM posts WHERE id = $1`

	var post models.Post
	if err := r.db.QueryRow(ctx, query, postId).Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
		return models.Post{}, fmt.Errorf("GetPostById: %w", err)
	}
	return post, nil
}
