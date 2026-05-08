package entity

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CategoryID   *uuid.UUID `gorm:"type:uuid" json:"category_id"`
	Category     *Category  `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Price       int        `gorm:"type:int;not null;check:price >= 0" json:"price"`
	Stock       int        `gorm:"type:int;default:0;check:stock >= 0" json:"stock"`
	ImageURL    *string    `gorm:"type:text" json:"image_url,omitempty"`
	IsAvailable bool       `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Menu) TableName() string {
	return "menus"
}