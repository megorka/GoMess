package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/megorka/goproject/authorization/internal/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, lastname, email, password, provider FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Password, &user.Provider)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("chat not found")
	} else if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}
	return &user, nil
}

func (r *Repository) FindByProviderID(provider, providerID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, lastname, email FROM users WHERE provider_id = $1`
	err := r.db.QueryRow(context.Background(), query, providerID).Scan(&user.ID, &user.Name, &user.LastName, &user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("FindByProviderID: %w", err)
	}
	return &user, nil
}

func (r *Repository) CreateOAuthUser(name, lastname, email, provider, providerID string) error {
	query := `INSERT INTO users (name, lastname, email, provider, provider_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query, name, lastname, email, provider, providerID)
	if err != nil {
		return fmt.Errorf("CreateOAuthUser: %w", err)
	}
	return nil
}
