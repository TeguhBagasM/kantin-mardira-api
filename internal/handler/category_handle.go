package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var request dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	category, err := h.categoryService.Create(request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryExists):
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Category already exists"})
		case errors.Is(err, service.ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid category"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create category"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Category created successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) FindAll(c *gin.Context) {
	categories, err := h.categoryService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Categories fetched successfully",
		"data":    categories,
	})
}

func (h *CategoryHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryService.FindByID(id)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Category not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid category ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category fetched successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var request dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	category, err := h.categoryService.Update(id, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Category not found"})
		case errors.Is(err, service.ErrCategoryExists):
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Category already exists"})
		case errors.Is(err, service.ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid category"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update category"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.categoryService.Delete(id); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Category not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid category ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
		"data":    gin.H{},
	})
}