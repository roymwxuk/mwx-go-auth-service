package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const prefix = "Bearer "

func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing access token",
			})
			return
		}

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
