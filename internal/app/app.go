package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roymwxuk/mwx-go-auth-service/config"
	"github.com/roymwxuk/mwx-go-auth-service/internal/auth"
	"github.com/roymwxuk/mwx-go-auth-service/internal/database"
	"github.com/roymwxuk/mwx-go-auth-service/internal/httpapi"
	"github.com/roymwxuk/mwx-go-auth-service/internal/oauth"
)

type App struct {
	server *http.Server
	pool   *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	googleProvider := oauth.NewGoogleProvider(cfg.GoogleClientID)

	privatePEM := strings.ReplaceAll(
		cfg.RsaPrivateKey,
		`\n`,
		"\n",
	)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePEM))
	if err != nil {
		return nil, err
	}

	publicPEM := strings.ReplaceAll(
		cfg.RsaPublicKey,
		`\n`,
		"\n",
	)

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicPEM))
	if err != nil {
		return nil, err
	}

	jwtService := auth.NewJWTService(
		publicKey,
		privateKey,
	)

	authService := auth.NewService(auth.NewRepository(pool), googleProvider, jwtService)
	authHandler := auth.NewHandler(authService)

	router := httpapi.NewRouter(authHandler, jwtService)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &App{
		server: server,
		pool:   pool,
	}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Close() {
	a.pool.Close()
}
