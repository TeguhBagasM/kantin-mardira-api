package entity

import (
	"time"

	"github.com/google/uuid"
)

type RevokedToken struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	JTI        string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"jti"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (RevokedToken) TableName() string {
	return "revoked_tokens"
}