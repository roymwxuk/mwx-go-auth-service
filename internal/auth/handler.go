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

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
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

	res := LoginResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
		ExpiresIn:    loginResult.ExpiresIn,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMe(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented: GetMe",
	})
}
