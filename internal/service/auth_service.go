package service

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/entity"
	"kantin-mardira-api/internal/repository"
	"kantin-mardira-api/internal/utils"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	Login(request dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(userID, jti string, expiresAt int64) error
}

type authService struct {
	userRepo repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s *authService) Login(request dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.ComparePassword(user.Password, request.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.LoginUserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *authService) Logout(userID, jti string, expiresAt int64) error {
	if jti == "" {
		return errors.New("invalid token")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	if s.tokenRepo == nil {
		return errors.New("token repository is not configured")
	}

	return s.tokenRepo.RevokeToken(&entity.RevokedToken{
		JTI:       jti,
		UserID:    parsedUserID,
		ExpiresAt: time.Unix(expiresAt, 0).UTC(),
	})
}