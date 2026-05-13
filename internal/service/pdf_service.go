package service

import (
	"fmt"
	"strings"
	"time"

	"kantin-mardira-api/internal/utils"
)

type PDFService interface {
	DailyReportPDF(date string) ([]byte, string, error)
	WeeklyReportPDF(startDate, endDate string) ([]byte, string, error)
	MonthlyReportPDF(month, year int) ([]byte, string, error)
	InvoicePDF(userID, role, transactionID string) ([]byte, string, error)
}

type pdfService struct {
	reportService     ReportService
	transactionService TransactionService
}

func NewPDFService(reportService ReportService, transactionService TransactionService) PDFService {
	return &pdfService{reportService: reportService, transactionService: transactionService}
}

func (s *pdfService) DailyReportPDF(date string) ([]byte, string, error) {
	data, err := s.reportService.Daily(date)
	if err != nil {
		return nil, "", err
	}
	transactions, err := s.reportService.TransactionsBetween(date, date)
	if err != nil {
		return nil, "", err
	}
	generator := utils.NewPDFGenerator("Daily Report")
	generator.AddTitle("Daily Report")
	generator.AddSubtitle(fmt.Sprintf("Period: %s", data.Date))
	generator.AddSubtitle(fmt.Sprintf("Generated at: %s", time.Now().Format("2006-01-02 15:04:05")))
	generator.AddSummarySection(
		map[string]string{"date": "Date", "total_transactions": "Total Transactions", "total_revenue": "Total Revenue", "total_items_sold": "Total Items Sold"},
		map[string]string{"date": data.Date, "total_transactions": fmt.Sprintf("%d", data.TotalTransactions), "total_revenue": generator.AddCurrency(int(data.TotalRevenue)), "total_items_sold": fmt.Sprintf("%d", data.TotalItemsSold)},
	)
	rows := make([][]string, 0, len(transactions))
	for _, tx := range transactions {
		cashierName := "-"
		if tx.Cashier != nil {
			cashierName = tx.Cashier.Name
		}
		customerName := "-"
		if tx.CustomerName != nil && strings.TrimSpace(*tx.CustomerName) != "" {
			customerName = *tx.CustomerName
		}
		rows = append(rows, []string{tx.TransactionCode, customerName, cashierName, tx.TransactionTime, generator.AddCurrency(tx.TotalAmount)})
	}
	generator.AddTitle("Transactions")
	if len(rows) == 0 {
		generator.AddSubtitle("No transactions found for this period.")
	} else {
		generator.AddSimpleTable([]string{"Code", "Customer", "Cashier", "Time", "Total"}, rows, []float64{40, 40, 40, 35, 35})
	}
	bytes, err := generator.OutputBytes()
	if err != nil {
		return nil, "", err
	}
	return bytes, fmt.Sprintf("daily-report-%s.pdf", data.Date), nil
}

func (s *pdfService) WeeklyReportPDF(startDate, endDate string) ([]byte, string, error) {
	data, err := s.reportService.Weekly(startDate, endDate)
	if err != nil {
		return nil, "", err
	}
	transactions, err := s.reportService.TransactionsBetween(startDate, endDate)
	if err != nil {
		return nil, "", err
	}
	generator := utils.NewPDFGenerator("Weekly Report")
	generator.AddTitle("Weekly Report")
	generator.AddSubtitle(fmt.Sprintf("Period: %s to %s", data.StartDate, data.EndDate))
	generator.AddSubtitle(fmt.Sprintf("Generated at: %s", time.Now().Format("2006-01-02 15:04:05")))
	generator.AddSummarySection(
		map[string]string{"total_transactions": "Total Transactions", "total_revenue": "Total Revenue", "total_items_sold": "Total Items Sold"},
		map[string]string{"total_transactions": fmt.Sprintf("%d", data.TotalTransactions), "total_revenue": generator.AddCurrency(int(data.TotalRevenue)), "total_items_sold": fmt.Sprintf("%d", data.TotalItemsSold)},
	)
	if len(transactions) == 0 {
		generator.AddSubtitle("No transactions found for this period.")
	} else {
		rows := make([][]string, 0, len(transactions))
		for _, tx := range transactions {
			cashierName := "-"
			if tx.Cashier != nil {
				cashierName = tx.Cashier.Name
			}
			customerName := "-"
			if tx.CustomerName != nil && strings.TrimSpace(*tx.CustomerName) != "" {
				customerName = *tx.CustomerName
			}
			rows = append(rows, []string{tx.TransactionCode, customerName, cashierName, tx.TransactionTime, generator.AddCurrency(tx.TotalAmount)})
		}
		generator.AddSimpleTable([]string{"Code", "Customer", "Cashier", "Time", "Total"}, rows, []float64{40, 40, 40, 35, 35})
	}
	bytes, err := generator.OutputBytes()
	if err != nil {
		return nil, "", err
	}
	return bytes, fmt.Sprintf("weekly-report-%s-to-%s.pdf", data.StartDate, data.EndDate), nil
}

