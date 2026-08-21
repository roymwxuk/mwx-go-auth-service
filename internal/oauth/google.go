package oauth

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

const ProviderName = "google"

type GoogleProvider struct {
	clientID string
}

func NewGoogleProvider(clientID string) *GoogleProvider {
	return &GoogleProvider{
		clientID: clientID,
	}
}

type GoogleUser struct {
	Provider       string
	ProviderUserID string

	Email         string
	EmailVerified bool
	DisplayName   string
	AvatarURL     string
}

func (p *GoogleProvider) VerifyIDToken(ctx context.Context, idToken string) (*GoogleUser, error) {
	payload, err := idtoken.Validate(ctx, idToken, p.clientID)
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("Missing claims.email")
	}
	avatarURL, ok := payload.Claims["picture"].(string)
	if !ok {
		return nil, errors.New("Missing claims.picture")
	}

	return &GoogleUser{
		Provider:       ProviderName,
		ProviderUserID: payload.Subject,
		Email:          email,
		EmailVerified:  payload.Claims["email_verified"].(bool),
		DisplayName:    payload.Claims["name"].(string),
		AvatarURL:      avatarURL,
	}, nil

}
