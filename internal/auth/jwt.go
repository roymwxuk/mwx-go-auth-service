package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signingMethod = jwt.SigningMethodHS256

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
		signingMethod,
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
		signingMethod,
		claims,
	)

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type AccessTokenClaim struct {
	UserID string
	Email  string
}

func (j *JWTService) VerifyAccessToken(
	tokenString string,
) (*AccessTokenClaim, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if token.Method != signingMethod {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(j.secret), nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("missing user_id")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("missing email")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	return &AccessTokenClaim{
		UserID: userID,
		Email:  email,
	}, nil
}

func (j *JWTService) AccessTokenExpiryInSec() int {
	return 15 * 60
}
