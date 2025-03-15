package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/megorka/goproject/user_service/internal/models"
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

func (r *Repository) DeleteFriend(ctx context.Context, userID, friendID int) error {
	query := `DELETE FROM friends WHERE (friend1_id = $1 AND friend2_id = $2) OR (friend1_id = $2 AND friend2_id = $1)`
	_, err := r.db.Exec(ctx, query, userID, friendID)
	if err != nil {
		return fmt.Errorf("DeleteFriend: %w", err)
	}
	return nil
}

func (r *Repository) GetFriends(ctx context.Context, userID int) ([]models.User, error) {
	var friends []models.User
	query := `SELECT id, name, lastname, email FROM users WHERE id IN (SELECT friend2_id FROM friends WHERE friend1_id = $1)`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetFriends: %w", err)
	}
	for rows.Next() {
		var friend models.User
		err := rows.Scan(&friend.ID, &friend.Name, &friend.LastName, &friend.Email)
		if err != nil {
			return nil, fmt.Errorf("GetFriends: %w", err)
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *Repository) UploadAvatar(ctx context.Context, userID int, avatarURL string) error {
	query := `UPDATE users SET avatar = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, avatarURL, userID)
	if err != nil {
		return fmt.Errorf("UploadAvatar: %w", err)
	}
	return nil
}

func (r *Repository) GetPostsOnProfile(ctx context.Context, userID int) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetPostsOnProfile: %w", err)
	}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Created_at)
		if err != nil {
			return nil, fmt.Errorf("GetPostsOnProfile: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
