package service

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/entity"
	"kantin-mardira-api/internal/repository"
	"kantin-mardira-api/internal/utils"
)

var (
	ErrMenuNotFound          = errors.New("menu not found")
	ErrMenuInvalid           = errors.New("invalid menu")
	ErrMenuCategoryNotFound = errors.New("category not found")
	ErrMenuImageRequired     = errors.New("menu image required")
	ErrMenuImageInvalid      = errors.New("invalid menu image")
)

type MenuService interface {
	Create(request dto.CreateMenuRequest, imageFile *multipart.FileHeader) (*dto.MenuResponse, error)
	FindAll() ([]dto.MenuResponse, error)
	FindByID(id string) (*dto.MenuResponse, error)
	Update(id string, request dto.UpdateMenuRequest, imageFile *multipart.FileHeader) (*dto.MenuResponse, error)
	Delete(id string) error
}

type menuService struct {
	menuRepo     repository.MenuRepository
	publicBaseURL string
}

func NewMenuService(menuRepo repository.MenuRepository, publicBaseURL string) MenuService {
	return &menuService{menuRepo: menuRepo, publicBaseURL: publicBaseURL}
}

func mapMenuResponse(menu *entity.Menu, publicBaseURL string) *dto.MenuResponse {
	response := &dto.MenuResponse{
		ID:          menu.ID.String(),
		Name:        menu.Name,
		Price:       menu.Price,
		Stock:       menu.Stock,
		ImageURL:    resolveMenuImageURL(publicBaseURL, menu.ImageURL),
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

func resolveMenuImageURL(publicBaseURL string, imageURL *string) *string {
	if imageURL == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*imageURL)
	if trimmed == "" {
		return nil
	}

	resolved := utils.ResolvePublicURL(publicBaseURL, trimmed)
	return &resolved
}

func normalizeMenuName(name string) string {
	return strings.TrimSpace(name)
}

func (s *menuService) Create(request dto.CreateMenuRequest, imageFile *multipart.FileHeader) (*dto.MenuResponse, error) {
	name := normalizeMenuName(request.Name)
	if name == "" || request.Price < 0 || request.Stock < 0 {
		return nil, ErrMenuInvalid
	}
	if imageFile == nil {
		return nil, ErrMenuImageRequired
	}

	category, err := s.menuRepo.FindCategoryByID(request.CategoryID)
	if err != nil {
		return nil, ErrMenuCategoryNotFound
	}

	categoryID, _ := uuid.Parse(request.CategoryID)
	imagePath, err := utils.SaveMenuImage(imageFile)
	if err != nil {
		return nil, ErrMenuImageInvalid
	}

	menu := &entity.Menu{
		ID:          uuid.New(),
		CategoryID:  &categoryID,
		Category:    category,
		Name:        name,
		Price:       request.Price,
		Stock:       request.Stock,
		ImageURL:    &imagePath,
		IsAvailable: request.IsAvailable,
	}

	if err := s.menuRepo.Create(menu); err != nil {
		_ = utils.DeleteUploadedFile(imagePath)
		return nil, err
	}

	return mapMenuResponse(menu, s.publicBaseURL), nil
}

func (s *menuService) FindAll() ([]dto.MenuResponse, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MenuResponse, 0, len(menus))
	for i := range menus {
		responses = append(responses, *mapMenuResponse(&menus[i], s.publicBaseURL))
	}

	return responses, nil
}

func (s *menuService) FindByID(id string) (*dto.MenuResponse, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	return mapMenuResponse(menu, s.publicBaseURL), nil
}

func (s *menuService) Update(id string, request dto.UpdateMenuRequest, imageFile *multipart.FileHeader) (*dto.MenuResponse, error) {
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

	var newImagePath string
	oldImagePath := ""
	if menu.ImageURL != nil {
		oldImagePath = *menu.ImageURL
	}

	if imageFile != nil {
		newImagePath, err = utils.SaveMenuImage(imageFile)
		if err != nil {
			return nil, ErrMenuImageInvalid
		}
		menu.ImageURL = &newImagePath
	}

	categoryID, _ := uuid.Parse(request.CategoryID)
	menu.CategoryID = &categoryID
	menu.Category = category
	menu.Name = name
	menu.Price = request.Price
	menu.Stock = request.Stock
	menu.IsAvailable = request.IsAvailable

	if err := s.menuRepo.Update(menu); err != nil {
		if newImagePath != "" {
			_ = utils.DeleteUploadedFile(newImagePath)
		}
		return nil, err
	}

	if newImagePath != "" && oldImagePath != "" && oldImagePath != newImagePath {
		_ = utils.DeleteUploadedFile(oldImagePath)
	}

	return mapMenuResponse(menu, s.publicBaseURL), nil
}

func (s *menuService) Delete(id string) error {
	if _, err := s.menuRepo.FindByID(id); err != nil {
		return ErrMenuNotFound
	}

	return s.menuRepo.Delete(id)
}