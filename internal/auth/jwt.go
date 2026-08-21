package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: secret,
	}
}

func (j *JWTService) GenerateAccessToken(
	user *User,
) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,

		"iat": time.Now().Unix(),
		"exp": time.Now().
			Add(15 * time.Minute).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(j.secret),
	)
}

func (j *JWTService) GenerateRefreshToken(user *User) (string, error) {
	// TBD: implement JWT generation logic here
	return "dummy_refresh_token", nil
}

func (j *JWTService) VerifyAccessToken(token string) (*User, error) {
	// TBD: implement JWT verification logic here
	return nil, nil
}

func (j *JWTService) AccessTokenExpiryInSec() int {
	return 15 * 60
}
