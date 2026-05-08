package repository

import (
	"time"

	"gorm.io/gorm"
)

type ReportAggregateRow struct {
	TotalTransactions int64
	TotalRevenue      int64
	TotalItemsSold    int64
}

type TopSellingMenuRow struct {
	MenuID           string
	MenuName         string
	TotalQuantitySold int64
	TotalRevenue      int64
}

type ReportRepository interface {
	DailyReport(date time.Time) (*ReportAggregateRow, error)
	WeeklyReport(startDate, endDate time.Time) (*ReportAggregateRow, error)
	MonthlyReport(month int, year int) (*ReportAggregateRow, error)
	SummaryReport(startDate, endDate time.Time) (*ReportAggregateRow, error)
	TopSellingMenus(limit int) ([]TopSellingMenuRow, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) aggregateBetween(startDate, endDate time.Time) (*ReportAggregateRow, error) {
	var row ReportAggregateRow
	if err := r.db.Table("transactions").
		Select(`
			COALESCE(COUNT(DISTINCT transactions.id), 0) AS total_transactions,
			COALESCE(SUM(transaction_items.subtotal), 0) AS total_revenue,
			COALESCE(SUM(transaction_items.quantity), 0) AS total_items_sold
		`).
		Joins("JOIN transaction_items ON transaction_items.transaction_id = transactions.id").
		Where("transactions.payment_status = ? AND transactions.transaction_time >= ? AND transactions.transaction_time < ?", "paid", startDate, endDate).
		Scan(&row).Error; err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *reportRepository) DailyReport(date time.Time) (*ReportAggregateRow, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	return r.aggregateBetween(start, end)
}

func (r *reportRepository) WeeklyReport(startDate, endDate time.Time) (*ReportAggregateRow, error) {
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location()).Add(24 * time.Hour)
	return r.aggregateBetween(start, end)
}

func (r *reportRepository) MonthlyReport(month int, year int) (*ReportAggregateRow, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	return r.aggregateBetween(start, end)
}

func (r *reportRepository) SummaryReport(startDate, endDate time.Time) (*ReportAggregateRow, error) {
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location()).Add(24 * time.Hour)
	return r.aggregateBetween(start, end)
}

func (r *reportRepository) TopSellingMenus(limit int) ([]TopSellingMenuRow, error) {
	var rows []TopSellingMenuRow
	if err := r.db.Table("transaction_items").
		Select(`
			menus.id AS menu_id,
			menus.name AS menu_name,
			COALESCE(SUM(transaction_items.quantity), 0) AS total_quantity_sold,
			COALESCE(SUM(transaction_items.subtotal), 0) AS total_revenue
		`).
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("JOIN menus ON menus.id = transaction_items.menu_id").
		Where("transactions.payment_status = ?", "paid").
		Group("menus.id, menus.name").
		Order("total_quantity_sold DESC, total_revenue DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}