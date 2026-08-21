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
		"type":    "access",

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
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"type":    "refresh",

		"iat": time.Now().Unix(),
		"exp": time.Now().
			Add(30 * 24 * time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTService) VerifyAccessToken(token string) (*User, error) {
	// TBD: implement JWT verification logic here
	return nil, nil
}

func (j *JWTService) AccessTokenExpiryInSec() int {
	return 15 * 60
}
