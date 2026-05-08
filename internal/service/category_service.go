package service

import (
	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/repository"
)

type CategoryService interface {
	Create(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &dto.CategoryResponse{
		Name: request.Name,
}
	return category, nil
}