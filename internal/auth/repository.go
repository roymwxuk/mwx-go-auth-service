package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

type NewUserParam struct {
	DisplayName string
	Email       string
	AvatarURL   string

	// for user_identities
	Provider       string
	ProviderUserID string
}

func (r *Repository) CreateUserWithIdentity(
	ctx context.Context,
	param *NewUserParam,
) (*User, error) {
	fmt.Printf("creating a new user with email: %s\n", param.Email)
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var user User

	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO users (
			email,
			display_name,
			avatar_url
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			email,
			display_name,
			avatar_url,
			created_at,
			updated_at
		`,
		param.Email,
		param.DisplayName,
		param.AvatarURL,
	).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO user_identities (
			user_id,
			provider,
			provider_user_id
		)
		VALUES ($1, $2, $3)
		`,
		user.ID,
		param.Provider,
		param.ProviderUserID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
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
		SELECT
			u.id,
			u.email,
			u.display_name,
			u.avatar_url
		FROM user_identities ui
		JOIN users u
			ON u.id = ui.user_id
		WHERE ui.provider = $1
		  AND ui.provider_user_id = $2
		`,
		provider,
		providerUserID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
