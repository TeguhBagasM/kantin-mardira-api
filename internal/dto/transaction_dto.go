package dto

type TransactionItemRequest struct {
	MenuID   string `json:"menu_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gte=1"`
}

type CreateTransactionRequest struct {
	CustomerName  string                   `json:"customer_name" binding:"omitempty,max=100"`
	PaymentMethod string                   `json:"payment_method" binding:"required,oneof=cash qris"`
	PaidAmount    int                      `json:"paid_amount"`
	Items         []TransactionItemRequest `json:"items" binding:"required,min=1"`
}

type MenuResponseMini struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransactionItemResponse struct {
	Menu     *MenuResponseMini `json:"menu,omitempty"`
	Quantity int               `json:"quantity"`
	Price    int               `json:"price"`
	Subtotal int               `json:"subtotal"`
}

type CashierResponseMini struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransactionResponse struct {
	ID              string               `json:"id"`
	TransactionCode string               `json:"transaction_code"`
	CustomerName    *string              `json:"customer_name,omitempty"`
	Cashier         *CashierResponseMini `json:"cashier,omitempty"`
	PaymentMethod   string               `json:"payment_method"`
	PaymentStatus   string               `json:"payment_status"`
	TotalAmount     int                  `json:"total_amount"`
	PaidAmount      int                  `json:"paid_amount"`
	ChangeAmount    int                  `json:"change_amount"`
	TransactionTime string               `json:"transaction_time,omitempty"`
}

type TransactionDetailResponse struct {
	ID              string                    `json:"id"`
	TransactionCode string                    `json:"transaction_code"`
	CustomerName    *string                   `json:"customer_name,omitempty"`
	Cashier         *CashierResponseMini      `json:"cashier,omitempty"`
	PaymentMethod   string                    `json:"payment_method"`
	PaymentStatus   string                    `json:"payment_status"`
	TotalAmount     int                       `json:"total_amount"`
	PaidAmount      int                       `json:"paid_amount"`
	ChangeAmount    int                       `json:"change_amount"`
	TransactionTime string                    `json:"transaction_time,omitempty"`
	Items           []TransactionItemResponse `json:"items"`
}