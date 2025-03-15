package service

import (
	"context"
	"fmt"
	"github.com/megorka/goproject/authorization/internal/models"
	"github.com/megorka/goproject/authorization/internal/oauth"
	"github.com/megorka/goproject/authorization/internal/repository"
	"github.com/megorka/goproject/authorization/pkg/auth"
	"golang.org/x/oauth2"
	oauthv2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type Service struct {
	repo     *repository.Repository
	oauthcfg *oauth2.Config
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo, oauthcfg: oauth.GoogleOauthConfig}
}

func (s *Service) CreateUser(ctx context.Context, name, lastname, email, password string) error {
	return s.repo.CreateUser(ctx, name, lastname, email, password)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user.Provider != "local" {
		return "", fmt.Errorf("post registered via %s, use OAuth login", user.Provider)
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}
	return token, nil
}

func (s *Service) HandleGoogleCallback(code string) (string, error) {
	token, err := s.oauthcfg.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	oauth2Service, err := oauthv2.NewService(context.Background(), option.WithTokenSource(s.oauthcfg.TokenSource(context.Background(), token)))
	if err != nil {
		return "", fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return "", fmt.Errorf("failed to get post info: %w", err)
	}

	user, err := s.repo.FindByProviderID("google", userInfo.Id)
	if err != nil {
		return "", fmt.Errorf("failed to find post by provider id: %w", err)
	}

	if user == nil {
		newUser := &models.User{
			Name:       userInfo.GivenName,
			LastName:   userInfo.FamilyName,
			Email:      userInfo.Email,
			Provider:   "google",
			ProviderID: userInfo.Id,
		}
		if err := s.repo.CreateOAuthUser(newUser.Name, newUser.LastName, newUser.Email, newUser.Provider, newUser.ProviderID); err != nil {
			return "", err
		}
		user = newUser
	}
	return auth.CreateToken(user.ID)
}
