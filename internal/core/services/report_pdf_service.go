package services

import (
	"bytes"
	"fmt"
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

// GenerateIndividualReportPDF creates a professional PDF report for an individual assessment
func (s *ReportPDFService) GenerateIndividualReportPDF(report *domain.IndividualReportDTO, companyName string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 25.0)
	pdf.AddPage()

	// Set up footer
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.Cell(0, 0, "Plataforma Entorno35 - Reporte de Cumplimiento NOM-035 - Pagina "+fmt.Sprintf("%d", pdf.PageNo()))
	})

	// 1. Header (Dark Slate Banner)
	pdf.SetFillColor(51, 65, 85)
	pdf.Rect(0, 0, 210, 30, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 16)
	pdf.Text(10, 12, "Plataforma Entorno35")
	pdf.SetFont("Arial", "", 10)
	pdf.Text(10, 20, fmt.Sprintf("Fecha de Generacion: %s", time.Now().Format("2006-01-02")))
	pdf.SetFont("Arial", "B", 12)
	pdf.Text(150, 12, "CONFIDENCIAL")

	// 2. Title Section
	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(40)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Evaluacion de Riesgo Psicosocial NOM-035")
	pdf.Ln(8)
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(100, 100, 100)
	guideText := fmt.Sprintf("Reporte Individual - Analisis Guia %s", string(report.GuideType))
	pdf.Cell(0, 10, guideText)
	pdf.Ln(15)

	// 3. Staff Info Box
	pdf.SetFillColor(245, 245, 245)
	pdf.Rect(10, pdf.GetY(), 190, 35, "F")
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(0, 0, 0)
	yStart := pdf.GetY() + 8

	pdf.SetXY(15, yStart)
	pdf.Cell(0, 0, fmt.Sprintf("Empleado: %s", report.StaffName))

	pdf.SetXY(15, yStart+8)
	department := report.Department
	if department == "" {
		department = "No especificado"
	}
	pdf.Cell(0, 0, fmt.Sprintf("Departamento: %s", department))

	pdf.SetXY(15, yStart+16)
	shift := report.Shift
	if shift == "" {
		shift = "No especificado"
	}
	pdf.Cell(0, 0, fmt.Sprintf("Turno: %s", shift))

	pdf.SetXY(15, yStart+24)
	pdf.Cell(0, 0, fmt.Sprintf("Periodo de Evaluacion: %d", report.Period))

	pdf.SetXY(110, yStart)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 0, fmt.Sprintf("Puntuacion Total: %.1f", report.TotalScore))

	pdf.SetXY(110, yStart+8)
	riskLabel := formatRiskLevel(string(report.RiskLevel))
	pdf.Cell(0, 0, fmt.Sprintf("Nivel de Riesgo: %s", riskLabel))

	pdf.SetY(yStart + 30)
	pdf.Ln(10)

	// 4. Risk Thermometer (Robust Primitive)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Visualizacion del Nivel de Riesgo")
	pdf.Ln(8)

	// Draw Background Bar
	barX, barY, barW := 15.0, pdf.GetY(), 180.0
	pdf.SetFillColor(220, 220, 220)
	pdf.Rect(barX, barY, barW, 6, "F")

	// Axis Labels (0 and 100) - Positioned relative to bar coordinates
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.SetXY(barX, barY+8)  // Position cursor at bar start, below bar
	pdf.Cell(0, 0, "0")      // Use Cell instead of Text for consistent positioning
	pdf.SetXY(barX+barW-5, barY+8)  // Position at bar end
	pdf.Cell(0, 0, "100")
	pdf.SetTextColor(0, 0, 0)

	// Draw Indicator
	score := report.TotalScore
	if score > 100 { score = 100 }
	if score < 0 { score = 0 }
	indicatorX := barX + (score/100.0 * barW)

	// Color Logic
	riskLevelStr := string(report.RiskLevel)
	if riskLevelStr == "nulo" || riskLevelStr == "bajo" {
		pdf.SetFillColor(0, 150, 0)
	} else if riskLevelStr == "medio" {
		pdf.SetFillColor(255, 165, 0)
	} else {
		pdf.SetFillColor(200, 0, 0)
	}
	pdf.Circle(indicatorX, barY+3, 3, "F")

	// Risk Label - Positioned above the bar, centered on indicator
	pdf.SetFont("Arial", "B", 9)
	labelWidth := pdf.GetStringWidth(strings.ToUpper(riskLevelStr))
	labelX := indicatorX - (labelWidth / 2)
	if labelX < 15 { labelX = 15 } // Keep within left margin
	if labelX + labelWidth > 195 { labelX = 195 - labelWidth } // Keep within right margin

	pdf.SetXY(labelX, barY-6) // Position above bar with proper spacing
	pdf.Cell(0, 0, strings.ToUpper(riskLevelStr))
	pdf.Ln(15) // Less spacing after risk label

	// 5. Data Tables
	if pdf.GetY() > 220 {
		pdf.AddPage()
	}

	drawTableWithRiskLevels(pdf, "Puntuacion por Categoria", report.CategoryScores, report.CategoryMaxScores, report.CategoryRiskLevels)
	pdf.Ln(10)

	if pdf.GetY() > 220 {
		pdf.AddPage()
	}

	drawTableWithRiskLevels(pdf, "Puntuacion por Dominio", report.DomainScores, report.DomainMaxScores, report.DomainRiskLevels)
	pdf.Ln(12)

	// 6. Recommendations - Only show if there are recommendations
	if len(report.Recommendations) > 0 {
		if pdf.GetY() > 200 {
			pdf.AddPage()
		}

		pdf.SetFont("Arial", "B", 14)
		pdf.SetFillColor(240, 240, 255)
		pdf.CellFormat(0, 12, " Recomendaciones", "0", 1, "", true, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.Ln(6)

		for i, rec := range report.Recommendations {
			currentY := pdf.GetY()
			if currentY > 270 {
				pdf.AddPage()
				pdf.SetFont("Arial", "", 10)
			}

			pdf.SetX(15)
			pdf.SetFillColor(100, 100, 100)
			pdf.Circle(17, pdf.GetY()+3, 1, "F")
			pdf.SetX(22)

			pdf.MultiCell(175, 7, fmt.Sprintf("%d. %s", i+1, rec), "", "", false)
			pdf.Ln(4)
		}
	} else {
		// For low risk, add a positive summary instead of empty section
		if pdf.GetY() > 200 {
			pdf.AddPage()
		}

		pdf.SetFont("Arial", "B", 14)
		pdf.SetFillColor(220, 252, 231) // Light green background
		pdf.CellFormat(0, 12, " Resultado de la Evaluacion", "0", 1, "", true, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.Ln(6)

		pdf.SetX(15)
		riskLabel := formatRiskLevel(string(report.RiskLevel))
		summaryText := fmt.Sprintf("El empleado presenta un nivel de riesgo %s. No se requieren acciones correctivas inmediatas. Se recomienda mantener las condiciones laborales actuales y continuar con el monitoreo periodico.", riskLabel)
		pdf.MultiCell(180, 6, summaryText, "", "", false)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF output: %w", err)
	}

	return buf.Bytes(), nil
}

func drawTableWithRiskLevels(pdf *gofpdf.Fpdf, title string, scores map[string]float64, maxScores map[string]float64, riskLevels map[string]string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, title)
	pdf.Ln(10)

	// Header - 3 columns with Spanish labels
	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(85, 10, "Elemento", "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 10, "Puntaje/Max", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 10, "Nivel de Riesgo", "1", 1, "C", true, 0, "")

	// Rows - 3 columns
	pdf.SetFont("Arial", "", 9)
	fill := false
	for k, v := range scores {
		// Check if we need a page break before this row
		if pdf.GetY() > 250 {
			pdf.AddPage()
			pdf.SetFont("Arial", "", 9)
			// Re-draw header on new page for continuity
			pdf.SetFillColor(230, 230, 230)
			pdf.SetFont("Arial", "B", 9)
			pdf.CellFormat(85, 10, "Elemento", "1", 0, "L", true, 0, "")
			pdf.CellFormat(45, 10, "Puntaje/Max", "1", 0, "C", true, 0, "")
			pdf.CellFormat(50, 10, "Nivel de Riesgo", "1", 1, "C", true, 0, "")
			pdf.SetFont("Arial", "", 9)
			fill = false // Reset fill pattern
		}

		if fill {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Ensure text fits in the cell - truncate if too long
		itemText := k
		if pdf.GetStringWidth(itemText) > 80 {
			// Truncate with ellipsis
			for len(itemText) > 0 && pdf.GetStringWidth(itemText+"...") > 80 {
				itemText = itemText[:len(itemText)-1]
			}
			itemText += "..."
		}

		// Get max score and risk level
		maxScore := maxScores[k]
		riskLevel := riskLevels[k]
		
		// Format risk level for display
		riskDisplay := formatRiskLevel(riskLevel)

		// 3 columns: Item, Score/Max, Risk Level
		pdf.CellFormat(85, 8, itemText, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(45, 8, fmt.Sprintf("%.1f / %.0f", v, maxScore), "1", 0, "C", fill, 0, "")
		
		// Color-code risk level text
		pdf.SetTextColor(getRiskColor(riskLevel))
		pdf.CellFormat(50, 8, riskDisplay, "1", 1, "C", fill, 0, "")
		pdf.SetTextColor(0, 0, 0) // Reset to black
		
		fill = !fill
	}
	pdf.Ln(5) // Add some space after table
}

// formatRiskLevel converts risk level to readable format
func formatRiskLevel(level string) string {
	switch level {
	case "nulo":
		return "Nulo"
	case "bajo":
		return "Bajo"
	case "medio":
		return "Medio"
	case "alto":
		return "Alto"
	case "muy_alto":
		return "Muy Alto"
	default:
		return strings.Title(level)
	}
}

// getRiskColor returns RGB color based on risk level
func getRiskColor(level string) (r, g, b int) {
	switch level {
	case "nulo":
		return 34, 197, 94 // Green
	case "bajo":
		return 132, 204, 22 // Light green
	case "medio":
		return 234, 179, 8 // Yellow
	case "alto":
		return 249, 115, 22 // Orange
	case "muy_alto":
		return 239, 68, 68 // Red
	default:
		return 0, 0, 0 // Black
	}
}

// Legacy function kept for compatibility (not used anymore)
func drawTable(pdf *gofpdf.Fpdf, title string, data map[string]float64) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, title)
	pdf.Ln(10)

	// Header - STRICTLY 2 columns only
	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(130, 10, "Item", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 10, "Score", "1", 1, "C", true, 0, "")

	// Rows - STRICTLY 2 columns only
	pdf.SetFont("Arial", "", 9)
	fill := false
	for k, v := range data {
		// Check if we need a page break before this row
		if pdf.GetY() > 250 {
			pdf.AddPage()
			pdf.SetFont("Arial", "", 9)
			// Re-draw header on new page for continuity
			pdf.SetFillColor(230, 230, 230)
			pdf.SetFont("Arial", "B", 10)
			pdf.CellFormat(130, 10, "Item", "1", 0, "L", true, 0, "")
			pdf.CellFormat(50, 10, "Score", "1", 1, "C", true, 0, "")
			pdf.SetFont("Arial", "", 9)
			fill = false // Reset fill pattern
		}

		if fill {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Ensure text fits in the cell - truncate if too long
		itemText := k
		if pdf.GetStringWidth(itemText) > 125 {
			// Truncate with ellipsis
			for len(itemText) > 0 && pdf.GetStringWidth(itemText+"...") > 125 {
				itemText = itemText[:len(itemText)-1]
			}
			itemText += "..."
		}

		// STRICTLY 2 columns: Item and Score only
		pdf.CellFormat(130, 8, itemText, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(50, 8, fmt.Sprintf("%.1f", v), "1", 1, "C", fill, 0, "")
		fill = !fill
	}
	pdf.Ln(5) // Add some space after table
}