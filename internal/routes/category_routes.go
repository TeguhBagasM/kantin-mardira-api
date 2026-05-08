package routes

import (
	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(group *gin.RouterGroup, categoryHandler *handler.CategoryHandler, authMiddleware gin.HandlerFunc) {
	categories := group.Group("/categories")
	categories.Use(authMiddleware)
	{
		categories.GET("", categoryHandler.FindAll)
		categories.GET("/:id", categoryHandler.FindByID)

		adminOnly := categories.Group("")
		adminOnly.Use(middleware.RoleMiddleware("admin"))
		{
			adminOnly.POST("", categoryHandler.Create)
			adminOnly.PUT("/:id", categoryHandler.Update)
			adminOnly.DELETE("/:id", categoryHandler.Delete)
		}
	}
}