package routes

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
	"kantin-mardira-api/internal/repository"
	"kantin-mardira-api/internal/service"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	router.Static("/uploads", "./uploads")

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
	menuService := service.NewMenuService(menuRepo, resolvePublicBaseURL())
	menuHandler := handler.NewMenuHandler(menuService)
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)
	pdfService := service.NewPDFService(reportService, transactionService)
	pdfHandler := handler.NewPDFHandler(pdfService)
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	RegisterAuthRoutes(apiV1, authHandler, authMiddleware)
	RegisterUserRoutes(apiV1, userHandler, authMiddleware)
	RegisterCategoryRoutes(apiV1, categoryHandler, authMiddleware)
	RegisterMenuRoutes(apiV1, menuHandler, authMiddleware)
	RegisterTransactionRoutes(apiV1, transactionHandler, authMiddleware)
	RegisterReportRoutes(apiV1, reportHandler, authMiddleware)
	RegisterPDFRoutes(apiV1, pdfHandler, authMiddleware)
}

func resolvePublicBaseURL() string {
	if value := strings.TrimSpace(os.Getenv("APP_PUBLIC_URL")); value != "" {
		return strings.TrimRight(value, "/")
	}

	port := strings.TrimSpace(os.Getenv("APP_PORT"))
	if port == "" {
		port = "8080"
	}

	return fmt.Sprintf("http://localhost:%s", port)
}