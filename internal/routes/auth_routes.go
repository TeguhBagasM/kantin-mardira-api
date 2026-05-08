package routes

import (
	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/handler"
)

func RegisterAuthRoutes(group *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware gin.HandlerFunc) {
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
	}

	protectedGroup := group.Group("")
	protectedGroup.Use(authMiddleware)
	{
		protectedGroup.GET("/profile", authHandler.Profile)
		protectedGroup.POST("/logout", authHandler.Logout)
	}
}