package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
)

type MenuHandler struct {
	menuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (h *MenuHandler) Create(c *gin.Context) {
	var request dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.menuService.Create(request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMenuCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Category not found"})
		case errors.Is(err, service.ErrMenuInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid menu"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create menu"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Menu created successfully", "data": response})
}

func (h *MenuHandler) FindAll(c *gin.Context) {
	responses, err := h.menuService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch menus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Menus fetched successfully", "data": responses})
}

func (h *MenuHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.menuService.FindByID(id)
	if err != nil {
		if errors.Is(err, service.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Menu not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid menu ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Menu fetched successfully", "data": response})
}

func (h *MenuHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var request dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.menuService.Update(id, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMenuNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Menu not found"})
		case errors.Is(err, service.ErrMenuCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Category not found"})
		case errors.Is(err, service.ErrMenuInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid menu"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update menu"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Menu updated successfully", "data": response})
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.menuService.Delete(id); err != nil {
		if errors.Is(err, service.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Menu not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid menu ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Menu deleted successfully", "data": gin.H{}})
}