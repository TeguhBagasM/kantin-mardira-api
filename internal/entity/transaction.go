package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TransactionCode   string              `gorm:"type:varchar(50);uniqueIndex;not null" json:"transaction_code"`
	CustomerName      *string             `gorm:"type:varchar(100)" json:"customer_name,omitempty"`
	CashierID         *uuid.UUID          `gorm:"type:uuid" json:"cashier_id"`
	Cashier           *User               `gorm:"foreignKey:CashierID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"cashier,omitempty"`
	PaymentMethod     string              `gorm:"type:varchar(20);not null" json:"payment_method"`
	PaymentStatus     string              `gorm:"type:varchar(20);not null" json:"payment_status"`
	TotalAmount       int                 `gorm:"type:int;not null;check:total_amount >= 0" json:"total_amount"`
	PaidAmount        int                 `gorm:"type:int;default:0" json:"paid_amount"`
	ChangeAmount      int                 `gorm:"type:int;default:0" json:"change_amount"`
	TransactionTime   time.Time           `gorm:"autoCreateTime" json:"transaction_time"`
	CreatedAt         time.Time           `gorm:"autoCreateTime" json:"created_at"`
	Items             []TransactionItem   `gorm:"foreignKey:TransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}