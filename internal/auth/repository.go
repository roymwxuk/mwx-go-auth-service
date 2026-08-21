package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetUserByIdentity(
	ctx context.Context,
	provider string,
	providerUserID string,
) (*User, error) {
	var user User
	err := r.pool.QueryRow(
		ctx,
		`
		SELECT ID, Email, DisplayName, AvatarURL
		FROM users
		WHERE provider=$1 AND provider_user_id = $2
		`,
		provider,
		providerUserID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
