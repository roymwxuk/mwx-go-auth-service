package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) LoginWithGoogle(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	loginResult, err := h.service.LoginWithGoogle(c.Request.Context(), req.IDToken)
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
		"",
		true, // Secure
		true, // HttpOnly
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
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	loginResult, err := h.service.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "refresh token is invalid",
		})
		return
	}

	res := LoginResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
		ExpiresIn:    loginResult.ExpiresIn,
	}

	c.JSON(http.StatusOK, res)
}
