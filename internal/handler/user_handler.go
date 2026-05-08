package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.userService.Create(request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExist):
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already exists"})
		case errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid role"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User created successfully", "data": response})
}

func (h *UserHandler) FindAll(c *gin.Context) {
	responses, err := h.userService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Users fetched successfully", "data": responses})
}

func (h *UserHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.userService.FindByID(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User fetched successfully", "data": response})
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	response, err := h.userService.Update(id, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		case errors.Is(err, service.ErrEmailAlreadyExist):
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already exists"})
		case errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid role"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update user"})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User updated successfully", "data": response})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.Delete(id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User deleted successfully", "data": gin.H{}})
}