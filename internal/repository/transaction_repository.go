package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/entity"
)

type TransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) error
	CreateTransactionItem(tx *gorm.DB, item *entity.TransactionItem) error
	FindAll() ([]entity.Transaction, error)
	FindAllByCashierID(cashierID string) ([]entity.Transaction, error)
	FindByID(id string) (*entity.Transaction, error)
	FindByIDByCashierID(id, cashierID string) (*entity.Transaction, error)
	FindMenuByID(menuID string) (*entity.Menu, error)
	UpdateMenuStock(tx *gorm.DB, menuID string, stock int) error
	GenerateTransactionCode() (string, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) error {
	if transaction.ID == uuid.Nil {
		transaction.ID = uuid.New()
	}
	return tx.Create(transaction).Error
}

func (r *transactionRepository) CreateTransactionItem(tx *gorm.DB, item *entity.TransactionItem) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return tx.Create(item).Error
}

func (r *transactionRepository) FindAll() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Preload("Cashier").Preload("Items.Menu").Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) FindAllByCashierID(cashierID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Preload("Cashier").Preload("Items.Menu").Where("cashier_id = ?", cashierID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) FindByID(id string) (*entity.Transaction, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var transaction entity.Transaction
	if err := r.db.Preload("Cashier").Preload("Items.Menu").First(&transaction, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByIDByCashierID(id, cashierID string) (*entity.Transaction, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var transaction entity.Transaction
	if err := r.db.Preload("Cashier").Preload("Items.Menu").Where("id = ? AND cashier_id = ?", parsedID, cashierID).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) FindMenuByID(menuID string) (*entity.Menu, error) {
	parsedID, err := uuid.Parse(menuID)
	if err != nil {
		return nil, err
	}

	var menu entity.Menu
	if err := r.db.First(&menu, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}

	return &menu, nil
}

func (r *transactionRepository) UpdateMenuStock(tx *gorm.DB, menuID string, stock int) error {
	parsedID, err := uuid.Parse(menuID)
	if err != nil {
		return err
	}

	return tx.Model(&entity.Menu{}).Where("id = ?", parsedID).Update("stock", stock).Error
}

func (r *transactionRepository) GenerateTransactionCode() (string, error) {
	today := time.Now().Format("20060102")
	prefix := fmt.Sprintf("TRX-%s-", today)

	var count int64
	if err := r.db.Model(&entity.Transaction{}).Where("transaction_code LIKE ?", prefix+"%").Count(&count).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("TRX-%s-%03d", today, count+1), nil
}