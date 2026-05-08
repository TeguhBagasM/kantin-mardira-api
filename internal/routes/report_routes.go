package routes

import (
	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/handler"
	"kantin-mardira-api/internal/middleware"
)

func RegisterReportRoutes(group *gin.RouterGroup, reportHandler *handler.ReportHandler, authMiddleware gin.HandlerFunc) {
	reports := group.Group("/reports")
	reports.Use(authMiddleware)
	reports.Use(middleware.RoleMiddleware("admin", "manager"))
	{
		reports.GET("/daily", reportHandler.Daily)
		reports.GET("/weekly", reportHandler.Weekly)
		reports.GET("/monthly", reportHandler.Monthly)
		reports.GET("/summary", reportHandler.Summary)
		reports.GET("/top-selling", reportHandler.TopSelling)
	}
}