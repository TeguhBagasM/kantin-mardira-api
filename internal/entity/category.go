package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID			uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name		string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt	time.Time `gorm:"autoCreateTime" json:"created_at"`
}