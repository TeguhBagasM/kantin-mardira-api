package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
	"kantin-mardira-api/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	response, err := h.authService.Login(request)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid credentials",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login success",
		"data":    response,
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "authorized",
		"data": dto.ProfileResponse{
			UserID: userID.(string),
			Role:   role.(string),
		},
	})
}

	func (h *AuthHandler) Logout(c *gin.Context) {
		claimsValue, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			return
		}

		claims, ok := claimsValue.(*utils.Claims)
		if !ok || claims == nil || claims.ExpiresAt == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			return
		}

		if err := h.authService.Logout(claims.UserID, claims.ID, claims.ExpiresAt.Time.Unix()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to logout",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Logout success",
			"data": gin.H{},
		})
	}