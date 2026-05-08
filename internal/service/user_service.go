package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/entity"
	"kantin-mardira-api/internal/repository"
	"kantin-mardira-api/internal/utils"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrInvalidRole       = errors.New("invalid role")
)

type UserService interface {
	Create(request dto.CreateUserRequest) (*dto.UserResponse, error)
	FindAll() ([]dto.UserResponse, error)
	FindByID(id string) (*dto.UserResponse, error)
	Update(id string, request dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}

func isValidRole(role string) bool {
	switch normalizeRole(role) {
	case "admin", "cashier", "manager":
		return true
	default:
		return false
	}
}

func mapUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func (s *userService) Create(request dto.CreateUserRequest) (*dto.UserResponse, error) {
	if !isValidRole(request.Role) {
		return nil, ErrInvalidRole
	}

	email := normalizeEmail(request.Email)
	if existingUser, err := s.userRepo.FindByEmail(email); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else if existingUser != nil {
		return nil, ErrEmailAlreadyExist
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.New(),
		Name:     strings.TrimSpace(request.Name),
		Email:    email,
		Password: hashedPassword,
		Role:     normalizeRole(request.Role),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return mapUserResponse(user), nil
}

func (s *userService) FindAll() ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		responses = append(responses, dto.UserResponse{
			ID:    users[i].ID.String(),
			Name:  users[i].Name,
			Email: users[i].Email,
			Role:  users[i].Role,
		})
	}

	return responses, nil
}

func (s *userService) FindByID(id string) (*dto.UserResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return mapUserResponse(user), nil
}

func (s *userService) Update(id string, request dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, err
	}

	if !isValidRole(request.Role) {
		return nil, ErrInvalidRole
	}

	currentUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	email := normalizeEmail(request.Email)
	if existingUser, err := s.userRepo.FindByEmail(email); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else if existingUser != nil && existingUser.ID != currentUser.ID {
		return nil, ErrEmailAlreadyExist
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	currentUser.Name = strings.TrimSpace(request.Name)
	currentUser.Email = email
	currentUser.Password = hashedPassword
	currentUser.Role = normalizeRole(request.Role)

	if err := s.userRepo.Update(currentUser); err != nil {
		return nil, err
	}

	return mapUserResponse(currentUser), nil
}

func (s *userService) Delete(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}

	if _, err := s.userRepo.FindByID(id); err != nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(id)
}