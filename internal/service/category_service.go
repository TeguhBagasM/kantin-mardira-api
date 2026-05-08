package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/entity"
	"kantin-mardira-api/internal/repository"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrInvalidCategory  = errors.New("invalid category")
)

type CategoryService interface {
	Create(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	FindAll() ([]dto.CategoryResponse, error)
	FindByID(id string) (*dto.CategoryResponse, error)
	Update(id string, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(id string) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func mapCategoryResponse(category *entity.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:   category.ID.String(),
		Name: category.Name,
	}
}

func normalizeCategoryName(name string) string {
	return strings.TrimSpace(name)
}

func (s *categoryService) Create(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	name := normalizeCategoryName(request.Name)
	if name == "" {
		return nil, ErrInvalidCategory
	}

	if existingCategory, err := s.categoryRepo.FindByName(name); err == nil && existingCategory != nil {
		return nil, ErrCategoryExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &entity.Category{
		ID:   uuid.New(),
		Name: name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return mapCategoryResponse(category), nil
}

func (s *categoryService) FindAll() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))
	for i := range categories {
		responses = append(responses, dto.CategoryResponse{
			ID:   categories[i].ID.String(),
			Name: categories[i].Name,
		})
	}

	return responses, nil
}

func (s *categoryService) FindByID(id string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	return mapCategoryResponse(category), nil
}

func (s *categoryService) Update(id string, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	name := normalizeCategoryName(request.Name)
	if name == "" {
		return nil, ErrInvalidCategory
	}

	if existingCategory, err := s.categoryRepo.FindByName(name); err == nil && existingCategory != nil && existingCategory.ID != category.ID {
		return nil, ErrCategoryExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category.Name = name
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return mapCategoryResponse(category), nil
}

func (s *categoryService) Delete(id string) error {
	if _, err := s.categoryRepo.FindByID(id); err != nil {
		return ErrCategoryNotFound
	}

	return s.categoryRepo.Delete(id)
}