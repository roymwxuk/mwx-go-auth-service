package auth

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: secret,
	}
}

func (j *JWTService) GenerateAccessToken(user *User) (string, error) {
	// TBD: implement JWT generation logic here
	return "dummy_access_token", nil
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
