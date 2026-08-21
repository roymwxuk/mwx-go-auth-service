package app

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roymwxuk/mwx-go-auth-service/config"
	"github.com/roymwxuk/mwx-go-auth-service/internal/auth"
	"github.com/roymwxuk/mwx-go-auth-service/internal/database"
	"github.com/roymwxuk/mwx-go-auth-service/internal/httpapi"
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

	authService := auth.NewService(auth.NewRepository(pool))
	authHandler := auth.NewHandler(authService)

	router := httpapi.NewRouter(authHandler)

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
