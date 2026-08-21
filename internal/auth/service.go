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
}

func NewService(repo *Repository, googleProvider *oauth.GoogleProvider) *Service {
	return &Service{
		repo:           repo,
		googleProvider: googleProvider,
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

	return nil, errors.New("TBD: not implemented")
}
