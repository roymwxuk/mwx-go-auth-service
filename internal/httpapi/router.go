package httpapi

import (
	"github.com/gin-gonic/gin"

	"github.com/roymwxuk/mwx-go-auth-service/internal/auth"
)

func NewRouter(authHandler *auth.Handler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	api.GET("/health", HealthHandler)

	// auth
	api.POST("/auth/google", authHandler.LoginWithGoogle)

	// users
	api.GET("/users/me", authHandler.GetMe)

	return r
}
