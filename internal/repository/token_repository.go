package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/entity"
)

type TokenRepository interface {
	IsRevoked(jti string) (bool, error)
	RevokeToken(token *entity.RevokedToken) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) IsRevoked(jti string) (bool, error) {
	var revoked entity.RevokedToken
	err := r.db.Where("jti = ? AND expires_at > ?", jti, time.Now()).First(&revoked).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *tokenRepository) RevokeToken(token *entity.RevokedToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}

	return r.db.Create(token).Error
}