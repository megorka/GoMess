package service

import (
	"context"
	"github.com/megorka/goproject/user_service/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateFriend(ctx context.Context, friend1_id, friend2_id int) error {
	return s.repo.CreateFriends(ctx, friend1_id, friend2_id)
}
