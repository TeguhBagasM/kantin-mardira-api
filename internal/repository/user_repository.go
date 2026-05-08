package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindAll() ([]entity.User, error)
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := r.db.First(&user, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.Delete(&entity.User{}, "id = ?", parsedID).Error
}