package service

import (
	"errors"
	"time"

	"kantin-mardira-api/internal/dto"
	"kantin-mardira-api/internal/repository"
)

var ErrInvalidDateFormat = errors.New("invalid date format")

type ReportService interface {
	Daily(date string) (*dto.DailyReportResponse, error)
	Weekly(startDate, endDate string) (*dto.WeeklyReportResponse, error)
	Monthly(month int, year int) (*dto.MonthlyReportResponse, error)
	Summary(startDate, endDate string) (*dto.SummaryReportResponse, error)
	TopSelling(limit int) ([]dto.TopSellingMenuResponse, error)
	TransactionsBetween(startDate, endDate string) ([]dto.ReportTransactionRow, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func (s *reportService) Daily(date string) (*dto.DailyReportResponse, error) {
	parsedDate, err := parseDate(date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	row, err := s.reportRepo.DailyReport(parsedDate)
	if err != nil {
		return nil, err
	}

	return &dto.DailyReportResponse{
		Date:              parsedDate.Format("2006-01-02"),
		TotalTransactions: row.TotalTransactions,
		TotalRevenue:      row.TotalRevenue,
		TotalItemsSold:    row.TotalItemsSold,
	}, nil
}

func (s *reportService) Weekly(startDate, endDate string) (*dto.WeeklyReportResponse, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	end, err := parseDate(endDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if end.Before(start) {
		return nil, ErrInvalidDateFormat
	}

	row, err := s.reportRepo.WeeklyReport(start, end)
	if err != nil {
		return nil, err
	}

	return &dto.WeeklyReportResponse{
		StartDate:         start.Format("2006-01-02"),
		EndDate:           end.Format("2006-01-02"),
		TotalTransactions: row.TotalTransactions,
		TotalRevenue:      row.TotalRevenue,
		TotalItemsSold:    row.TotalItemsSold,
	}, nil
}

func (s *reportService) Monthly(month int, year int) (*dto.MonthlyReportResponse, error) {
	if month < 1 || month > 12 || year < 2000 {
		return nil, ErrInvalidDateFormat
	}

	row, err := s.reportRepo.MonthlyReport(month, year)
	if err != nil {
		return nil, err
	}

	return &dto.MonthlyReportResponse{
		Month:             month,
		Year:              year,
		TotalTransactions: row.TotalTransactions,
		TotalRevenue:      row.TotalRevenue,
		TotalItemsSold:    row.TotalItemsSold,
	}, nil
}

func (s *reportService) Summary(startDate, endDate string) (*dto.SummaryReportResponse, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	end, err := parseDate(endDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if end.Before(start) {
		return nil, ErrInvalidDateFormat
	}

	row, err := s.reportRepo.SummaryReport(start, end)
	if err != nil {
		return nil, err
	}

	average := int64(0)
	if row.TotalTransactions > 0 {
		average = row.TotalRevenue / row.TotalTransactions
	}

	return &dto.SummaryReportResponse{
		TotalTransactions:  row.TotalTransactions,
		TotalRevenue:      row.TotalRevenue,
		AverageTransaction: average,
		TotalItemsSold:    row.TotalItemsSold,
	}, nil
}

func (s *reportService) TopSelling(limit int) ([]dto.TopSellingMenuResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.reportRepo.TopSellingMenus(limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TopSellingMenuResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, dto.TopSellingMenuResponse{
			MenuID:           row.MenuID,
			MenuName:         row.MenuName,
			TotalQuantitySold: row.TotalQuantitySold,
			TotalRevenue:      row.TotalRevenue,
		})
	}

	return responses, nil
}

func (s *reportService) TransactionsBetween(startDate, endDate string) ([]dto.ReportTransactionRow, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	end, err := parseDate(endDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if end.Before(start) {
		return nil, ErrInvalidDateFormat
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location()).Add(24 * time.Hour)

	rows, err := s.reportRepo.TransactionsBetween(start, end)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReportTransactionRow, 0, len(rows))
	for _, row := range rows {
		response := dto.ReportTransactionRow{
			ID:              row.ID,
			TransactionCode: row.TransactionCode,
			CustomerName:    row.CustomerName,
			PaymentMethod:   row.PaymentMethod,
			PaymentStatus:   row.PaymentStatus,
			TotalAmount:     row.TotalAmount,
			PaidAmount:      row.PaidAmount,
			ChangeAmount:    row.ChangeAmount,
			TransactionTime: row.TransactionTime.Format(time.RFC3339),
		}
		if row.CashierID != "" || row.CashierName != "" {
			response.Cashier = &dto.CashierResponseMini{ID: row.CashierID, Name: row.CashierName}
		}
		responses = append(responses, response)
	}

	return responses, nil
}