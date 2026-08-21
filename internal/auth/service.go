package auth

import (
	"context"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

func (s *Service) LoginWithGoogle(
	ctx context.Context,
	idToken string,
) (*LoginResult, error) {
	return nil, errors.New("TBD: not implemented")
}
