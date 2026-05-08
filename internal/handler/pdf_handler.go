package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
)

type PDFHandler struct {
	pdfService service.PDFService
}

func NewPDFHandler(pdfService service.PDFService) *PDFHandler {
	return &PDFHandler{pdfService: pdfService}
}

func sendPDF(c *gin.Context, content []byte, filename string) {
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/pdf", content)
}

func (h *PDFHandler) Daily(c *gin.Context) {
	var query dto.DailyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	content, filename, err := h.pdfService.DailyReportPDF(query.Date)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate daily PDF"})
		return
	}

	sendPDF(c, content, filename)
}

func (h *PDFHandler) Weekly(c *gin.Context) {
	var query dto.WeeklyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	content, filename, err := h.pdfService.WeeklyReportPDF(query.StartDate, query.EndDate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate weekly PDF"})
		return
	}

	sendPDF(c, content, filename)
}

func (h *PDFHandler) Monthly(c *gin.Context) {
	var query dto.MonthlyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	content, filename, err := h.pdfService.MonthlyReportPDF(query.Month, query.Year)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate monthly PDF"})
		return
	}

	sendPDF(c, content, filename)
}

func (h *PDFHandler) Invoice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	transactionID := c.Param("id")
	content, filename, err := h.pdfService.InvoicePDF(userID.(string), role.(string), transactionID)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate invoice PDF"})
		return
	}

	sendPDF(c, content, filename)
}