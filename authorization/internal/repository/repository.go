package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/megorka/goproject/authorization/internal/models"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, name, lastname, email, password string) error {
	query := `INSERT INTO users (name, lastname, email, password) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, name, lastname, email, password)
	if err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	return nil
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, lastname, email, password FROM users WHERE email = $1`
	err := r.db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}
	return &user, nil
}
