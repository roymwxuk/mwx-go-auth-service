package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const CookieDomain = "roymwxuk.uk"
const SecureAttribute = true

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (h *Handler) LoginWithGoogle(c *gin.Context) {
	ctx := c.Request.Context()

	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	loginResult, err := h.service.LoginWithGoogle(ctx, req.IDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// todo: update for production
	// c.SetSameSite(http.SameSiteLaxMode)
	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		"access_token",
		loginResult.AccessToken,
		loginResult.ExpiresIn,
		"/",
		"roymwxuk.uk",
		SecureAttribute, // Secure
		true,            // HttpOnly
	)

	c.SetCookie(
		"refresh_token",
		loginResult.RefreshToken,
		60*60*24*30,
		"/auth/refresh",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"expires_in": loginResult.ExpiresIn,
	})
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "refresh token is required",
		})
		return
	}

	loginResult, err := h.service.Refresh(
		ctx,
		refreshToken,
	)
	if err != nil {
		log.Printf("RefreshToken error: %v\n", err)

		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "refresh token is invalid",
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"access_token",
		loginResult.AccessToken,
		loginResult.ExpiresIn,
		"/",
		CookieDomain,
		SecureAttribute,
		true,
	)

	res := LoginResponse{
		ExpiresIn: loginResult.ExpiresIn,
		ExpiresAt: loginResult.ExpiresAt,
	}

	c.JSON(http.StatusOK, res)
}
