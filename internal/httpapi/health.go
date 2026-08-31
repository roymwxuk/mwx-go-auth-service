package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "mw-auth-service is running. version: 1.0.1",
	})
}
