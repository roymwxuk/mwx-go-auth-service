package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const prefix = "Bearer "

func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid Authorization header",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)

		claim, err := jwtService.VerifyAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid access token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claim.UserID)

		c.Next()
	}
}
