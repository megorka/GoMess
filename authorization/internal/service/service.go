package service

import (
	"context"
	"github.com/megorka/goproject/authorization/internal/models"
	"github.com/megorka/goproject/authorization/internal/repository"
	"github.com/megorka/goproject/authorization/pkg/auth"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, name, lastname, email, password string) error {
	return s.repo.CreateUser(ctx, name, lastname, email, password)
}

func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", nil
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
