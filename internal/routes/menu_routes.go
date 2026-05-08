package routes

import (
	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
)

func RegisterMenuRoutes(group *gin.RouterGroup, menuHandler *handler.MenuHandler, authMiddleware gin.HandlerFunc) {
	menus := group.Group("/menus")
	menus.Use(authMiddleware)
	{
		menus.GET("", menuHandler.FindAll)
		menus.GET("/:id", menuHandler.FindByID)

		adminOnly := menus.Group("")
		adminOnly.Use(middleware.RoleMiddleware("admin"))
		{
			adminOnly.POST("", menuHandler.Create)
			adminOnly.PUT("/:id", menuHandler.Update)
			adminOnly.DELETE("/:id", menuHandler.Delete)
		}
	}
}