package utils

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct {
	pdf *gofpdf.Fpdf
}

func NewPDFGenerator(title string) *PDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(title, false)
	pdf.SetAuthor("Kantin Mardira API", false)
	pdf.SetMargins(15, 25, 15)
	pdf.SetAutoPageBreak(true, 20)

	generator := &PDFGenerator{pdf: pdf}
	pdf.SetHeaderFunc(func() {
		generator.header()
	})
	pdf.SetFooterFunc(func() {
		generator.footer()
	})
	pdf.AddPage()
	return generator
}

func (g *PDFGenerator) header() {
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.CellFormat(0, 7, "KANTIN MARDIRA", "", 1, "L", false, 0, "")
	g.pdf.SetFont("Arial", "", 10)
	g.pdf.CellFormat(0, 5, "Jl. Soekarno-Hatta No. 211, Leuwipanjang.", "", 1, "L", false, 0, "")
	g.pdf.Ln(2)
	g.pdf.Line(15, 22, 195, 22)
	g.pdf.Ln(6)
}

func (g *PDFGenerator) footer() {
	g.pdf.SetY(-15)
	g.pdf.SetFont("Arial", "I", 8)
	g.pdf.Line(15, g.pdf.GetY(), 195, g.pdf.GetY())
	g.pdf.CellFormat(0, 5, fmt.Sprintf("Page %d", g.pdf.PageNo()), "", 0, "C", false, 0, "")
}

func (g *PDFGenerator) AddTitle(title string) {
	g.pdf.SetFont("Arial", "B", 16)
	g.pdf.CellFormat(0, 10, strings.ToUpper(title), "", 1, "C", false, 0, "")
	g.pdf.Ln(2)
}

func (g *PDFGenerator) AddSubtitle(text string) {
	g.pdf.SetFont("Arial", "", 10)
	g.pdf.CellFormat(0, 6, text, "", 1, "C", false, 0, "")
	g.pdf.Ln(2)
}

func (g *PDFGenerator) AddSummarySection(labels map[string]string, values map[string]string) {
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.CellFormat(0, 7, "Summary", "", 1, "L", false, 0, "")
	g.pdf.SetFont("Arial", "", 10)
	for key, label := range labels {
		g.pdf.CellFormat(50, 7, label, "1", 0, "L", false, 0, "")
		g.pdf.CellFormat(0, 7, values[key], "1", 1, "L", false, 0, "")
	}
	g.pdf.Ln(2)
}

func (g *PDFGenerator) AddSimpleTable(headers []string, rows [][]string, widths []float64) {
	g.pdf.SetFont("Arial", "B", 10)
	for i, header := range headers {
		g.pdf.CellFormat(widths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	g.pdf.Ln(-1)
	g.pdf.SetFont("Arial", "", 9)
	for _, row := range rows {
		for i, cell := range row {
			align := "L"
			if i > 0 {
				align = "C"
			}
			g.pdf.CellFormat(widths[i], 7, cell, "1", 0, align, false, 0, "")
		}
		g.pdf.Ln(-1)
	}
	g.pdf.Ln(2)
}

func (g *PDFGenerator) AddCurrency(value int) string {
	return fmt.Sprintf("Rp %s", formatNumber(int64(value)))
}

func formatNumber(value int64) string {
	negative := value < 0
	if negative {
		value = -value
	}
	str := fmt.Sprintf("%d", value)
	if len(str) <= 3 {
		if negative {
			return "-" + str
		}
		return str
	}

	var parts []string
	for len(str) > 3 {
		parts = append([]string{str[len(str)-3:]}, parts...)
		str = str[:len(str)-3]
	}
	if str != "" {
		parts = append([]string{str}, parts...)
	}
	result := strings.Join(parts, ".")
	if negative {
		result = "-" + result
	}
	return result
}

func (g *PDFGenerator) OutputBytes() ([]byte, error) {
	var buffer bytes.Buffer
	if err := g.pdf.Output(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (g *PDFGenerator) AddInvoiceTable(rows [][]string) {
	g.AddSimpleTable([]string{"Menu", "Qty", "Price", "Subtotal"}, rows, []float64{80, 20, 45, 45})
}

func (g *PDFGenerator) AddPaginatedNote(totalRows int) {
	g.pdf.SetFont("Arial", "I", 8)
	g.pdf.CellFormat(0, 5, fmt.Sprintf("Total rows: %d", int(math.Abs(float64(totalRows)))), "", 1, "R", false, 0, "")
}