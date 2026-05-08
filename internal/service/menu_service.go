package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/entity"
	"kantin-mardira-api/internal/repository"
)

var (
	ErrMenuNotFound      = errors.New("menu not found")
	ErrMenuInvalid       = errors.New("invalid menu")
	ErrMenuCategoryNotFound = errors.New("category not found")
)

type MenuService interface {
	Create(request dto.CreateMenuRequest) (*dto.MenuResponse, error)
	FindAll() ([]dto.MenuResponse, error)
	FindByID(id string) (*dto.MenuResponse, error)
	Update(id string, request dto.UpdateMenuRequest) (*dto.MenuResponse, error)
	Delete(id string) error
}

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{menuRepo: menuRepo}
}

func mapMenuResponse(menu *entity.Menu) *dto.MenuResponse {
	response := &dto.MenuResponse{
		ID:          menu.ID.String(),
		Name:        menu.Name,
		Price:       menu.Price,
		Stock:       menu.Stock,
		ImageURL:    menu.ImageURL,
		IsAvailable: menu.IsAvailable,
	}

	if menu.Category != nil {
		response.Category = &dto.CategoryResponseMini{
			ID:   menu.Category.ID.String(),
			Name: menu.Category.Name,
		}
	}

	return response
}

func normalizeMenuName(name string) string {
	return strings.TrimSpace(name)
}

func (s *menuService) Create(request dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	name := normalizeMenuName(request.Name)
	if name == "" || request.Price < 0 || request.Stock < 0 {
		return nil, ErrMenuInvalid
	}

	category, err := s.menuRepo.FindCategoryByID(request.CategoryID)
	if err != nil {
		return nil, ErrMenuCategoryNotFound
	}

	categoryID, _ := uuid.Parse(request.CategoryID)
	menu := &entity.Menu{
		ID:          uuid.New(),
		CategoryID:  &categoryID,
		Category:    category,
		Name:        name,
		Price:       request.Price,
		Stock:       request.Stock,
		ImageURL:    request.ImageURL,
		IsAvailable: request.IsAvailable,
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}

	return mapMenuResponse(menu), nil
}

func (s *menuService) FindAll() ([]dto.MenuResponse, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MenuResponse, 0, len(menus))
	for i := range menus {
		responses = append(responses, *mapMenuResponse(&menus[i]))
	}

	return responses, nil
}

func (s *menuService) FindByID(id string) (*dto.MenuResponse, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	return mapMenuResponse(menu), nil
}

func (s *menuService) Update(id string, request dto.UpdateMenuRequest) (*dto.MenuResponse, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	name := normalizeMenuName(request.Name)
	if name == "" || request.Price < 0 || request.Stock < 0 {
		return nil, ErrMenuInvalid
	}

	category, err := s.menuRepo.FindCategoryByID(request.CategoryID)
	if err != nil {
		return nil, ErrMenuCategoryNotFound
	}

	categoryID, _ := uuid.Parse(request.CategoryID)
	menu.CategoryID = &categoryID
	menu.Category = category
	menu.Name = name
	menu.Price = request.Price
	menu.Stock = request.Stock
	menu.ImageURL = request.ImageURL
	menu.IsAvailable = request.IsAvailable

	if err := s.menuRepo.Update(menu); err != nil {
		return nil, err
	}

	return mapMenuResponse(menu), nil
}

func (s *menuService) Delete(id string) error {
	if _, err := s.menuRepo.FindByID(id); err != nil {
		return ErrMenuNotFound
	}

	return s.menuRepo.Delete(id)
}