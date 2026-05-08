package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var request dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.transactionService.Create(userID.(string), request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientStock):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Insufficient stock"})
		case errors.Is(err, service.ErrMenuNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Menu not found"})
		case errors.Is(err, service.ErrMenuCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Category not found"})
		case errors.Is(err, service.ErrInvalidPaymentAmount):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid payment amount"})
		case errors.Is(err, service.ErrInvalidPaymentMethod):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid payment method"})
		case errors.Is(err, service.ErrInvalidTransactionItem):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid transaction item"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create transaction"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Transaction created successfully", "data": response})
}

func (h *TransactionHandler) FindAll(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	responses, err := h.transactionService.FindAll(userID.(string), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transactions fetched successfully", "data": responses})
}

func (h *TransactionHandler) FindByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	id := c.Param("id")
	response, err := h.transactionService.FindByID(userID.(string), role.(string), id)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Transaction not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid transaction ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transaction fetched successfully", "data": response})
}