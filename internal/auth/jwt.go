package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signingMethod = jwt.SigningMethodRS384

const ISSUER = "mwx-go-auth-service"

type JWTService struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJWTService(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *JWTService {
	return &JWTService{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (j *JWTService) GenerateAccessToken(
	user *User,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"type": "access",

		"iss": ISSUER,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().
			Add(15 * time.Minute).
			Unix(),
	}

	token := jwt.NewWithClaims(
		signingMethod,
		claims,
	)

	return token.SignedString(j.privateKey)
}

func (j *JWTService) GenerateRefreshToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"type": "refresh",

		"iss": ISSUER,
		"iat": time.Now().Unix(),
		"exp": time.Now().
			Add(30 * 24 * time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		signingMethod,
		claims,
	)

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type AccessTokenClaim struct {
	UserID string
}

func (j *JWTService) VerifyAccessToken(
	tokenString string,
) (*AccessTokenClaim, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if token.Method != signingMethod {
			return nil, errors.New("unexpected signing method")
		}

		return j.publicKey, nil
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

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("missing sub")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	return &AccessTokenClaim{
		UserID: userID,
	}, nil
}

func (j *JWTService) AccessTokenExpiryInSec() int {
	return 15 * 60
}
