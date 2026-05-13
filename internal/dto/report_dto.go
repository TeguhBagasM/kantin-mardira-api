package dto

import "kantin-mardira-api/internal/entity"

type DailyReportQuery struct {
	Date string `form:"date" binding:"required"`
}

type WeeklyReportQuery struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type MonthlyReportQuery struct {
	Month int `form:"month" binding:"required,min=1,max=12"`
	Year  int `form:"year" binding:"required"`
}

type SummaryReportQuery struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type TopSellingQuery struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type DailyReportResponse struct {
	Date              string `json:"date"`
	TotalTransactions int64  `json:"total_transactions"`
	TotalRevenue      int64  `json:"total_revenue"`
	TotalItemsSold    int64  `json:"total_items_sold"`
}

type WeeklyReportResponse struct {
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	TotalTransactions int64  `json:"total_transactions"`
	TotalRevenue      int64  `json:"total_revenue"`
	TotalItemsSold    int64  `json:"total_items_sold"`
}

type MonthlyReportResponse struct {
	Month             int   `json:"month"`
	Year              int   `json:"year"`
	TotalTransactions int64 `json:"total_transactions"`
	TotalRevenue      int64 `json:"total_revenue"`
	TotalItemsSold    int64 `json:"total_items_sold"`
}

type SummaryReportResponse struct {
	TotalTransactions  int64 `json:"total_transactions"`
	TotalRevenue       int64 `json:"total_revenue"`
	AverageTransaction int64 `json:"average_transaction"`
	TotalItemsSold     int64 `json:"total_items_sold"`
}

type TopSellingMenuResponse struct {
	MenuID            string `json:"menu_id"`
	MenuName          string `json:"menu_name"`
	TotalQuantitySold int64  `json:"total_quantity_sold"`
	TotalRevenue      int64  `json:"total_revenue"`
}

type ReportTransactionRow struct {
	ID              string                      `json:"id"`
	TransactionCode string                      `json:"transaction_code"`
	CustomerName    *string                     `json:"customer_name,omitempty"`
	Cashier         *CashierResponseMini        `json:"cashier,omitempty"`
	PaymentMethod   string                      `json:"payment_method"`
	PaymentStatus   string                      `json:"payment_status"`
	TotalAmount     int                         `json:"total_amount"`
	PaidAmount      int                         `json:"paid_amount"`
	ChangeAmount    int                         `json:"change_amount"`
	TransactionTime string                      `json:"transaction_time,omitempty"`
	Items           []TransactionItemResponse   `json:"items,omitempty"`
}

type ReportPDFData struct {
	Title             string
	Period            string
	GeneratedAt       string
	SummaryTitle      string
	SummaryLabels     map[string]string
	SummaryValues     map[string]string
	Transactions      []ReportTransactionRow
	TopSellingMenus   []TopSellingMenuResponse
	EmptyMessage      string
	ShowTopSelling    bool
	ShowTransactions  bool
	DateLabel         string
}

type InvoicePDFData struct {
	TransactionCode string
	CustomerName    *string
	CashierName     string
	TransactionDate string
	PaymentMethod   string
	PaymentStatus   string
	Items           []TransactionItemResponse
	TotalAmount     int
	PaidAmount      int
	ChangeAmount    int
	GeneratedAt     string
	Title           string
	HeaderName      string
	AddressLine     string
	Note            string
	Transaction     *entity.Transaction
}