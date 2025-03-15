package service

import (
	"context"
	"github.com/megorka/goproject/post_service/internal/models"
	"github.com/megorka/goproject/post_service/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePost(ctx context.Context, userId int, title, content string) error {
	return s.repo.CreatePost(ctx, userId, title, content)
}

func (s *Service) UpdatePost(ctx context.Context, id, postId int, title, content string) error {
	return s.repo.UpdatePost(ctx, postId, title, content)
}

func (s *Service) DeletePost(ctx context.Context, postId int) error {
	return s.repo.DeletePost(ctx, postId)
}

func (s *Service) GetPostById(ctx context.Context, postId int) (models.Post, error) {
	return s.repo.GetPostById(ctx, postId)
}
