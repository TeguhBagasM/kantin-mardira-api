package entity

import (
	"github.com/google/uuid"
)

type TransactionItem struct {
	ID            uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TransactionID uuid.UUID    `gorm:"type:uuid;not null" json:"transaction_id"`
	Transaction   *Transaction `gorm:"foreignKey:TransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"transaction,omitempty"`
	MenuID        *uuid.UUID   `gorm:"type:uuid" json:"menu_id"`
	Menu          *Menu        `gorm:"foreignKey:MenuID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"menu,omitempty"`
	Quantity      int          `gorm:"type:int;not null;check:quantity > 0" json:"quantity"`
	Price         int          `gorm:"type:int;not null;check:price >= 0" json:"price"`
	Subtotal      int          `gorm:"type:int;not null;check:subtotal >= 0" json:"subtotal"`
}

func (TransactionItem) TableName() string {
	return "transaction_items"
}