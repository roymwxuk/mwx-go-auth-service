package httpapi

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	healthHandler := NewHealthHandler()

	api.GET("/health", healthHandler.Health)

	return r
}
