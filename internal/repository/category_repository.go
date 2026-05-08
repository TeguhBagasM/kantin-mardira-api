package repository

import (
	"errors"
	"strings"

	"kantin-mardira-api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
	FindAll() ([]entity.Category, error)
	FindByID(id string) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	Update(category *entity.Category) error
	Delete(id string) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *entity.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.db.Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(id string) (*entity.Category, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var category entity.Category
	if err := r.db.First(&category, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.Where("LOWER(name) = ?", strings.ToLower(strings.TrimSpace(name))).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Model(&entity.Category{}).
		Where("id = ?", category.ID).
		Update("name", category.Name).Error
}

func (r *categoryRepository) Delete(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.Delete(&entity.Category{}, "id = ?", parsedID).Error
}