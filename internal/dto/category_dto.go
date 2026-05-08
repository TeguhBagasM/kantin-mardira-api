package dto

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}