func (s *pdfService) MonthlyReportPDF(month, year int) ([]byte, string, error) {
	data, err := s.reportService.Monthly(month, year)
	if err != nil {
		return nil, "", err
	}
	monthLabel := fmt.Sprintf("%04d-%02d", year, month)
	generator := utils.NewPDFGenerator("Monthly Report")
	generator.AddTitle("Monthly Report")
	generator.AddSubtitle(fmt.Sprintf("Period: %s", monthLabel))
	generator.AddSubtitle(fmt.Sprintf("Generated at: %s", time.Now().Format("2006-01-02 15:04:05")))
	generator.AddSummarySection(
		map[string]string{"month": "Month", "total_transactions": "Total Transactions", "total_revenue": "Total Revenue", "total_items_sold": "Total Items Sold"},
		map[string]string{"month": monthLabel, "total_transactions": fmt.Sprintf("%d", data.TotalTransactions), "total_revenue": generator.AddCurrency(int(data.TotalRevenue)), "total_items_sold": fmt.Sprintf("%d", data.TotalItemsSold)},
	)
	transactions, err := s.reportService.TransactionsBetween(fmt.Sprintf("%04d-%02d-01", year, month), time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1).Format("2006-01-02"))
	if err == nil && len(transactions) > 0 {
		rows := make([][]string, 0, len(transactions))
		for _, tx := range transactions {
			cashierName := "-"
			if tx.Cashier != nil {
				cashierName = tx.Cashier.Name
			}
			customerName := "-"
			if tx.CustomerName != nil && strings.TrimSpace(*tx.CustomerName) != "" {
				customerName = *tx.CustomerName
			}
			rows = append(rows, []string{tx.TransactionCode, customerName, cashierName, tx.TransactionTime, generator.AddCurrency(tx.TotalAmount)})
		}
		generator.AddTitle("Transactions")
		generator.AddSimpleTable([]string{"Code", "Customer", "Cashier", "Time", "Total"}, rows, []float64{40, 40, 40, 35, 35})
	}
	topMenus, err := s.reportService.TopSelling(10)
	if err == nil && len(topMenus) > 0 {
		generator.AddTitle("Top Selling Menu")
		rows := make([][]string, 0, len(topMenus))
		for _, menu := range topMenus {
			rows = append(rows, []string{menu.MenuName, fmt.Sprintf("%d", menu.TotalQuantitySold), generator.AddCurrency(int(menu.TotalRevenue))})
		}
		generator.AddSimpleTable([]string{"Menu", "Qty Sold", "Revenue"}, rows, []float64{95, 35, 50})
	}
	bytes, err := generator.OutputBytes()
	if err != nil {
		return nil, "", err
	}
	return bytes, fmt.Sprintf("monthly-report-%04d-%02d.pdf", year, month), nil
}

func (s *pdfService) InvoicePDF(userID, role, transactionID string) ([]byte, string, error) {
	transaction, err := s.transactionService.FindByID(userID, role, transactionID)
	if err != nil {
		return nil, "", err
	}
	generator := utils.NewPDFGenerator("Invoice")
	generator.AddTitle("Invoice")
	generator.AddSubtitle("KANTIN MARDIRA")
	generator.AddSubtitle(fmt.Sprintf("Invoice: %s", transaction.TransactionCode))
	customerName := "-"
	if transaction.CustomerName != nil && strings.TrimSpace(*transaction.CustomerName) != "" {
		customerName = *transaction.CustomerName
	}
	generator.AddSubtitle(fmt.Sprintf("Customer: %s", customerName))
	cashierName := "-"
	if transaction.Cashier != nil && strings.TrimSpace(transaction.Cashier.Name) != "" {
		cashierName = transaction.Cashier.Name
	}
	generator.AddSubtitle(fmt.Sprintf("Cashier: %s", cashierName))
	generator.AddSubtitle(fmt.Sprintf("Date: %s", transaction.TransactionTime))
	generator.AddSubtitle(fmt.Sprintf("Payment Method: %s", transaction.PaymentMethod))
	generator.AddSubtitle(fmt.Sprintf("Payment Status: %s", transaction.PaymentStatus))
	rows := make([][]string, 0, len(transaction.Items))
	for _, item := range transaction.Items {
		menuName := "-"
		if item.Menu != nil {
			menuName = item.Menu.Name
		}
		rows = append(rows, []string{menuName, fmt.Sprintf("%d", item.Quantity), generator.AddCurrency(item.Price), generator.AddCurrency(item.Subtotal)})
	}
	generator.AddInvoiceTable(rows)
	generator.AddSummarySection(
		map[string]string{"total": "Total Amount", "paid": "Paid Amount", "change": "Change Amount"},
		map[string]string{"total": generator.AddCurrency(transaction.TotalAmount), "paid": generator.AddCurrency(transaction.PaidAmount), "change": generator.AddCurrency(transaction.ChangeAmount)},
	)
	bytes, err := generator.OutputBytes()
	if err != nil {
		return nil, "", err
	}
	return bytes, fmt.Sprintf("invoice-%s.pdf", transaction.TransactionCode), nil
}