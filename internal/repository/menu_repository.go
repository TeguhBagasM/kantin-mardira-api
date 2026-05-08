package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/entity"
)

type MenuRepository interface {
	Create(menu *entity.Menu) error
	FindAll() ([]entity.Menu, error)
	FindByID(id string) (*entity.Menu, error)
	Update(menu *entity.Menu) error
	Delete(id string) error
	FindCategoryByID(categoryID string) (*entity.Category, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(menu *entity.Menu) error {
	if menu.ID == uuid.Nil {
		menu.ID = uuid.New()
	}
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	if err := r.db.Preload("Category").Order("created_at DESC").Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) FindByID(id string) (*entity.Menu, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var menu entity.Menu
	if err := r.db.Preload("Category").First(&menu, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}

	return &menu, nil
}

func (r *menuRepository) Update(menu *entity.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.Delete(&entity.Menu{}, "id = ?", parsedID).Error
}

func (r *menuRepository) FindCategoryByID(categoryID string) (*entity.Category, error) {
	parsedID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	var category entity.Category
	if err := r.db.First(&category, "id = ?", parsedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, err
	}

	return &category, nil
}