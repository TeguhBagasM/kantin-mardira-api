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
	categoryRepo := repository.NewCategoryRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	reportRepo := repository.NewReportRepository(db)
	authService := service.NewAuthService(userRepo, tokenRepo)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	RegisterAuthRoutes(apiV1, authHandler, authMiddleware)
	RegisterUserRoutes(apiV1, userHandler, authMiddleware)
	RegisterCategoryRoutes(apiV1, categoryHandler, authMiddleware)
	RegisterMenuRoutes(apiV1, menuHandler, authMiddleware)
	RegisterTransactionRoutes(apiV1, transactionHandler, authMiddleware)
	RegisterReportRoutes(apiV1, reportHandler, authMiddleware)
}