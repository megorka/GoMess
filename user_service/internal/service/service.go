package service

import (
	"context"
	"github.com/megorka/goproject/user_service/internal/models"
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

func (s *Service) DeleteFriend(ctx context.Context, friend1_id, friend2_id int) error {
	return s.repo.DeleteFriend(ctx, friend1_id, friend2_id)
}

func (s *Service) GetFriends(ctx context.Context, userID int) ([]models.User, error) {
	return s.repo.GetFriends(ctx, userID)
}

func (s *Service) UploadAvatar(ctx context.Context, userID int, avatarURL string) error {
	return s.repo.UploadAvatar(ctx, userID, avatarURL)
}

func (s *Service) GetPostsOnProfile(ctx context.Context, userID int) ([]models.Post, error) {
	return s.repo.GetPostsOnProfile(ctx, userID)
}
