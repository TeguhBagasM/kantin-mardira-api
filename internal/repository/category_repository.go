package repository

import (
	"kantin-mardira-api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) categoryRepository {
	return categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *entity.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return r.db.Create(category).Error
}