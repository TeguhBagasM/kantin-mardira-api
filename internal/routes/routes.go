package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
	"kantin-mardira-api/internal/repository"
	"kantin-mardira-api/internal/service"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	apiV1 := router.Group("/api/v1")

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	authService := service.NewAuthService(userRepo, tokenRepo)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	RegisterAuthRoutes(apiV1, authHandler, authMiddleware)
	RegisterUserRoutes(apiV1, userHandler, authMiddleware)
}