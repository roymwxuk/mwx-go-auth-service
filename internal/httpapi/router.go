package httpapi

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/roymwxuk/mwx-go-auth-service/internal/auth"
)

func NewRouter(authHandler *auth.Handler) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           60 * time.Second,
	}))

	api := r.Group("/api")

	api.GET("/health", HealthHandler)

	// auth
	api.POST("/auth/google", authHandler.LoginWithGoogle)

	// users
	api.GET("/users/me", authHandler.GetMe)

	return r
}
