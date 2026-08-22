package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/roymwxuk/mwx-go-auth-service/internal/oauth"
)

type Service struct {
	repo           *Repository
	googleProvider *oauth.GoogleProvider // inject
	jwtService     *JWTService           // inject
}

func NewService(repo *Repository, googleProvider *oauth.GoogleProvider, jwtService *JWTService) *Service {
	return &Service{
		repo:           repo,
		googleProvider: googleProvider,
		jwtService:     jwtService,
	}
}

func (s *Service) LoginWithGoogle(
	ctx context.Context,
	idToken string,
) (*LoginResult, error) {
	googleUser, err := s.googleProvider.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	fmt.Printf("googleUserId: %s, username: %s\n", googleUser.ProviderUserID, googleUser.DisplayName)

	// search the user+userIdentity
	var user *User
	user, err = s.repo.GetUserByIdentity(ctx, googleUser.Provider, googleUser.ProviderUserID)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			fmt.Printf("NOT user not found error, check plz.")
			return nil, err
		}
		// create user if user is not existing
		user, err = s.repo.CreateUserWithIdentity(ctx, &NewUserParam{
			Email:          googleUser.Email,
			AvatarURL:      googleUser.AvatarURL,
			DisplayName:    googleUser.DisplayName,
			Provider:       googleUser.Provider,
			ProviderUserID: googleUser.ProviderUserID,
		})
		if err != nil {
			fmt.Printf("error creating user")
			return nil, err
		}
	}
	fmt.Printf("Found user:%v\n", user)

	// generate JWT

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtService.AccessTokenExpiryInSec(),
	}, nil
}

func (s *Service) GetUserByID(
	ctx context.Context,
	userID string,
) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (*LoginResult, error) {
	// Verify JWT
	claim, err := s.jwtService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check refresh token record
	tokenRecord, err := s.repo.GetRefreshTokenByJTI(
		ctx,
		claim.JTI,
	)
	if err != nil {
		return nil, err
	}

	if tokenRecord.RevokedAt != nil {
		return nil, ErrRefreshTokenRevoked
	}

	// Get user
	user, err := s.repo.GetUserByID(
		ctx,
		claim.UserID,
	)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	newAccessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken, // reuse existing refresh token
		ExpiresIn:    s.jwtService.AccessTokenExpiryInSec(),
	}, nil
}
