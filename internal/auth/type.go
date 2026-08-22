package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	DisplayName string
	AvatarURL   *string // for handling null values
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserIdentity struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	JTI       string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

var ErrRefreshTokenRevoked = errors.New("refresh token revoked")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}
