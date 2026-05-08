package routes

import (
	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
)

func RegisterPDFRoutes(group *gin.RouterGroup, pdfHandler *handler.PDFHandler, authMiddleware gin.HandlerFunc) {
	reports := group.Group("/reports")
	reports.Use(authMiddleware)
	reports.Use(middleware.RoleMiddleware("admin", "manager"))
	{
		reports.GET("/daily/pdf", pdfHandler.Daily)
		reports.GET("/weekly/pdf", pdfHandler.Weekly)
		reports.GET("/monthly/pdf", pdfHandler.Monthly)
	}

	transactions := group.Group("/transactions")
	transactions.Use(authMiddleware)
	{
		transactions.GET("/:id/invoice", pdfHandler.Invoice)
	}
}