package service

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/entity"
	"kantin-mardira-api/internal/repository"
)

var (
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrInvalidPaymentMethod   = errors.New("invalid payment method")
	ErrInvalidPaymentAmount   = errors.New("invalid payment amount")
	ErrInvalidTransactionItem = errors.New("invalid transaction item")
)

type TransactionService interface {
	Create(userID string, request dto.CreateTransactionRequest) (*dto.TransactionDetailResponse, error)
	FindAll(userID, role string) ([]dto.TransactionResponse, error)
	FindByID(userID, role, id string) (*dto.TransactionDetailResponse, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	db              *gorm.DB
}

func NewTransactionService(transactionRepo repository.TransactionRepository, db *gorm.DB) TransactionService {
	return &transactionService{transactionRepo: transactionRepo, db: db}
}

func mapTransactionSummary(transaction *entity.Transaction) *dto.TransactionResponse {
	response := &dto.TransactionResponse{
		ID:              transaction.ID.String(),
		TransactionCode: transaction.TransactionCode,
		CustomerName:    transaction.CustomerName,
		PaymentMethod:   transaction.PaymentMethod,
		PaymentStatus:   transaction.PaymentStatus,
		TotalAmount:     transaction.TotalAmount,
		PaidAmount:      transaction.PaidAmount,
		ChangeAmount:    transaction.ChangeAmount,
		TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
	}

	if transaction.Cashier != nil {
		response.Cashier = &dto.CashierResponseMini{
			ID:   transaction.Cashier.ID.String(),
			Name: transaction.Cashier.Name,
		}
	}

	return response
}

func mapTransactionDetail(transaction *entity.Transaction) *dto.TransactionDetailResponse {
	detail := &dto.TransactionDetailResponse{
		ID:              transaction.ID.String(),
		TransactionCode: transaction.TransactionCode,
		CustomerName:    transaction.CustomerName,
		PaymentMethod:   transaction.PaymentMethod,
		PaymentStatus:   transaction.PaymentStatus,
		TotalAmount:     transaction.TotalAmount,
		PaidAmount:      transaction.PaidAmount,
		ChangeAmount:    transaction.ChangeAmount,
		TransactionTime: transaction.TransactionTime.Format(time.RFC3339),
		Items:           make([]dto.TransactionItemResponse, 0, len(transaction.Items)),
	}

	if transaction.Cashier != nil {
		detail.Cashier = &dto.CashierResponseMini{
			ID:   transaction.Cashier.ID.String(),
			Name: transaction.Cashier.Name,
		}
	}

	for i := range transaction.Items {
		item := transaction.Items[i]
		itemResponse := dto.TransactionItemResponse{
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: item.Subtotal,
		}
		if item.Menu != nil {
			itemResponse.Menu = &dto.MenuResponseMini{
				ID:   item.Menu.ID.String(),
				Name: item.Menu.Name,
			}
		}
		detail.Items = append(detail.Items, itemResponse)
	}

	return detail
}

func parseUUIDFromContext(userID string) (uuid.UUID, error) {
	return uuid.Parse(strings.TrimSpace(userID))
}

func normalizeOptionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func mapCreateTransactionRequestToEntity(transactionCode string, cashierID uuid.UUID, request dto.CreateTransactionRequest) *entity.Transaction {
	return &entity.Transaction{
		ID:             uuid.New(),
		TransactionCode: transactionCode,
		CustomerName:    normalizeOptionalString(request.CustomerName),
		CashierID:       &cashierID,
		PaymentMethod:   strings.ToLower(strings.TrimSpace(request.PaymentMethod)),
		PaymentStatus:   "pending",
	}
}

func (s *transactionService) Create(userID string, request dto.CreateTransactionRequest) (*dto.TransactionDetailResponse, error) {
	if len(request.Items) == 0 {
		return nil, ErrInvalidTransactionItem
	}

	cashierID, err := parseUUIDFromContext(userID)
	if err != nil {
		return nil, err
	}

	transactionCode, err := s.transactionRepo.GenerateTransactionCode()
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	transaction := mapCreateTransactionRequestToEntity(transactionCode, cashierID, request)

	var totalAmount int
	items := make([]entity.TransactionItem, 0, len(request.Items))

	for _, itemRequest := range request.Items {
		if itemRequest.Quantity <= 0 {
			tx.Rollback()
			return nil, ErrInvalidTransactionItem
		}

		menu, err := s.transactionRepo.FindMenuByID(itemRequest.MenuID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if !menu.IsAvailable {
			tx.Rollback()
			return nil, ErrInvalidTransactionItem
		}

		if menu.Stock < itemRequest.Quantity {
			tx.Rollback()
			return nil, ErrInsufficientStock
		}

		subtotal := menu.Price * itemRequest.Quantity
		totalAmount += subtotal

		menuID, _ := uuid.Parse(itemRequest.MenuID)
		items = append(items, entity.TransactionItem{
			ID:            uuid.New(),
			TransactionID: transaction.ID,
			MenuID:        &menuID,
			Quantity:      itemRequest.Quantity,
			Price:         menu.Price,
			Subtotal:      subtotal,
		})

		if err := s.transactionRepo.UpdateMenuStock(tx, itemRequest.MenuID, menu.Stock-itemRequest.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	paymentMethod := strings.ToLower(strings.TrimSpace(request.PaymentMethod))
	paidAmount := request.PaidAmount
	changeAmount := 0
	paymentStatus := "paid"

	switch paymentMethod {
	case "cash":
		if paidAmount < totalAmount {
			tx.Rollback()
			return nil, ErrInvalidPaymentAmount
		}
		changeAmount = paidAmount - totalAmount
	case "qris":
		paidAmount = 0
		changeAmount = 0
	default:
		tx.Rollback()
		return nil, ErrInvalidPaymentMethod
	}

	transaction.PaymentMethod = paymentMethod
	transaction.PaymentStatus = paymentStatus
	transaction.TotalAmount = totalAmount
	transaction.PaidAmount = paidAmount
	transaction.ChangeAmount = changeAmount

	if err := s.transactionRepo.CreateTransaction(tx, transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range items {
		items[i].TransactionID = transaction.ID
		if err := s.transactionRepo.CreateTransactionItem(tx, &items[i]); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	createdTransaction, err := s.transactionRepo.FindByID(transaction.ID.String())
	if err != nil {
		return nil, err
	}

	return mapTransactionDetail(createdTransaction), nil
}

func (s *transactionService) FindAll(userID, role string) ([]dto.TransactionResponse, error) {
	var (
		transactions []entity.Transaction
		err          error
	)

	switch strings.ToLower(strings.TrimSpace(role)) {
	case "cashier":
		transactions, err = s.transactionRepo.FindAllByCashierID(userID)
	default:
		transactions, err = s.transactionRepo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	responses := make([]dto.TransactionResponse, 0, len(transactions))
	for i := range transactions {
		responses = append(responses, *mapTransactionSummary(&transactions[i]))
	}

	return responses, nil
}

func (s *transactionService) FindByID(userID, role, id string) (*dto.TransactionDetailResponse, error) {
	var (
		transaction *entity.Transaction
		err         error
	)

	switch strings.ToLower(strings.TrimSpace(role)) {
	case "cashier":
		transaction, err = s.transactionRepo.FindByIDByCashierID(id, userID)
	default:
		transaction, err = s.transactionRepo.FindByID(id)
	}

	if err != nil {
		return nil, ErrTransactionNotFound
	}

	return mapTransactionDetail(transaction), nil
}