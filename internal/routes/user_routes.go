package routes

import (
	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
)

func RegisterUserRoutes(group *gin.RouterGroup, userHandler *handler.UserHandler, authMiddleware gin.HandlerFunc) {
	users := group.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("", userHandler.FindAll)
		users.GET("/:id", userHandler.FindByID)

		adminOnly := users.Group("")
		adminOnly.Use(middleware.RoleMiddleware("admin"))
		{
			adminOnly.POST("", userHandler.Create)
			adminOnly.PUT("/:id", userHandler.Update)
			adminOnly.DELETE("/:id", userHandler.Delete)
		}
	}
}