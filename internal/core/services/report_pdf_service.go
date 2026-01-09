package services

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/entorno35/backend/internal/domain"
	"github.com/jung-kurt/gofpdf"
)

// ReportPDFService handles professional PDF generation for assessment reports
type ReportPDFService struct{}

// NewReportPDFService creates a new PDF service instance
func NewReportPDFService() *ReportPDFService {
	return &ReportPDFService{}
}

// GenerateIndividualReportPDF creates a high-fidelity PDF report for an individual assessment
func (s *ReportPDFService) GenerateIndividualReportPDF(report *domain.IndividualReportDTO, companyName string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "fonts/")
	pdf.AddPage()

	// Set up fonts with UTF-8 support
	pdf.AddUTF8Font("DejaVuSans", "", "DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVuSans", "B", "DejaVuSans-Bold.ttf")

	// Brand Header
	s.drawBrandHeader(pdf, companyName)

	// Title Section
	s.drawTitleSection(pdf)

	// Assessment Info
	s.drawAssessmentInfo(pdf, report)

	// Risk Thermometer
	s.drawRiskThermometer(pdf, report)

	// Category Scores Table
	s.drawCategoryTable(pdf, report)

	// Domain Scores Table
	s.drawDomainTable(pdf, report)

	// Recommendations Section
	s.drawRecommendationsSection(pdf, report)

	// Footer
	s.drawFooter(pdf)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// drawBrandHeader creates a branded header with company info
func (s *ReportPDFService) drawBrandHeader(pdf *gofpdf.Fpdf, companyName string) {
	// Dark Slate colored banner (R:51, G:65, B:85)
	pdf.SetFillColor(51, 65, 85)
	pdf.Rect(0, 0, 210, 25, "F")

	// White text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("DejaVuSans", "B", 12)

	// Company name on the left
	pdf.SetXY(15, 8)
	pdf.Cell(0, 0, companyName)

	// CONFIDENTIAL label on the right
	pdf.SetXY(130, 8)
	pdf.Cell(0, 0, "CONFIDENTIAL")

	// Date below company name
	pdf.SetFont("DejaVuSans", "", 8)
	pdf.SetXY(15, 15)
	pdf.Cell(0, 0, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02")))

	// Reset text color to black
	pdf.SetTextColor(0, 0, 0)
}

// drawTitleSection creates the main title
func (s *ReportPDFService) drawTitleSection(pdf *gofpdf.Fpdf) {
	pdf.SetY(35)
	pdf.SetFont("DejaVuSans", "B", 16)
	pdf.Cell(0, 0, "NOM-035 Psychosocial Risk Assessment Report")

	pdf.SetY(45)
	pdf.SetFont("DejaVuSans", "", 12)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "Guía II - Identificación y análisis de los factores de riesgo psicosocial")
	pdf.SetTextColor(0, 0, 0)
}

// drawAssessmentInfo displays staff and assessment details
func (s *ReportPDFService) drawAssessmentInfo(pdf *gofpdf.Fpdf, report *domain.IndividualReportDTO) {
	yPos := 60.0

	// Assessment Information Box
	pdf.SetFillColor(245, 245, 245)
	pdf.Rect(15, yPos, 180, 35, "F")

	pdf.SetFont("DejaVuSans", "B", 12)
	pdf.SetXY(20, yPos+8)
	pdf.Cell(0, 0, "Assessment Information")

	pdf.SetFont("DejaVuSans", "", 10)
	pdf.SetXY(20, yPos+18)
	pdf.Cell(0, 0, fmt.Sprintf("Staff Member: %s", report.StaffName))

	if report.Department != "" {
		pdf.SetXY(20, yPos+25)
		pdf.Cell(0, 0, fmt.Sprintf("Department: %s", report.Department))
	}

	if report.Shift != "" {
		pdf.SetXY(20, yPos+32)
		pdf.Cell(0, 0, fmt.Sprintf("Shift: %s", report.Shift))
	}

	pdf.SetXY(110, yPos+18)
	pdf.Cell(0, 0, fmt.Sprintf("Period: %d", report.Period))

	if report.CompletedAt != nil {
		pdf.SetXY(110, yPos+25)
		pdf.Cell(0, 0, fmt.Sprintf("Completed: %s", report.CompletedAt.Format("2006-01-02")))
	}
}

