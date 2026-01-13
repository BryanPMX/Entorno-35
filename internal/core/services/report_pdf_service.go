package services

import (
	"bytes"
	"fmt"
	"sort"
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
	pdf.SetAutoPageBreak(true, 20.0)
	pdf.AddPage()

	// Professional footer
	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(128, 128, 128)
		footerText := fmt.Sprintf("Entorno35 | Reporte NOM-035 | %s | Pagina %d", time.Now().Format("02/01/2006"), pdf.PageNo())
		pdf.CellFormat(0, 10, footerText, "", 0, "C", false, 0, "")
	})

	// ========== HEADER SECTION ==========
	// Modern gradient header
	pdf.SetFillColor(15, 23, 42) // Slate-900
	pdf.Rect(0, 0, 210, 35, "F")

	// Accent stripe
	pdf.SetFillColor(59, 130, 246) // Blue-500
	pdf.Rect(0, 35, 210, 3, "F")

	// Header text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 18)
	pdf.SetXY(15, 10)
	pdf.Cell(0, 0, "Reporte de Riesgo Psicosocial")

	pdf.SetFont("Arial", "", 11)
	pdf.SetXY(15, 20)
	pdf.Cell(0, 0, fmt.Sprintf("NOM-035-STPS-2018 | Guia %s", string(report.GuideType)))

	// Confidential badge
	pdf.SetFillColor(239, 68, 68) // Red-500
	pdf.RoundedRect(155, 8, 42, 8, 2, "1234", "F")
	pdf.SetFont("Arial", "B", 8)
	pdf.SetXY(155, 10)
	pdf.CellFormat(42, 0, "CONFIDENCIAL", "", 0, "C", false, 0, "")

	// Date
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(155, 22)
	pdf.Cell(0, 0, time.Now().Format("02 de Enero, 2006"))

	pdf.SetY(48)

	// ========== EMPLOYEE INFO CARD ==========
	pdf.SetFillColor(248, 250, 252) // Slate-50
	pdf.RoundedRect(15, pdf.GetY(), 180, 32, 4, "1234", "F")

	infoY := pdf.GetY() + 6
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(22, infoY)
	pdf.Cell(0, 0, report.StaffName)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 116, 139)

	// Row 1: Department & Shift
	pdf.SetXY(22, infoY+10)
	department := report.Department
	if department == "" {
		department = "Sin especificar"
	}
	pdf.Cell(0, 0, fmt.Sprintf("Departamento: %s", department))

	pdf.SetXY(110, infoY+10)
	shift := report.Shift
	if shift == "" {
		shift = "Sin especificar"
	}
	pdf.Cell(0, 0, fmt.Sprintf("Turno: %s", shift))

	// Row 2: Period
	pdf.SetXY(22, infoY+18)
	pdf.Cell(0, 0, fmt.Sprintf("Periodo de Evaluacion: %d", report.Period))

	pdf.SetY(infoY + 35)

	// ========== SCORE SUMMARY CARD ==========
	riskLevel := string(report.RiskLevel)
	riskConfig := getRiskColorConfig(riskLevel)

	// Main score card with risk-based accent
	pdf.SetFillColor(255, 255, 255)
	pdf.RoundedRect(15, pdf.GetY(), 180, 50, 4, "1234", "F")

	scoreY := pdf.GetY() + 8

	// Score section
	pdf.SetFont("Arial", "B", 36)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(28, scoreY)
	pdf.Cell(0, 0, fmt.Sprintf("%.0f", report.TotalScore))

	// Max score
	maxScore := report.TotalMaxScore
	if maxScore == 0 {
		maxScore = 184
	}
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(148, 163, 184)
	scoreWidth := pdf.GetStringWidth(fmt.Sprintf("%.0f", report.TotalScore))
	pdf.SetXY(28+scoreWidth+3, scoreY+8)
	pdf.Cell(0, 0, fmt.Sprintf("/ %.0f puntos", maxScore))

	// Percentage
	percentage := (report.TotalScore / maxScore) * 100
	if percentage > 100 {
		percentage = 100
	}
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(28, scoreY+18)
	pdf.Cell(0, 0, fmt.Sprintf("%.1f%% del maximo posible", percentage))

	// Risk badge on right
	badgeX := 130.0
	badgeY := scoreY + 5

	pdf.SetFillColor(riskConfig.bgR, riskConfig.bgG, riskConfig.bgB)
	pdf.RoundedRect(badgeX, badgeY, 55, 28, 3, "1234", "F")

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(riskConfig.textR, riskConfig.textG, riskConfig.textB)
	riskLabel := formatRiskLevel(riskLevel)
	labelWidth := pdf.GetStringWidth(riskLabel)
	pdf.SetXY(badgeX+(55-labelWidth)/2, badgeY+8)
	pdf.Cell(0, 0, riskLabel)

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(badgeX, badgeY+18)
	pdf.CellFormat(55, 0, "Nivel de Riesgo", "", 0, "C", false, 0, "")

	// Progress bar
	barY := scoreY + 32
	pdf.SetFillColor(226, 232, 240)
	pdf.RoundedRect(28, barY, 155, 6, 2, "1234", "F")

	fillWidth := (percentage / 100.0) * 155
	if fillWidth > 0 {
		pdf.SetFillColor(riskConfig.barR, riskConfig.barG, riskConfig.barB)
		pdf.RoundedRect(28, barY, fillWidth, 6, 2, "1234", "F")
	}

	pdf.SetY(pdf.GetY() + 60)

	// Medical attention warning
	if report.RequiresMedical {
		pdf.SetFillColor(254, 226, 226)
		pdf.RoundedRect(15, pdf.GetY(), 180, 14, 3, "1234", "F")
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(185, 28, 28)
		pdf.SetXY(22, pdf.GetY()+5)
		pdf.Cell(0, 0, "ATENCION: Esta evaluacion indica necesidad de atencion medica ocupacional")
		pdf.Ln(20)
	}

	// ========== CATEGORY ANALYSIS ==========
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}

	drawModernSection(pdf, "Analisis por Categorias", "Puntuacion y nivel de riesgo por categoria NOM-035")
	pdf.Ln(8)

	drawModernTable(pdf, report.CategoryScores, report.CategoryMaxScores, report.CategoryRiskLevels)

	// ========== DOMAIN ANALYSIS ==========
	if pdf.GetY() > 180 {
		pdf.AddPage()
	}

	pdf.Ln(8)
	drawModernSection(pdf, "Analisis por Dominios", "Puntuacion detallada por dominio de riesgo psicosocial")
	pdf.Ln(8)

	drawModernTable(pdf, report.DomainScores, report.DomainMaxScores, report.DomainRiskLevels)

	// ========== RECOMMENDATIONS ==========
	if len(report.Recommendations) > 0 {
		if pdf.GetY() > 180 {
			pdf.AddPage()
		}

		pdf.Ln(8)
		drawModernSection(pdf, "Recomendaciones", "Acciones sugeridas segun la NOM-035")
		pdf.Ln(8)

		for i, rec := range report.Recommendations {
			if pdf.GetY() > 260 {
				pdf.AddPage()
			}

			// Recommendation card
			pdf.SetFillColor(248, 250, 252)
			pdf.RoundedRect(15, pdf.GetY(), 180, 16, 3, "1234", "F")

			// Number badge
			pdf.SetFillColor(59, 130, 246)
			pdf.Circle(25, pdf.GetY()+8, 5, "F")
			pdf.SetFont("Arial", "B", 9)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetXY(22, pdf.GetY()+6)
			pdf.Cell(6, 0, fmt.Sprintf("%d", i+1))

			// Text
			pdf.SetFont("Arial", "", 9)
			pdf.SetTextColor(51, 65, 85)
			pdf.SetXY(35, pdf.GetY()+4)
			pdf.MultiCell(155, 5, rec, "", "L", false)

			pdf.Ln(4)
		}
	} else {
		// Positive result for low risk
		if pdf.GetY() > 220 {
			pdf.AddPage()
		}

		pdf.Ln(8)
		pdf.SetFillColor(220, 252, 231)
		pdf.RoundedRect(15, pdf.GetY(), 180, 30, 4, "1234", "F")

		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(21, 128, 61)
		pdf.SetXY(22, pdf.GetY()+8)
		pdf.Cell(0, 0, "Resultado Favorable")

		pdf.SetFont("Arial", "", 10)
		pdf.SetXY(22, pdf.GetY()+10)
		riskLabel := formatRiskLevel(riskLevel)
		pdf.MultiCell(165, 5, fmt.Sprintf("El empleado presenta un nivel de riesgo %s. No se requieren acciones correctivas inmediatas.", riskLabel), "", "L", false)
	}

	// ========== QUESTION-LEVEL DETAILS ==========
	if len(report.QuestionResponses) > 0 {
		pdf.AddPage()
		drawModernSection(pdf, "Detalle de Respuestas por Pregunta", "Respuestas individuales organizadas por categoria")
		pdf.Ln(4)

		// Group questions by category
		categoryQuestions := make(map[string][]domain.QuestionResponseDetail)
		for _, qr := range report.QuestionResponses {
			categoryName := qr.Category
			if categoryName == "" {
				categoryName = "Sin categoria"
			}
			categoryQuestions[categoryName] = append(categoryQuestions[categoryName], qr)
		}

		// Sort category names
		categoryNames := make([]string, 0, len(categoryQuestions))
		for name := range categoryQuestions {
			categoryNames = append(categoryNames, name)
		}
		sort.Strings(categoryNames)

		// Draw questions by category
		for _, categoryName := range categoryNames {
			questions := categoryQuestions[categoryName]

			if pdf.GetY() > 220 {
				pdf.AddPage()
			}

			// Category header
			pdf.SetFillColor(241, 245, 249)
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(30, 41, 59)
			pdf.SetX(15)
			pdf.CellFormat(180, 8, categoryName, "", 1, "L", true, 0, "")
			pdf.Ln(2)

			// Questions table header
			pdf.SetFillColor(248, 250, 252)
			pdf.SetFont("Arial", "B", 8)
			pdf.SetTextColor(71, 85, 105)
			pdf.SetX(15)
			pdf.CellFormat(10, 7, "#", "", 0, "C", true, 0, "")
			pdf.CellFormat(115, 7, "Pregunta", "", 0, "L", true, 0, "")
			pdf.CellFormat(20, 7, "Resp.", "", 0, "C", true, 0, "")
			pdf.CellFormat(20, 7, "Punt.", "", 0, "C", true, 0, "")
			pdf.CellFormat(15, 7, "Max", "", 1, "C", true, 0, "")

			// Questions rows
			pdf.SetFont("Arial", "", 8)
			fill := false
			for _, q := range questions {
				if pdf.GetY() > 265 {
					pdf.AddPage()
					// Redraw category header
					pdf.SetFillColor(241, 245, 249)
					pdf.SetFont("Arial", "B", 10)
					pdf.SetTextColor(30, 41, 59)
					pdf.SetX(15)
					pdf.CellFormat(180, 8, categoryName+" (continuacion)", "", 1, "L", true, 0, "")
					pdf.Ln(2)
					// Redraw table header
					pdf.SetFillColor(248, 250, 252)
					pdf.SetFont("Arial", "B", 8)
					pdf.SetTextColor(71, 85, 105)
					pdf.SetX(15)
					pdf.CellFormat(10, 7, "#", "", 0, "C", true, 0, "")
					pdf.CellFormat(115, 7, "Pregunta", "", 0, "L", true, 0, "")
					pdf.CellFormat(20, 7, "Resp.", "", 0, "C", true, 0, "")
					pdf.CellFormat(20, 7, "Punt.", "", 0, "C", true, 0, "")
					pdf.CellFormat(15, 7, "Max", "", 1, "C", true, 0, "")
					pdf.SetFont("Arial", "", 8)
					fill = false
				}

				if fill {
					pdf.SetFillColor(248, 250, 252)
				} else {
					pdf.SetFillColor(255, 255, 255)
				}

				// Truncate question text if too long
				questionText := q.QuestionText
				maxTextWidth := 110.0
				if pdf.GetStringWidth(questionText) > maxTextWidth {
					for len(questionText) > 0 && pdf.GetStringWidth(questionText+"...") > maxTextWidth {
						questionText = questionText[:len(questionText)-1]
					}
					questionText += "..."
				}

				pdf.SetTextColor(30, 41, 59)
				pdf.SetX(15)
				pdf.CellFormat(10, 6, fmt.Sprintf("%d", q.QuestionNumber), "", 0, "C", fill, 0, "")
				pdf.CellFormat(115, 6, questionText, "", 0, "L", fill, 0, "")
				pdf.CellFormat(20, 6, fmt.Sprintf("%d", q.SelectedValue), "", 0, "C", fill, 0, "")

				// Color-code the score
				scorePercent := float64(q.CalculatedScore) / float64(q.MaxScore)
				if scorePercent >= 0.75 {
					pdf.SetTextColor(185, 28, 28) // Red for high risk scores
				} else if scorePercent >= 0.5 {
					pdf.SetTextColor(161, 98, 7) // Amber for medium
				} else {
					pdf.SetTextColor(21, 128, 61) // Green for low
				}
				pdf.CellFormat(20, 6, fmt.Sprintf("%d", q.CalculatedScore), "", 0, "C", fill, 0, "")

				pdf.SetTextColor(100, 116, 139)
				pdf.CellFormat(15, 6, fmt.Sprintf("%d", q.MaxScore), "", 1, "C", fill, 0, "")

				fill = !fill
			}

			pdf.Ln(6)
		}
	}

	// ========== LEGAL FOOTER ==========
	if pdf.GetY() > 240 {
		pdf.AddPage()
	}

	pdf.Ln(15)
	pdf.SetFillColor(241, 245, 249)
	pdf.RoundedRect(15, pdf.GetY(), 180, 25, 3, "1234", "F")

	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(22, pdf.GetY()+6)
	pdf.MultiCell(165, 4, "Este documento contiene informacion confidencial protegida por la Ley Federal de Proteccion de Datos Personales. Su distribucion no autorizada esta prohibida. Generado automaticamente por la plataforma Entorno35 conforme a la NOM-035-STPS-2018.", "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("error al generar PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// drawModernSection draws a section header with improved spacing
func drawModernSection(pdf *gofpdf.Fpdf, title, subtitle string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetX(15)
	pdf.Cell(0, 8, title)
	pdf.Ln(10) // Increased from 6 to 10 for more space between title and subtitle

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetX(15)
	pdf.Cell(0, 5, subtitle)
	pdf.Ln(8)
}

// drawModernTable draws a modern styled table
func drawModernTable(pdf *gofpdf.Fpdf, scores, maxScores map[string]float64, riskLevels map[string]string) {
	// Sort keys
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Header
	pdf.SetFillColor(241, 245, 249)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(71, 85, 105)

	pdf.SetX(15)
	pdf.CellFormat(95, 10, "Elemento", "", 0, "L", true, 0, "")
	pdf.CellFormat(40, 10, "Puntuacion", "", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Nivel de Riesgo", "", 1, "C", true, 0, "")

	// Rows
	pdf.SetFont("Arial", "", 9)
	fill := false

	for _, k := range keys {
		if pdf.GetY() > 250 {
			pdf.AddPage()
			// Redraw header
			pdf.SetFillColor(241, 245, 249)
			pdf.SetFont("Arial", "B", 9)
			pdf.SetTextColor(71, 85, 105)
			pdf.SetX(15)
			pdf.CellFormat(95, 10, "Elemento", "", 0, "L", true, 0, "")
			pdf.CellFormat(40, 10, "Puntuacion", "", 0, "C", true, 0, "")
			pdf.CellFormat(45, 10, "Nivel de Riesgo", "", 1, "C", true, 0, "")
			pdf.SetFont("Arial", "", 9)
			fill = false
		}

		score := scores[k]
		maxScore := maxScores[k]
		riskLevel := riskLevels[k]
		config := getRiskColorConfig(riskLevel)

		if fill {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Item name (truncate if needed)
		itemText := k
		if pdf.GetStringWidth(itemText) > 90 {
			for len(itemText) > 0 && pdf.GetStringWidth(itemText+"...") > 90 {
				itemText = itemText[:len(itemText)-1]
			}
			itemText += "..."
		}

		pdf.SetTextColor(30, 41, 59)
		pdf.SetX(15)
		pdf.CellFormat(95, 9, itemText, "", 0, "L", fill, 0, "")

		// Score
		pdf.SetTextColor(71, 85, 105)
		pdf.CellFormat(40, 9, fmt.Sprintf("%.1f / %.0f", score, maxScore), "", 0, "C", fill, 0, "")

		// Risk level badge
		riskLabel := formatRiskLevel(riskLevel)
		pdf.SetTextColor(config.textR, config.textG, config.textB)
		pdf.CellFormat(45, 9, riskLabel, "", 1, "C", fill, 0, "")

		fill = !fill
	}
}

// riskColorConfig holds color configuration for risk levels
type riskColorConfig struct {
	bgR, bgG, bgB       int
	textR, textG, textB int
	barR, barG, barB    int
}

// getRiskColorConfig returns color configuration for a risk level
func getRiskColorConfig(level string) riskColorConfig {
	configs := map[string]riskColorConfig{
		"nulo": {
			bgR: 220, bgG: 252, bgB: 231,
			textR: 21, textG: 128, textB: 61,
			barR: 34, barG: 197, barB: 94,
		},
		"bajo": {
			bgR: 236, bgG: 252, bgB: 203,
			textR: 77, textG: 124, textB: 15,
			barR: 132, barG: 204, barB: 22,
		},
		"medio": {
			bgR: 254, bgG: 243, bgB: 199,
			textR: 161, textG: 98, textB: 7,
			barR: 245, barG: 158, barB: 11,
		},
		"alto": {
			bgR: 255, bgG: 237, bgB: 213,
			textR: 194, textG: 65, textB: 12,
			barR: 249, barG: 115, barB: 22,
		},
		"muy_alto": {
			bgR: 254, bgG: 226, bgB: 226,
			textR: 185, textG: 28, textB: 28,
			barR: 239, barG: 68, barB: 68,
		},
	}

	if config, exists := configs[level]; exists {
		return config
	}
	return configs["nulo"]
}

// formatRiskLevel converts risk level to readable format
func formatRiskLevel(level string) string {
	labels := map[string]string{
		"nulo":     "Nulo",
		"bajo":     "Bajo",
		"medio":    "Medio",
		"alto":     "Alto",
		"muy_alto": "Muy Alto",
	}
	if label, exists := labels[level]; exists {
		return label
	}
	return strings.Title(level)
}

// getRiskColor returns RGB text color based on risk level
func getRiskColor(level string) (r, g, b int) {
	config := getRiskColorConfig(level)
	return config.textR, config.textG, config.textB
}

// getRiskBgColor returns light background RGB color based on risk level
func getRiskBgColor(level string) (r, g, b int) {
	config := getRiskColorConfig(level)
	return config.bgR, config.bgG, config.bgB
}

// getRiskBarColor returns solid bar RGB color based on risk level
func getRiskBarColor(level string) (r, g, b int) {
	config := getRiskColorConfig(level)
	return config.barR, config.barG, config.barB
}

// GenerateGeneralReportPDF creates a professional PDF report for company-wide assessment statistics
func (s *ReportPDFService) GenerateGeneralReportPDF(report *domain.GeneralReportDTO) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20.0)
	pdf.AddPage()

	// Professional footer with consistent format
	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(128, 128, 128)
		footerText := fmt.Sprintf("Entorno35 | Reporte General NOM-035 | %s | Pagina %d", time.Now().Format("02/01/2006"), pdf.PageNo())
		pdf.CellFormat(0, 10, footerText, "", 0, "C", false, 0, "")
	})

	// ========== HEADER ==========
	// Dark blue header background (#1A1F36 = 26, 31, 54)
	pdf.SetFillColor(26, 31, 54)
	pdf.Rect(0, 0, 210, 40, "F")

	// Main title - white text, bold, 24pt
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetXY(15, 12)
	pdf.Cell(0, 0, "Reporte General de Cumplimiento")

	// Subtitle - lighter gray, 14pt
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(200, 200, 200)
	pdf.SetXY(15, 24)
	periodText := "Todos los Periodos"
	if report.Period != nil {
		periodText = fmt.Sprintf("Periodo %d", *report.Period)
	}
	pdf.Cell(0, 0, fmt.Sprintf("NOM-035-STPS-2018 | %s", periodText))

	// Red CONFIDENCIAL pill-shaped button (#FF4D4F = 255, 77, 79)
	pdf.SetFillColor(255, 77, 79)
	badgeX := 155.0
	badgeY := 10.0
	badgeWidth := 45.0
	badgeHeight := 10.0
	pdf.RoundedRect(badgeX, badgeY, badgeWidth, badgeHeight, 5, "1234", "F")
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(badgeX, badgeY)
	pdf.CellFormat(badgeWidth, badgeHeight, "CONFIDENCIAL", "", 0, "CM", false, 0, "")

	pdf.SetY(48) // Spacing after header

	// ========== COMPANY INFO CARD ==========
	// Centered company details in gray text (12pt)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(100, 116, 139) // Gray text
	
	// Center company name
	companyNameWidth := pdf.GetStringWidth(report.CompanyName)
	companyNameX := (210 - companyNameWidth) / 2
	pdf.SetXY(companyNameX, pdf.GetY())
	pdf.Cell(0, 0, report.CompanyName)
	
	// Center company ID
	companyIDStr := report.CompanyID.String()
	if len(companyIDStr) > 12 {
		companyIDStr = companyIDStr[:12] + "..."
	}
	idText := fmt.Sprintf("ID: %s", companyIDStr)
	idWidth := pdf.GetStringWidth(idText)
	idX := (210 - idWidth) / 2
	pdf.SetXY(idX, pdf.GetY()+8)
	pdf.Cell(0, 0, idText)
	
	// Center date
	dateText := time.Now().Format("02 de Enero, 2006")
	dateWidth := pdf.GetStringWidth(dateText)
	dateX := (210 - dateWidth) / 2
	pdf.SetXY(dateX, pdf.GetY()+8)
	pdf.Cell(0, 0, fmt.Sprintf("Fecha: %s", dateText))

	pdf.SetY(pdf.GetY() + 20) // Spacing after company details

	// ========== EXECUTIVE SUMMARY ==========
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}
	
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Resumen Ejecutivo")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	summaryText := fmt.Sprintf("Este reporte presenta un analisis integral del cumplimiento de la NOM-035-STPS-2018 para %s. "+
		"El analisis incluye %d empleados, de los cuales %d han completado sus evaluaciones, "+
		"representando una tasa de participacion del %.1f%%. Los resultados muestran la distribucion de riesgos psicosociales "+
		"identificados y proporcionan recomendaciones para la prevencion y mejora continua.",
		report.CompanyName, report.TotalStaff, report.CompletedAssessments, report.ParticipationRate)
	pdf.MultiCell(180, 5, summaryText, "", "L", false)
	pdf.Ln(8)

	// ========== KEY METRICS ==========
	// Blue subtitle (16pt)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(59, 130, 246) // Blue
	pdf.SetX(15)
	pdf.Cell(0, 8, "Metricas Clave")
	pdf.Ln(10)

	metricsY := pdf.GetY()
	metricWidth := 58.0
	metricHeight := 50.0
	// Calculate spacing for even distribution
	availableWidth := 180.0
	totalMetricsWidth := 3 * metricWidth
	totalSpacing := availableWidth - totalMetricsWidth
	spacing := totalSpacing / 4.0
	startX := 15.0 + spacing

	// Metric 1: Total Staff - Light blue (#E6F0FF = 230, 240, 255)
	drawEnhancedMetricCard(pdf, startX, metricsY, metricWidth, metricHeight,
		fmt.Sprintf("%d", report.TotalStaff), "Personal Total",
		30, 41, 59, // Dark text
		230, 240, 255) // Light blue background

	// Metric 2: Completed - Light green (#E6FFE6 = 230, 255, 230)
	drawEnhancedMetricCard(pdf, startX+metricWidth+spacing, metricsY, metricWidth, metricHeight,
		fmt.Sprintf("%d", report.CompletedAssessments), "Evaluaciones Completas",
		30, 41, 59, // Dark text
		230, 255, 230) // Light green background

	// Metric 3: Participation - Light yellow (#FFFBE6 = 255, 251, 230)
	drawEnhancedMetricCard(pdf, startX+2*(metricWidth+spacing), metricsY, metricWidth, metricHeight,
		fmt.Sprintf("%.1f%%", report.ParticipationRate), "Tasa de Participacion",
		30, 41, 59, // Dark text
		255, 251, 230) // Light yellow background

	pdf.SetY(metricsY + metricHeight + 15)

	// ========== RISK DISTRIBUTION ==========
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}
	
	// Subtitle
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(59, 130, 246) // Blue
	pdf.SetX(15)
	pdf.Cell(0, 8, "Distribucion de Riesgo")
	pdf.Ln(10)
	
	// Interpretive analysis paragraph
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	analysisText := "La distribucion de riesgo muestra la cantidad de evaluaciones clasificadas en cada nivel segun los "+
		"criterios de la NOM-035-STPS-2018. Esta informacion es fundamental para identificar areas de atencion prioritaria "+
		"y desarrollar estrategias de prevencion efectivas."
	pdf.MultiCell(180, 5, analysisText, "", "L", false)
	pdf.Ln(6)

	if len(report.RiskDistribution) > 0 {
		riskOrder := []string{"nulo", "bajo", "medio", "alto", "muy_alto"}
		riskMap := make(map[string]int64)
		maxCount := int64(1)
		totalCount := int64(0)

		// Calculate total count and max count
		for _, rd := range report.RiskDistribution {
			riskMap[string(rd.RiskLevel)] = rd.Count
			totalCount += rd.Count
			if rd.Count > maxCount {
				maxCount = rd.Count
			}
		}

		// Improved spacing and layout - optimized for page width (210mm)
		labelWidth := 50.0
		barMaxWidth := 80.0
		barHeight := 18.0
		countWidth := 20.0
		percentWidth := 20.0
		spacing := 5.0
		barStartX := 15 + labelWidth + spacing

		// Professional table header with solid background
		pdf.SetFillColor(249, 250, 251) // #F9FAFB - light gray background
		pdf.SetFont("Arial", "B", 9)
		pdf.SetTextColor(30, 41, 59) // Dark text for contrast
		headerHeight := 12.0
		pdf.SetX(15)
		pdf.CellFormat(labelWidth, headerHeight, "Nivel de Riesgo", "1", 0, "L", true, 0, "")
		pdf.CellFormat(barMaxWidth, headerHeight, "Distribucion", "1", 0, "C", true, 0, "")
		pdf.CellFormat(countWidth, headerHeight, "Cantidad", "1", 0, "R", true, 0, "") // Right-aligned for numbers
		pdf.CellFormat(percentWidth, headerHeight, "Porcentaje", "1", 1, "R", true, 0, "") // Right-aligned for numbers
		pdf.Ln(2)

		for _, level := range riskOrder {
			count, exists := riskMap[level]
			if !exists {
				continue
			}

			// Calculate percentage
			percentage := float64(0)
			if totalCount > 0 {
				percentage = (float64(count) / float64(totalCount)) * 100
			}

			rowY := pdf.GetY()
			rowHeight := barHeight

			// Label - left-aligned
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(30, 41, 59)
			pdf.SetX(15)
			pdf.CellFormat(labelWidth, rowHeight, formatRiskLevel(level), "LR", 0, "L", false, 0, "")

			// Bar background with padding
			barY := rowY + 2
			pdf.SetFillColor(241, 245, 249)
			pdf.RoundedRect(barStartX, barY, barMaxWidth, barHeight-4, 3, "1234", "F")

			// Bar fill - Green (#28A745 = 40, 167, 69) for all risk levels
			barWidth := (float64(count) / float64(maxCount)) * barMaxWidth
			if barWidth > 0 {
				// Use green color for bars as specified
				pdf.SetFillColor(40, 167, 69) // #28A745
				pdf.RoundedRect(barStartX, barY, barWidth, barHeight-4, 3, "1234", "F")
			}

			// Count - right-aligned for numbers
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(30, 41, 59)
			countX := barStartX + barMaxWidth + spacing
			countY := rowY + (barHeight-4)/2 - 2
			pdf.SetXY(countX, countY)
			pdf.CellFormat(countWidth, 0, fmt.Sprintf("%d", count), "", 0, "R", false, 0, "")

			// Percentage - right-aligned for numbers
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(100, 116, 139)
			percentX := countX + countWidth + spacing
			pdf.SetXY(percentX, countY)
			pdf.CellFormat(percentWidth, 0, fmt.Sprintf("%.1f%%", percentage), "", 0, "R", false, 0, "")

			// Draw bottom border for row separation - light gray borders
			pdf.SetDrawColor(229, 231, 235) // #E5E7EB
			pdf.SetLineWidth(0.5)
			pdf.Line(15, rowY+rowHeight, 15+labelWidth+barMaxWidth+countWidth+percentWidth+spacing*2, rowY+rowHeight)

			pdf.SetY(rowY + rowHeight + 2)
		}

		// Total summary row with proper borders and right-aligned numbers
		if totalCount > 0 {
			pdf.Ln(2)
			pdf.SetFillColor(249, 250, 251) // Light gray background
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(30, 41, 59)
			pdf.SetX(15)
			pdf.CellFormat(labelWidth, 12, "Total", "1", 0, "L", true, 0, "")
			pdf.CellFormat(barMaxWidth, 12, "", "1", 0, "C", true, 0, "")
			pdf.CellFormat(countWidth, 12, fmt.Sprintf("%d", totalCount), "1", 0, "R", true, 0, "") // Right-aligned
			pdf.CellFormat(percentWidth, 12, "100.0%", "1", 1, "R", true, 0, "") // Right-aligned
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(148, 163, 184)
		pdf.SetX(15)
		pdf.Cell(0, 10, "No hay datos de distribucion de riesgo disponibles")
		pdf.Ln(15)
	}

	// ========== DEPARTMENT HEATMAP ==========
	if len(report.DepartmentHeatmap) > 0 {
		if pdf.GetY() > 180 {
			pdf.AddPage()
		}

		pdf.Ln(8)
		// Subtitle
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(59, 130, 246) // Blue
		pdf.SetX(15)
		pdf.Cell(0, 8, "Mapa de Calor por Departamento")
		pdf.Ln(10)
		
		// Interpretive analysis paragraph
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(15)
		heatmapAnalysis := "El mapa de calor por departamento permite identificar areas organizacionales con mayor concentracion "+
			"de riesgos psicosociales. Esta informacion es crucial para dirigir recursos y estrategias de prevencion de manera "+
			"efectiva y priorizada."
		pdf.MultiCell(180, 5, heatmapAnalysis, "", "L", false)
		pdf.Ln(6)

		// Group by department
		deptMap := make(map[string]map[string]int64)
		for _, hd := range report.DepartmentHeatmap {
			if _, exists := deptMap[hd.Department]; !exists {
				deptMap[hd.Department] = make(map[string]int64)
			}
			deptMap[hd.Department][string(hd.RiskLevel)] = hd.Count
		}

		departments := make([]string, 0, len(deptMap))
		for dept := range deptMap {
			departments = append(departments, dept)
		}
		sort.Strings(departments)

		// Professional table header with light blue background
		pdf.SetFillColor(230, 240, 255) // Light blue background
		pdf.SetFont("Arial", "B", 9)
		pdf.SetTextColor(30, 41, 59) // Dark text for contrast
		headerHeight := 12.0
		pdf.SetX(15)
		pdf.CellFormat(55, headerHeight, "Departamento", "1", 0, "L", true, 0, "")
		pdf.CellFormat(25, headerHeight, "Nulo", "1", 0, "R", true, 0, "") // Right-aligned for numbers
		pdf.CellFormat(25, headerHeight, "Bajo", "1", 0, "R", true, 0, "")
		pdf.CellFormat(25, headerHeight, "Medio", "1", 0, "R", true, 0, "")
		pdf.CellFormat(25, headerHeight, "Alto", "1", 0, "R", true, 0, "")
		pdf.CellFormat(25, headerHeight, "Muy Alto", "1", 1, "R", true, 0, "")

		// Rows
		pdf.SetFont("Arial", "", 9)
		fill := false
		riskLevels := []string{"nulo", "bajo", "medio", "alto", "muy_alto"}

		for _, dept := range departments {
			if pdf.GetY() > 250 {
				pdf.AddPage()
				// Redraw header with light blue background
				pdf.SetFillColor(230, 240, 255) // Light blue background
				pdf.SetFont("Arial", "B", 9)
				pdf.SetTextColor(30, 41, 59)
				pdf.SetX(15)
				pdf.CellFormat(55, headerHeight, "Departamento", "1", 0, "L", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Nulo", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Bajo", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Medio", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Alto", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Muy Alto", "1", 1, "R", true, 0, "")
				pdf.SetFont("Arial", "", 9)
				fill = false
			}

			rowY := pdf.GetY()
			rowHeight := 10.0

			if fill {
				pdf.SetFillColor(249, 250, 251) // Alternating row color
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			deptName := dept
			if len(deptName) > 22 {
				deptName = deptName[:19] + "..."
			}

			// Department name - left-aligned
			pdf.SetTextColor(30, 41, 59)
			pdf.SetFont("Arial", "", 9)
			pdf.SetX(15)
			pdf.CellFormat(55, rowHeight, deptName, "LR", 0, "L", fill, 0, "")

			// Counts - right-aligned for numbers with heat map colors
			for _, level := range riskLevels {
				count := deptMap[dept][level]
				countStr := "-"
				cellFillColor := fill // Preserve row fill state
				
				if count > 0 {
					countStr = fmt.Sprintf("%d", count)
					// Apply heat map color gradients based on risk level
					switch level {
					case "nulo":
						// Very light green (almost white) for no risk
						pdf.SetFillColor(245, 255, 245)
						pdf.SetTextColor(21, 128, 61)
						cellFillColor = true
					case "bajo":
						// Green (#28A745 = 40, 167, 69) background for low risk
						pdf.SetFillColor(40, 167, 69)
						pdf.SetTextColor(255, 255, 255) // White text for contrast
						cellFillColor = true
					case "medio":
						// Yellow/amber for medium risk
						pdf.SetFillColor(255, 235, 59)
						pdf.SetTextColor(30, 41, 59)
						cellFillColor = true
					case "alto":
						// Orange for high risk
						pdf.SetFillColor(255, 152, 0)
						pdf.SetTextColor(255, 255, 255)
						cellFillColor = true
					case "muy_alto":
						// Red for very high risk
						pdf.SetFillColor(244, 67, 54)
						pdf.SetTextColor(255, 255, 255)
						cellFillColor = true
					default:
						pdf.SetTextColor(30, 41, 59)
					}
				} else {
					pdf.SetTextColor(100, 116, 139)
				}
				pdf.CellFormat(25, rowHeight, countStr, "LR", 0, "R", cellFillColor, 0, "")
				// Reset colors
				pdf.SetTextColor(30, 41, 59)
				if fill {
					pdf.SetFillColor(249, 250, 251)
				} else {
					pdf.SetFillColor(255, 255, 255)
				}
			}
			
			// Draw bottom border for row separation (1px solid #E5E7EB = 229, 231, 235)
			pdf.SetDrawColor(229, 231, 235)
			pdf.SetLineWidth(0.5)
			pdf.Line(15, rowY+rowHeight, 15+55+25*5, rowY+rowHeight)
			
			pdf.Ln(rowHeight)
			fill = !fill
		}
	}

	// ========== DEMOGRAPHIC SUMMARY ==========
	hasDemo := len(report.AgeDistribution) > 0 ||
		len(report.ShiftTypeDistribution) > 0 ||
		len(report.ExperienceDistribution) > 0 ||
		len(report.MaritalStatusDistribution) > 0

	if hasDemo {
		pdf.AddPage()
		// Subtitle
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(59, 130, 246) // Blue
		pdf.SetX(15)
		pdf.Cell(0, 8, "Analisis Demografico")
		pdf.Ln(10)
		
		// Interpretive analysis paragraph
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(15)
		demoAnalysis := "El analisis demografico proporciona insights sobre la composicion de la fuerza laboral, "+
			"permitiendo identificar patrones y desarrollar estrategias personalizadas de prevencion de riesgos psicosociales."
		pdf.MultiCell(180, 5, demoAnalysis, "", "L", false)
		pdf.Ln(6)

		chartWidth := 85.0
		chartHeight := 40.0
		chartSpacing := 10.0

		// Row 1: Age Distribution and Shift Type - Horizontal bar charts
		startY := pdf.GetY()

		if len(report.AgeDistribution) > 0 {
			drawDemoHorizontalBar(pdf, 15, startY, chartWidth, chartHeight, "Distribucion por Edad", report.AgeDistribution)
		}

		if len(report.ShiftTypeDistribution) > 0 {
			drawDemoHorizontalBar(pdf, 15+chartWidth+chartSpacing, startY, chartWidth, chartHeight, "Distribucion por Turno", report.ShiftTypeDistribution)
		}

		pdf.SetY(startY + chartHeight + 12)

		// Row 2: Experience and Marital Status - Horizontal bar charts
		startY = pdf.GetY()

		if len(report.ExperienceDistribution) > 0 {
			drawDemoHorizontalBar(pdf, 15, startY, chartWidth, chartHeight, "Experiencia Laboral", report.ExperienceDistribution)
		}

		if len(report.MaritalStatusDistribution) > 0 {
			drawDemoHorizontalBar(pdf, 15+chartWidth+chartSpacing, startY, chartWidth, chartHeight, "Estado Civil", report.MaritalStatusDistribution)
		}

		pdf.SetY(startY + chartHeight + 12)
	}

	// ========== PAGE 4: RECOMMENDATIONS AND ACTION PLAN ==========
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(59, 130, 246) // Blue
	pdf.SetX(15)
	pdf.Cell(0, 8, "Recomendaciones y Plan de Accion")
	pdf.Ln(10)
	
	// Numbered list of actionable items
	actionItems := []string{
		"Incrementar la participacion al 100%% para Q2 2026 mediante campanas lideradas por RRHH, incluyendo comunicaciones dirigidas y seguimiento individual",
		"Implementar capacitacion continua sobre factores psicosociales para supervisores y personal de RRHH, con enfoque en identificacion temprana y prevencion",
		"Establecer comite de prevencion de riesgos psicosociales con representacion de empleados y directivos, con reuniones trimestrales",
		"Desarrollar politicas de conciliacion trabajo-familia con horarios flexibles y apoyo para el equilibrio laboral",
		"Crear canales de comunicacion abiertos y confidenciales para reportar situaciones de riesgo, con respuesta garantizada en 48 horas",
		"Implementar programas de reconocimiento y desarrollo profesional para mejorar el sentido de pertenencia y retencion del talento",
		"Realizar evaluaciones de seguimiento semestrales para monitorear la efectividad de las medidas implementadas",
		"Establecer indicadores de desempeno (KPIs) para medir el impacto de las acciones preventivas en la reduccion de riesgos",
	}
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(71, 85, 105)
	for i, item := range actionItems {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		pdf.SetX(20)
		pdf.Cell(5, 6, fmt.Sprintf("%d.", i+1))
		pdf.SetX(25)
		pdf.MultiCell(165, 5, item, "", "L", false)
		pdf.Ln(3)
	}

	// ========== PSYCHOSOCIAL RISK PREVENTION POLICY STATEMENT ==========
	if pdf.GetY() > 220 {
		pdf.AddPage()
	}
	
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(59, 130, 246)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Declaracion de Politica de Prevencion de Riesgos Psicosociales")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	// NOM-035 specific commitments
	policyText := fmt.Sprintf("%s se compromete a prevenir y controlar los factores de riesgo psicosocial en el ambiente de trabajo, "+
		"en estricto cumplimiento con la NOM-035-STPS-2018. Esta politica establece nuestros compromisos especificos conforme a la norma: "+
		"(1) Identificacion: Realizar evaluaciones periodicas para identificar factores de riesgo psicosocial en el ambiente laboral, "+
		"seguiendo los procedimientos establecidos en la NOM-035-STPS-2018. (2) Prevencion: Implementar medidas de prevencion y control para "+
		"reducir o eliminar los factores de riesgo identificados, con enfasis en factores organizacionales y del entorno. "+
		"(3) Promocion: Promover entornos organizacionales favorables mediante politicas de trabajo, comunicacion efectiva, "+
		"y reconocimiento del desempeno. (4) Difusion: Difundir los resultados de las evaluaciones y las medidas implementadas "+
		"a todos los niveles organizacionales, garantizando la confidencialidad y proteccion de datos personales conforme a la "+
		"Ley Federal de Proteccion de Datos Personales en Posesion de los Particulares (LFPDPPP). Referencia NOM-035: ",
		report.CompanyName)
	
	// Write policy text
	pdf.MultiCell(180, 5, policyText, "", "L", false)
	
	// Add hyperlink to NOM-035 reference
	nom035URL := "http://dof.gob.mx/normasOficiales/7544/stps2/stps2.htm"
	pdf.SetTextColor(59, 130, 246) // Blue color for link
	pdf.SetFont("Arial", "U", 10)  // Underlined for link appearance
	linkText := nom035URL
	pdf.SetX(15)
	pdf.Cell(0, 5, linkText)
	pdf.SetFont("Arial", "", 10)  // Reset font
	pdf.SetTextColor(71, 85, 105) // Reset text color

	// ========== LONG-TERM TRACKING NOTES ==========
	if pdf.GetY() > 230 {
		pdf.AddPage()
	}
	
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(59, 130, 246)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Notas de Seguimiento a Largo Plazo")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	// Calculate next biennial assessment date (2 years from current date)
	nextAssessmentDate := time.Now().AddDate(2, 0, 0)
	nextAssessmentStr := nextAssessmentDate.Format("Enero de 2006")
	
	trackingText := fmt.Sprintf("Para garantizar la efectividad continua de las medidas de prevencion conforme a la NOM-035-STPS-2018: "+
		"(1) Revaluaciones bienales programadas para %s, con preparacion iniciando 3 meses antes. "+
		"(2) Mantener registros detallados de incidentes, acciones correctivas y su efectividad en sistema documentado. "+
		"(3) Comparar resultados entre periodos para identificar tendencias y patrones de riesgo. "+
		"(4) Ajustar estrategias basadas en datos cuantitativos y feedback cualitativo del personal. "+
		"(5) Documentar mejoras implementadas y su impacto medible en la reduccion de riesgos psicosociales.",
		nextAssessmentStr)
	pdf.MultiCell(180, 5, trackingText, "", "L", false)

	// ========== EMPLOYEE ENGAGEMENT TIPS ==========
	if pdf.GetY() > 230 {
		pdf.AddPage()
	}
	
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(59, 130, 246)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Consejos para el Compromiso de los Empleados")
	pdf.Ln(10)
	
	engagementTips := []string{
		"Compartir este reporte a traves de portales internos con canales de retroalimentacion habilitados para comentarios y sugerencias",
		"Organizar sesiones informativas departamentales para explicar los resultados y las medidas de prevencion implementadas",
		"Establecer buzones de sugerencias digitales para que los empleados propongan mejoras al ambiente laboral",
		"Fomentar la participacion activa en programas de bienestar y prevencion mediante incentivos y reconocimientos",
		"Reconocer y celebrar los logros individuales y de equipo que contribuyan a un ambiente laboral saludable",
		"Proporcionar oportunidades de desarrollo profesional y crecimiento mediante programas de capacitacion y mentoring",
		"Establecer espacios de dialogo mensuales entre empleados y directivos para discutir preocupaciones y propuestas",
		"Involucrar a los empleados en la toma de decisiones que les afectan directamente mediante comites representativos",
	}
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(71, 85, 105)
	for _, tip := range engagementTips {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		pdf.SetX(20)
		pdf.Cell(5, 6, "•")
		pdf.SetX(25)
		pdf.MultiCell(165, 5, tip, "", "L", false)
		pdf.Ln(2)
	}

	// ========== CONFIDENTIALITY NOTE BOX ==========
	if pdf.GetY() > 240 {
		pdf.AddPage()
	}

	pdf.Ln(10)
	pdf.SetFillColor(230, 240, 255) // Light blue background
	pdf.RoundedRect(15, pdf.GetY(), 180, 30, 4, "1234", "F")

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(22, pdf.GetY()+6)
	pdf.Cell(0, 0, "Nota de Confidencialidad")
	
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetXY(22, pdf.GetY()+8)
	pdf.MultiCell(165, 4, "Este documento contiene informacion confidencial protegida por la Ley Federal de Proteccion de Datos Personales. "+
		"Su distribucion no autorizada esta prohibida. Generado automaticamente por la plataforma Entorno35 conforme a la NOM-035-STPS-2018.", "", "L", false)

	// ========== LEGAL FOOTER ==========
	if pdf.GetY() > 240 {
		pdf.AddPage()
	}

	pdf.Ln(15)
	pdf.SetFillColor(241, 245, 249)
	pdf.RoundedRect(15, pdf.GetY(), 180, 25, 3, "1234", "F")

	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(22, pdf.GetY()+6)
	pdf.MultiCell(165, 4, "Este documento contiene informacion confidencial protegida por la Ley Federal de Proteccion de Datos Personales. Su distribucion no autorizada esta prohibida. Generado automaticamente por la plataforma Entorno35 conforme a la NOM-035-STPS-2018.", "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("error al generar PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// drawMetricCard draws a metric card
func drawMetricCard(pdf *gofpdf.Fpdf, x, y, w, h float64, value, label string, textR, textG, textB, bgR, bgG, bgB int) {
	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.RoundedRect(x, y, w, h, 4, "1234", "F")

	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(textR, textG, textB)
	valueWidth := pdf.GetStringWidth(value)
	pdf.SetXY(x+(w-valueWidth)/2, y+10)
	pdf.Cell(0, 0, value)

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(71, 85, 105)
	labelWidth := pdf.GetStringWidth(label)
	pdf.SetXY(x+(w-labelWidth)/2, y+h-10)
	pdf.Cell(0, 0, label)
}

// drawEnhancedMetricCard draws an enhanced metric card with large numbers (32pt) and subtitle (12pt)
func drawEnhancedMetricCard(pdf *gofpdf.Fpdf, x, y, w, h float64, value, label string, textR, textG, textB, bgR, bgG, bgB int) {
	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.RoundedRect(x, y, w, h, 4, "1234", "F")

	// Large bold number (32pt)
	pdf.SetFont("Arial", "B", 32)
	pdf.SetTextColor(textR, textG, textB)
	valueWidth := pdf.GetStringWidth(value)
	pdf.SetXY(x+(w-valueWidth)/2, y+18)
	pdf.Cell(0, 0, value)

	// Subtitle (12pt)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(71, 85, 105)
	labelWidth := pdf.GetStringWidth(label)
	pdf.SetXY(x+(w-labelWidth)/2, y+h-12)
	pdf.Cell(0, 0, label)
}

// drawDemoSummaryInline draws a demographic summary inline
func drawDemoSummaryInline(pdf *gofpdf.Fpdf, title string, data []domain.DemographicDistribution) {
	if len(data) == 0 {
		return
	}

	if pdf.GetY() > 260 {
		pdf.AddPage()
	}

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	pdf.Cell(50, 6, title)

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(30, 41, 59)

	items := make([]string, 0)
	for i, d := range data {
		if i >= 5 {
			break
		}
		items = append(items, fmt.Sprintf("%s: %d", d.Category, d.Count))
	}

	pdf.Cell(0, 6, strings.Join(items, " | "))
	pdf.Ln(10)
}

// drawDemoPieChart draws a demographic distribution as a visual chart with legend
func drawDemoPieChart(pdf *gofpdf.Fpdf, x, y, w, h float64, title string, data []domain.DemographicDistribution, colors []struct{ r, g, b int }) {
	if len(data) == 0 {
		return
	}

	// Card background
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(x, y, w, h, 4, "1234", "F")

	// Title
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(x+4, y+4)
	pdf.Cell(w-8, 5, title)

	// Calculate total
	var total int64
	for _, d := range data {
		total += d.Count
	}
	if total == 0 {
		return
	}

	// Draw horizontal bar segments
	barX := x + 4
	barY := y + 14
	barWidth := w - 8
	barHeight := 12.0

	// Draw background bar
	pdf.SetFillColor(226, 232, 240)
	pdf.RoundedRect(barX, barY, barWidth, barHeight, 3, "1234", "F")

	// Draw segments
	currentX := barX
	for i, d := range data {
		if i >= len(colors) {
			break
		}
		if d.Count == 0 {
			continue
		}

		segmentWidth := (float64(d.Count) / float64(total)) * barWidth
		if segmentWidth > 0.5 {
			pdf.SetFillColor(colors[i].r, colors[i].g, colors[i].b)
			if i == 0 {
				// First segment - round left corners
				pdf.RoundedRect(currentX, barY, segmentWidth, barHeight, 3, "13", "F")
			} else if i == len(data)-1 || currentX+segmentWidth >= barX+barWidth-1 {
				// Last segment - round right corners
				pdf.RoundedRect(currentX, barY, segmentWidth, barHeight, 3, "24", "F")
			} else {
				// Middle segment - no rounding
				pdf.Rect(currentX, barY, segmentWidth, barHeight, "F")
			}
			currentX += segmentWidth
		}
	}

	// Draw legend
	legendY := barY + barHeight + 4
	legendX := x + 4
	itemsPerRow := 2
	legendItemWidth := (w - 8) / float64(itemsPerRow)

	for i, d := range data {
		if i >= 6 || i >= len(colors) {
			break
		}

		row := i / itemsPerRow
		col := i % itemsPerRow
		itemX := legendX + float64(col)*legendItemWidth
		itemY := legendY + float64(row)*10

		// Color box
		pdf.SetFillColor(colors[i].r, colors[i].g, colors[i].b)
		pdf.Rect(itemX, itemY, 6, 6, "F")

		// Label
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetXY(itemX+8, itemY)

		// Truncate category name if too long
		categoryName := d.Category
		if len(categoryName) > 10 {
			categoryName = categoryName[:8] + ".."
		}

		percentage := float64(d.Count) / float64(total) * 100
		pdf.Cell(legendItemWidth-12, 6, fmt.Sprintf("%s: %.0f%%", categoryName, percentage))
	}
}

// drawDemoHorizontalBar draws a horizontal bar chart with blue bars (#007BFF) and white text labels inside
func drawDemoHorizontalBar(pdf *gofpdf.Fpdf, x, y, w, h float64, title string, data []domain.DemographicDistribution) {
	if len(data) == 0 {
		return
	}

	// Card background
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(x, y, w, h, 4, "1234", "F")

	// Title
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(x+4, y+4)
	pdf.Cell(w-8, 5, title)

	// Calculate total
	var total int64
	for _, d := range data {
		total += d.Count
	}
	if total == 0 {
		return
	}

	// Draw horizontal bar - find the largest category (100%)
	maxCount := int64(0)
	maxCategory := ""
	for _, d := range data {
		if d.Count > maxCount {
			maxCount = d.Count
			maxCategory = d.Category
		}
	}

	// Bar dimensions
	barX := x + 4
	barY := y + 14
	barWidth := w - 8
	barHeight := 20.0

	// Draw background bar
	pdf.SetFillColor(226, 232, 240)
	pdf.RoundedRect(barX, barY, barWidth, barHeight, 3, "1234", "F")

	// Draw blue bar at 100% (#007BFF = 0, 123, 255)
	pdf.SetFillColor(0, 123, 255)
	pdf.RoundedRect(barX, barY, barWidth, barHeight, 3, "1234", "F")

	// White text label inside bar
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(255, 255, 255)
	labelText := fmt.Sprintf("%s: 100%%", maxCategory)
	labelWidth := pdf.GetStringWidth(labelText)
	labelX := barX + (barWidth-labelWidth)/2
	pdf.SetXY(labelX, barY+6)
	pdf.Cell(0, 0, labelText)
}

