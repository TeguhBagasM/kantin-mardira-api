package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) Daily(c *gin.Context) {
	var query dto.DailyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.reportService.Daily(query.Date)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch daily report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Daily report fetched successfully", "data": response})
}

func (h *ReportHandler) Weekly(c *gin.Context) {
	var query dto.WeeklyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.reportService.Weekly(query.StartDate, query.EndDate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch weekly report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Weekly report fetched successfully", "data": response})
}

func (h *ReportHandler) Monthly(c *gin.Context) {
	var query dto.MonthlyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.reportService.Monthly(query.Month, query.Year)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch monthly report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Monthly report fetched successfully", "data": response})
}

func (h *ReportHandler) Summary(c *gin.Context) {
	var query dto.SummaryReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.reportService.Summary(query.StartDate, query.EndDate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch summary report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Summary report fetched successfully", "data": response})
}

func (h *ReportHandler) TopSelling(c *gin.Context) {
	var query dto.TopSellingQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	response, err := h.reportService.TopSelling(query.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch top selling menus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Top selling menus fetched successfully", "data": response})
}