// drawRiskThermometer creates a visual risk level indicator
func (s *ReportPDFService) drawRiskThermometer(pdf *gofpdf.Fpdf, report *domain.IndividualReportDTO) {
	yPos := 105.0

	pdf.SetFont("DejaVuSans", "B", 12)
	pdf.SetXY(15, yPos)
	pdf.Cell(0, 0, "Risk Assessment")

	// Total Score
	pdf.SetFont("DejaVuSans", "", 14)
	pdf.SetXY(15, yPos+10)
	pdf.Cell(0, 0, fmt.Sprintf("Total Score: %.0f/100", report.TotalScore))

	// Risk Level Text
	pdf.SetXY(15, yPos+18)
	pdf.Cell(0, 0, fmt.Sprintf("Risk Level: %s", strings.ToUpper(string(report.RiskLevel))))

	// Thermometer Scale (horizontal bar)
	scaleStartX := 15.0
	scaleEndX := 195.0
	scaleY := yPos + 25
	scaleHeight := 8.0

	// Draw background scale
	pdf.SetFillColor(240, 240, 240)
	pdf.Rect(scaleStartX, scaleY, scaleEndX-scaleStartX, scaleHeight, "F")

	// Calculate marker position (0-100 scale)
	scorePercent := report.TotalScore / 100.0
	markerX := scaleStartX + (scorePercent * (scaleEndX - scaleStartX))

	// Set marker color based on risk level
	switch report.RiskLevel {
	case "nulo", "bajo":
		pdf.SetFillColor(0, 128, 0) // Green
	case "medio":
		pdf.SetFillColor(255, 165, 0) // Orange
	case "alto", "muy_alto":
		pdf.SetFillColor(220, 20, 60) // Red
	default:
		pdf.SetFillColor(128, 128, 128) // Gray
	}

	// Draw marker (circle)
	markerRadius := 5.0
	pdf.Circle(markerX, scaleY+scaleHeight/2, markerRadius, "F")

	// Draw scale labels
	pdf.SetFont("DejaVuSans", "", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.SetXY(scaleStartX, scaleY+scaleHeight+3)
	pdf.Cell(0, 0, "0")
	pdf.SetXY(scaleEndX-5, scaleY+scaleHeight+3)
	pdf.Cell(0, 0, "100")
	pdf.SetTextColor(0, 0, 0)
}

// drawCategoryTable creates a professional table for category scores
func (s *ReportPDFService) drawCategoryTable(pdf *gofpdf.Fpdf, report *domain.IndividualReportDTO) {
	yPos := 140.0

	pdf.SetFont("DejaVuSans", "B", 12)
	pdf.SetXY(15, yPos)
	pdf.Cell(0, 0, "Category Scores")

	yPos += 10

	// Table headers
	headers := []string{"Categoría", "Puntaje", "Nivel de Riesgo"}
	colWidths := []float64{80, 30, 50}

	// Header row
	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("DejaVuSans", "B", 10)
	xPos := 15.0

	for i, header := range headers {
		pdf.Rect(xPos, yPos, colWidths[i], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, header)
		xPos += colWidths[i]
	}

	yPos += 8

	// Data rows
	row := 0
	for category, score := range report.CategoryScores {
		// Alternate row colors
		if row%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		riskLevel := s.calculateRiskLevel(score, 20) // Max score for categories is 20

		pdf.SetFont("DejaVuSans", "", 9)
		xPos = 15.0

		// Category
		pdf.Rect(xPos, yPos, colWidths[0], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, category)
		xPos += colWidths[0]

		// Score
		pdf.Rect(xPos, yPos, colWidths[1], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, fmt.Sprintf("%.1f", score))
		xPos += colWidths[1]

		// Risk Level
		pdf.Rect(xPos, yPos, colWidths[2], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, riskLevel)

		yPos += 8
		row++
	}
}

// drawDomainTable creates a professional table for domain scores
func (s *ReportPDFService) drawDomainTable(pdf *gofpdf.Fpdf, report *domain.IndividualReportDTO) {
	yPos := 200.0

	pdf.SetFont("DejaVuSans", "B", 12)
	pdf.SetXY(15, yPos)
	pdf.Cell(0, 0, "Domain Scores")

	yPos += 10

	// Table headers
	headers := []string{"Dominio", "Puntaje", "Nivel de Riesgo"}
	colWidths := []float64{80, 30, 50}

	// Header row
	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("DejaVuSans", "B", 10)
	xPos := 15.0

	for i, header := range headers {
		pdf.Rect(xPos, yPos, colWidths[i], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, header)
		xPos += colWidths[i]
	}

	yPos += 8

	// Data rows
	row := 0
	for domain, score := range report.DomainScores {
		// Alternate row colors
		if row%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		riskLevel := s.calculateRiskLevel(score, 15) // Max score for domains is 15

		pdf.SetFont("DejaVuSans", "", 9)
		xPos = 15.0

		// Domain
		pdf.Rect(xPos, yPos, colWidths[0], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, domain)
		xPos += colWidths[0]

		// Score
		pdf.Rect(xPos, yPos, colWidths[1], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, fmt.Sprintf("%.1f", score))
		xPos += colWidths[1]

		// Risk Level
		pdf.Rect(xPos, yPos, colWidths[2], 8, "FD")
		pdf.SetXY(xPos+2, yPos+2)
		pdf.Cell(0, 0, riskLevel)

		yPos += 8
		row++
	}
}

// drawRecommendationsSection creates a formatted recommendations section
func (s *ReportPDFService) drawRecommendationsSection(pdf *gofpdf.Fpdf, report *domain.IndividualReportDTO) {
	yPos := 260.0

	pdf.SetFont("DejaVuSans", "B", 12)
	pdf.SetXY(15, yPos)
	pdf.Cell(0, 0, "Recommendations")

	// Background box
	boxHeight := float64(len(report.Recommendations)*8 + 20)
	pdf.SetFillColor(248, 248, 248)
	pdf.Rect(15, yPos+5, 180, boxHeight, "F")

	// Warning icon for high risk
	if report.RiskLevel == "alto" || report.RiskLevel == "muy_alto" {
		pdf.SetFillColor(220, 20, 60)
		pdf.Circle(20, yPos+12, 3, "F")
		pdf.SetFont("DejaVuSans", "B", 8)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetXY(18.5, yPos+10)
		pdf.Cell(0, 0, "!")
		pdf.SetTextColor(0, 0, 0)
	}

	yPos += 15
	pdf.SetFont("DejaVuSans", "", 9)

	for i, rec := range report.Recommendations {
		// Checkbox
		pdf.Rect(25, yPos+1, 3, 3, "D")

		// Recommendation text
		pdf.SetXY(35, yPos+2)
		pdf.MultiCell(155, 5, fmt.Sprintf("%d. %s", i+1, rec), "", "", false)
		yPos += 8
	}
}

// drawFooter adds page footer with branding
func (s *ReportPDFService) drawFooter(pdf *gofpdf.Fpdf) {
	// Add footer on every page
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("DejaVuSans", "", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.Cell(0, 0, "Entorno35 Platform - NOM-035 Compliance Report - Page "+strconv.Itoa(pdf.PageNo()))
	})
}

// calculateRiskLevel determines risk level based on score and max score
func (s *ReportPDFService) calculateRiskLevel(score, maxScore float64) string {
	percentage := score / maxScore

	switch {
	case percentage < 0.25:
		return "Nulo"
	case percentage < 0.50:
		return "Bajo"
	case percentage < 0.75:
		return "Medio"
	case percentage <= 1.0:
		return "Alto"
	default:
		return "N/A"
	}
}