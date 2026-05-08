package routes

import (
	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
)

func RegisterTransactionRoutes(group *gin.RouterGroup, transactionHandler *handler.TransactionHandler, authMiddleware gin.HandlerFunc) {
	transactions := group.Group("/transactions")
	transactions.Use(authMiddleware)
	{
		adminCashier := transactions.Group("")
		adminCashier.Use(middleware.RoleMiddleware("admin", "cashier"))
		{
			adminCashier.POST("", transactionHandler.Create)
		}

		readAll := transactions.Group("")
		readAll.Use(middleware.RoleMiddleware("admin", "manager", "cashier"))
		{
			readAll.GET("", transactionHandler.FindAll)
			readAll.GET("/:id", transactionHandler.FindByID)
		}
	}
}