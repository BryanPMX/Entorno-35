package services

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
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

// getFontPath resolves the absolute path to a font file
// It tries multiple strategies to find the fonts directory:
// 1. Relative to project root (by finding go.mod) - PRIORITY
// 2. Relative to current working directory (fonts/filename)
// 3. Relative to executable location
func getFontPath(filename string) (string, error) {
	var attemptedPaths []string

	// PRIORITY: Try from project root by finding go.mod first
	// This works regardless of where the server is started from
	cwd, err := os.Getwd()
	if err == nil {
		dir := cwd
		for {
			goModPath := filepath.Join(dir, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				// Found project root - construct absolute path directly
				fontPath := filepath.Join(dir, "fonts", filename)
				// Since dir comes from os.Getwd(), it should already be absolute
				// But ensure it's absolute and clean
				if !filepath.IsAbs(fontPath) {
					var absErr error
					fontPath, absErr = filepath.Abs(fontPath)
					if absErr != nil {
						return "", fmt.Errorf("failed to get absolute path: %v", absErr)
					}
				}
				// Clean the path to remove any redundant separators
				fontPath = filepath.Clean(fontPath)
				attemptedPaths = append(attemptedPaths, fontPath)

				// Verify the file exists and path is absolute
				if _, err := os.Stat(fontPath); err == nil {
					// Double-check it's absolute before returning
					if !filepath.IsAbs(fontPath) {
						return "", fmt.Errorf("path is not absolute after cleaning: %s", fontPath)
					}
					return fontPath, nil
				}
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break // Reached filesystem root
			}
			dir = parent
		}
	}

	// Try relative path (works when running from project root)
	relativePath := filepath.Join("fonts", filename)
	attemptedPaths = append(attemptedPaths, relativePath)
	if _, err := os.Stat(relativePath); err == nil {
		absPath, absErr := filepath.Abs(relativePath)
		if absErr != nil {
			return "", fmt.Errorf("failed to get absolute path for %s: %v", relativePath, absErr)
		}
		return absPath, nil
	}

	// Try relative to executable location
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		fontPath := filepath.Join(execDir, "fonts", filename)
		attemptedPaths = append(attemptedPaths, fontPath)
		// Ensure absolute path
		if !filepath.IsAbs(fontPath) {
			var absErr error
			fontPath, absErr = filepath.Abs(fontPath)
			if absErr != nil {
				attemptedPaths = append(attemptedPaths, fmt.Sprintf("%s (abs failed: %v)", fontPath, absErr))
			} else {
				attemptedPaths[len(attemptedPaths)-1] = fontPath
			}
		}
		if _, err := os.Stat(fontPath); err == nil {
			return fontPath, nil
		}
		// Try going up from bin/ directory
		if filepath.Base(execDir) == "bin" {
			projectRoot := filepath.Dir(execDir)
			fontPath := filepath.Join(projectRoot, "fonts", filename)
			attemptedPaths = append(attemptedPaths, fontPath)
			// Ensure absolute path
			if !filepath.IsAbs(fontPath) {
				var absErr error
				fontPath, absErr = filepath.Abs(fontPath)
				if absErr != nil {
					attemptedPaths = append(attemptedPaths, fmt.Sprintf("%s (abs failed: %v)", fontPath, absErr))
				} else {
					attemptedPaths[len(attemptedPaths)-1] = fontPath
				}
			}
			if _, err := os.Stat(fontPath); err == nil {
				return fontPath, nil
			}
		}
	}

	// Get current working directory for error message
	cwdStr := "unknown"
	if cwd, err := os.Getwd(); err == nil {
		cwdStr = cwd
	}

	// Return error with all attempted paths
	return "", fmt.Errorf("font file '%s' not found. Attempted paths: %v. Current working directory: %s",
		filename, attemptedPaths, cwdStr)
}

// GenerateIndividualReportPDF creates a professional PDF report for an individual assessment
func (s *ReportPDFService) GenerateIndividualReportPDF(report *domain.IndividualReportDTO, companyName string) ([]byte, error) {
	// Load UTF-8 fonts for proper Spanish character support
	regularFontPath, err := getFontPath("DejaVuSans.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to locate regular font: %v", err)
	}
	boldFontPath, err := getFontPath("DejaVuSans-Bold.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to locate bold font: %v", err)
	}

	// Verify paths are absolute before using
	if !filepath.IsAbs(regularFontPath) {
		return nil, fmt.Errorf("regular font path is not absolute: %s", regularFontPath)
	}
	if !filepath.IsAbs(boldFontPath) {
		return nil, fmt.Errorf("bold font path is not absolute: %s", boldFontPath)
	}

	// Clean paths to ensure they're properly formatted
	regularFontPath = filepath.Clean(regularFontPath)
	boldFontPath = filepath.Clean(boldFontPath)

	// Get font directory and filenames
	// Pass font directory as 4th parameter to gofpdf.New() to avoid path manipulation issues
	fontDir := filepath.Dir(regularFontPath)
	regularFontName := filepath.Base(regularFontPath)
	boldFontName := filepath.Base(boldFontPath)

	// Ensure font directory path uses forward slashes and is absolute
	fontDir = filepath.ToSlash(fontDir)
	if !filepath.IsAbs(fontDir) {
		var absErr error
		fontDir, absErr = filepath.Abs(fontDir)
		if absErr != nil {
			return nil, fmt.Errorf("failed to get absolute font directory: %v", absErr)
		}
		fontDir = filepath.ToSlash(fontDir)
	}
	if len(fontDir) > 0 && fontDir[0] != '/' {
		fontDir = "/" + fontDir
	}

	// Initialize PDF with font directory as 4th parameter
	pdf := gofpdf.New("P", "mm", "A4", fontDir)

	// Add fonts using just the filename (relative to font directory)
	pdf.AddUTF8Font("DejaVu", "", regularFontName)
	pdf.AddUTF8Font("DejaVu", "B", boldFontName)

	pdf.SetAutoPageBreak(true, 30.0)
	pdf.AddPage()

	// Professional footer with confidential notice
	pdf.SetFooterFunc(func() {
		// Confidential notice at bottom
		pdf.SetY(-20)
		pdf.SetFont("DejaVu", "", 7)
		pdf.SetTextColor(100, 116, 139)
		confidentialText := "Este documento contiene informacion confidencial protegida por la Ley Federal de Proteccion de Datos Personales. Su distribucion no autorizada esta prohibida."
		pdf.MultiCell(180, 3, confidentialText, "", "C", false)

		// Page info above confidential notice
		pdf.SetY(-12)
		pdf.SetFont("DejaVu", "", 8)
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

	// Header text - use UTF-8 font for Spanish characters
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("DejaVu", "B", 18)
	pdf.SetXY(15, 10)
	pdf.Cell(0, 0, "Reporte de Riesgo Psicosocial")

	pdf.SetFont("DejaVu", "", 11)
	pdf.SetXY(15, 20)
	pdf.Cell(0, 0, fmt.Sprintf("NOM-035-STPS-2018 | Guia %s", string(report.GuideType)))

	// Confidential badge - properly centered
	pdf.SetFillColor(239, 68, 68) // Red-500
	badgeX := 155.0
	badgeY := 8.0
	badgeWidth := 42.0
	badgeHeight := 8.0
	pdf.RoundedRect(badgeX, badgeY, badgeWidth, badgeHeight, 2, "1234", "F")
	pdf.SetFont("DejaVu", "B", 8)
	pdf.SetTextColor(255, 255, 255)
	// Use CellFormat with badgeHeight to properly center text vertically
	pdf.SetXY(badgeX, badgeY)
	pdf.CellFormat(badgeWidth, badgeHeight, "CONFIDENCIAL", "", 0, "CM", false, 0, "")

	// Date
	pdf.SetFont("DejaVu", "", 9)
	pdf.SetXY(155, 22)
	pdf.Cell(0, 0, time.Now().Format("02 de Enero, 2006"))

	pdf.SetY(48)

	// ========== EMPLOYEE INFO CARD ==========
	pdf.SetFillColor(248, 250, 252) // Slate-50
	pdf.RoundedRect(15, pdf.GetY(), 180, 32, 4, "1234", "F")

	infoY := pdf.GetY() + 6
	pdf.SetFont("DejaVu", "B", 12)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(22, infoY)
	pdf.Cell(0, 0, report.StaffName)

	pdf.SetFont("DejaVu", "", 10)
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
	pdf.SetFont("DejaVu", "B", 36)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(28, scoreY)
	pdf.Cell(0, 0, fmt.Sprintf("%.0f", report.TotalScore))

	// Max score
	maxScore := report.TotalMaxScore
	if maxScore == 0 {
		maxScore = 184
	}
	pdf.SetFont("DejaVu", "", 14)
	pdf.SetTextColor(148, 163, 184)
	scoreWidth := pdf.GetStringWidth(fmt.Sprintf("%.0f", report.TotalScore))
	pdf.SetXY(28+scoreWidth+3, scoreY+8)
	pdf.Cell(0, 0, fmt.Sprintf("/ %.0f puntos", maxScore))

	// Percentage
	percentage := (report.TotalScore / maxScore) * 100
	if percentage > 100 {
		percentage = 100
	}
	pdf.SetFont("DejaVu", "", 10)
	pdf.SetXY(28, scoreY+18)
	pdf.Cell(0, 0, fmt.Sprintf("%.1f%% del maximo posible", percentage))

	// Risk badge on right
	riskBadgeX := 130.0
	riskBadgeY := scoreY + 5

	pdf.SetFillColor(riskConfig.bgR, riskConfig.bgG, riskConfig.bgB)
	pdf.RoundedRect(riskBadgeX, riskBadgeY, 55, 28, 3, "1234", "F")

	pdf.SetFont("DejaVu", "B", 12)
	pdf.SetTextColor(riskConfig.textR, riskConfig.textG, riskConfig.textB)
	riskLabel := formatRiskLevel(riskLevel)
	labelWidth := pdf.GetStringWidth(riskLabel)
	pdf.SetXY(riskBadgeX+(55-labelWidth)/2, riskBadgeY+8)
	pdf.Cell(0, 0, riskLabel)

	pdf.SetFont("DejaVu", "", 9)
	pdf.SetXY(riskBadgeX, riskBadgeY+18)
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
		pdf.SetFont("DejaVu", "B", 10)
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

			// Professional numbered list format
			pdf.SetX(20)
			pdf.SetFont("DejaVu", "B", 11)
			pdf.SetTextColor(59, 130, 246) // Blue for numbers
			pdf.Cell(8, 6, fmt.Sprintf("%d.", i+1))

			// Recommendation text with better spacing
			pdf.SetFont("DejaVu", "", 10)
			pdf.SetTextColor(51, 65, 85)
			pdf.SetX(32)
			pdf.MultiCell(158, 5, rec, "", "L", false)

			pdf.Ln(6) // Extra spacing between recommendations
		}
	} else {
		// Positive result for low risk
		if pdf.GetY() > 220 {
			pdf.AddPage()
		}

		pdf.Ln(8)
		pdf.SetFillColor(220, 252, 231)
		pdf.RoundedRect(15, pdf.GetY(), 180, 30, 4, "1234", "F")

		pdf.SetFont("DejaVu", "B", 11)
		pdf.SetTextColor(21, 128, 61)
		pdf.SetXY(22, pdf.GetY()+8)
		pdf.Cell(0, 0, "Resultado Favorable")

		pdf.SetFont("DejaVu", "", 10)
		pdf.SetXY(22, pdf.GetY()+4)
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
			pdf.SetFont("DejaVu", "B", 10)
			pdf.SetTextColor(30, 41, 59)
			pdf.SetX(15)
			pdf.CellFormat(180, 8, categoryName, "", 1, "L", true, 0, "")
			pdf.Ln(2)

			// Questions table header
			pdf.SetFillColor(248, 250, 252)
			pdf.SetFont("DejaVu", "B", 8)
			pdf.SetTextColor(71, 85, 105)
			pdf.SetX(15)
			pdf.CellFormat(10, 7, "#", "", 0, "C", true, 0, "")
			pdf.CellFormat(115, 7, "Pregunta", "", 0, "L", true, 0, "")
			pdf.CellFormat(20, 7, "Resp.", "", 0, "C", true, 0, "")
			pdf.CellFormat(20, 7, "Punt.", "", 0, "C", true, 0, "")
			pdf.CellFormat(15, 7, "Max", "", 1, "C", true, 0, "")

			// Questions rows
			pdf.SetFont("DejaVu", "", 8)
			fill := false
			for _, q := range questions {
				if pdf.GetY() > 265 {
					pdf.AddPage()
					// Redraw category header
					pdf.SetFillColor(241, 245, 249)
					pdf.SetFont("DejaVu", "B", 10)
					pdf.SetTextColor(30, 41, 59)
					pdf.SetX(15)
					pdf.CellFormat(180, 8, categoryName+" (continuacion)", "", 1, "L", true, 0, "")
					pdf.Ln(2)
					// Redraw table header
					pdf.SetFillColor(248, 250, 252)
					pdf.SetFont("DejaVu", "B", 8)
					pdf.SetTextColor(71, 85, 105)
					pdf.SetX(15)
					pdf.CellFormat(10, 7, "#", "", 0, "C", true, 0, "")
					pdf.CellFormat(115, 7, "Pregunta", "", 0, "L", true, 0, "")
					pdf.CellFormat(20, 7, "Resp.", "", 0, "C", true, 0, "")
					pdf.CellFormat(20, 7, "Punt.", "", 0, "C", true, 0, "")
					pdf.CellFormat(15, 7, "Max", "", 1, "C", true, 0, "")
					pdf.SetFont("DejaVu", "", 8)
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

	// ========== PAGE 5: OFFICIAL DOCUMENTATION APPENDIX ==========
	pdf.AddPage()

	// Header
	pdf.SetFillColor(15, 23, 42) // Dark blue #1A1F36 (26, 31, 54)
	pdf.Rect(0, 0, 210, 35, "F")

	// Accent stripe
	pdf.SetFillColor(59, 130, 246) // Blue-500
	pdf.Rect(0, 35, 210, 3, "F")

	// Header text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("DejaVu", "B", 18)
	pdf.SetXY(15, 10)
	pdf.Cell(0, 0, "Apéndice de Documentación Oficial")

	pdf.SetFont("DejaVu", "", 11)
	pdf.SetXY(15, 20)
	pdf.Cell(0, 0, "NOM-035-STPS-2018 | Referencias y Glosario")

	// Confidential badge
	pdf.SetFillColor(255, 77, 79) // Red #FF4D4F
	badgeX = 155.0
	badgeY = 8.0
	badgeWidth = 42.0
	badgeHeight = 8.0
	pdf.RoundedRect(badgeX, badgeY, badgeWidth, badgeHeight, 2, "1234", "F")
	pdf.SetFont("DejaVu", "B", 8)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(badgeX, badgeY)
	pdf.CellFormat(badgeWidth, badgeHeight, "CONFIDENCIAL", "", 0, "CM", false, 0, "")

	pdf.SetY(48)

	// ========== GLOSSARY ==========
	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Glosario de Términos")
	pdf.Ln(10)

	glossary := []struct {
		term       string
		definition string
	}{
		{
			term:       "Riesgo Psicosocial",
			definition: "Factores que pueden provocar trastornos de ansiedad, estrés, depresión u otros problemas de salud mental en los trabajadores.",
		},
		{
			term:       "NOM-035-STPS-2018",
			definition: "Norma Oficial Mexicana que establece los elementos para identificar, analizar y prevenir los factores de riesgo psicosocial en los centros de trabajo.",
		},
		{
			term:       "Nivel de Riesgo",
			definition: "Clasificación del riesgo psicosocial según los umbrales establecidos en la NOM-035: Nulo (<20), Bajo (20-<45), Medio (45-<70), Alto (70-<90), Muy Alto (≥90).",
		},
		{
			term:       "Factores de Riesgo Organizacional",
			definition: "Condiciones peligrosas o nocivas inherentes al trabajo, relacionadas con la organización del mismo.",
		},
		{
			term:       "Atención Médica Ocupacional",
			definition: "Servicios médicos especializados para trabajadores expuestos a factores de riesgo psicosocial.",
		},
	}

	pdf.SetFont("DejaVu", "", 10)
	for _, item := range glossary {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		pdf.SetFont("DejaVu", "B", 10)
		pdf.SetTextColor(30, 41, 59)
		pdf.SetX(20)
		pdf.Cell(0, 6, item.term+":")
		pdf.Ln(8)
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(25)
		pdf.MultiCell(160, 4, item.definition, "", "L", false)
		pdf.Ln(4)
	}

	pdf.Ln(8)

	// ========== OFFICIAL LINKS ==========
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}

	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Enlaces a Documentación Oficial")
	pdf.Ln(10)

	// Single NOM-035 link
	nom035Title := "NOM-035-STPS-2018 - Texto Completo"
	nom035URL := "https://www.gob.mx/stps/articulos/norma-oficial-mexicana-nom-035-stps-2018-factores-de-riesgo-psicosocial-en-el-trabajo-identificacion-analisis-y-prevencion"

	pdf.SetFont("DejaVu", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(20)
	pdf.Cell(5, 6, "1.")
	pdf.SetX(25)
	pdf.SetTextColor(59, 130, 246) // Blue for links
	pdf.SetFont("DejaVu", "U", 10) // Underlined
	pdf.MultiCell(165, 5, nom035Title, "", "L", false)
	pdf.SetFont("DejaVu", "", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetX(25)
	pdf.MultiCell(165, 4, nom035URL, "", "L", false)
	pdf.Ln(3)
	pdf.SetFont("DejaVu", "", 10)

	pdf.Ln(8)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("error al generar PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// drawModernSection draws a section header with improved spacing
func drawModernSection(pdf *gofpdf.Fpdf, title, subtitle string) {
	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetX(15)
	pdf.Cell(0, 8, title)
	pdf.Ln(10) // Increased from 6 to 10 for more space between title and subtitle

	pdf.SetFont("DejaVu", "", 10)
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
	pdf.SetFont("DejaVu", "B", 9)
	pdf.SetTextColor(71, 85, 105)

	pdf.SetX(15)
	pdf.CellFormat(70, 10, "Elemento", "", 0, "L", true, 0, "")
	pdf.CellFormat(52, 10, "Puntuacion", "", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Nivel de Riesgo", "", 1, "C", true, 0, "")

	// Rows
	pdf.SetFont("DejaVu", "", 9)
	fill := false

	for _, k := range keys {
		if pdf.GetY() > 250 {
			pdf.AddPage()
			// Redraw header
			pdf.SetFillColor(241, 245, 249)
			pdf.SetFont("DejaVu", "B", 9)
			pdf.SetTextColor(71, 85, 105)
			pdf.SetX(15)
			pdf.CellFormat(70, 10, "Elemento", "", 0, "L", true, 0, "")
			pdf.CellFormat(52, 10, "Puntuacion", "", 0, "C", true, 0, "")
			pdf.CellFormat(45, 10, "Nivel de Riesgo", "", 1, "C", true, 0, "")
			pdf.SetFont("DejaVu", "", 9)
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
		maxItemWidth := 65.0
		if pdf.GetStringWidth(itemText) > maxItemWidth {
			for len(itemText) > 0 && pdf.GetStringWidth(itemText+"...") > maxItemWidth {
				itemText = itemText[:len(itemText)-1]
			}
			itemText += "..."
		}

		pdf.SetTextColor(30, 41, 59)
		pdf.SetX(15)
		pdf.CellFormat(70, 9, itemText, "", 0, "L", fill, 0, "")

		// Score
		pdf.SetTextColor(71, 85, 105)
		scoreText := fmt.Sprintf("%.1f / %.0f", score, maxScore)
		// Truncate score text if too long
		if pdf.GetStringWidth(scoreText) > 50 {
			scoreText = fmt.Sprintf("%.0f/%.0f", score, maxScore)
		}
		pdf.CellFormat(52, 9, scoreText, "", 0, "C", fill, 0, "")

		// Risk level badge
		riskLabel := formatRiskLevel(riskLevel)
		// Truncate risk label if too long to prevent cutoff
		maxRiskLabelWidth := 43.0
		if pdf.GetStringWidth(riskLabel) > maxRiskLabelWidth {
			for len(riskLabel) > 0 && pdf.GetStringWidth(riskLabel) > maxRiskLabelWidth {
				riskLabel = riskLabel[:len(riskLabel)-1]
			}
		}
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
	// Load UTF-8 fonts for proper Spanish character support
	regularFontPath, err := getFontPath("DejaVuSans.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to locate regular font: %v", err)
	}
	boldFontPath, err := getFontPath("DejaVuSans-Bold.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to locate bold font: %v", err)
	}

	// Verify paths are absolute before using
	if !filepath.IsAbs(regularFontPath) {
		return nil, fmt.Errorf("regular font path is not absolute: %s", regularFontPath)
	}
	if !filepath.IsAbs(boldFontPath) {
		return nil, fmt.Errorf("bold font path is not absolute: %s", boldFontPath)
	}

	// Clean paths to ensure they're properly formatted
	regularFontPath = filepath.Clean(regularFontPath)
	boldFontPath = filepath.Clean(boldFontPath)

	// Get font directory and filenames
	// Pass font directory as 4th parameter to gofpdf.New() to avoid path manipulation issues
	fontDir := filepath.Dir(regularFontPath)
	regularFontName := filepath.Base(regularFontPath)
	boldFontName := filepath.Base(boldFontPath)

	// Ensure font directory path uses forward slashes and is absolute
	fontDir = filepath.ToSlash(fontDir)
	if !filepath.IsAbs(fontDir) {
		var absErr error
		fontDir, absErr = filepath.Abs(fontDir)
		if absErr != nil {
			return nil, fmt.Errorf("failed to get absolute font directory: %v", absErr)
		}
		fontDir = filepath.ToSlash(fontDir)
	}
	if len(fontDir) > 0 && fontDir[0] != '/' {
		fontDir = "/" + fontDir
	}

	// Initialize PDF with font directory as 4th parameter
	pdf := gofpdf.New("P", "mm", "A4", fontDir)

	// Add fonts using just the filename (relative to font directory)
	pdf.AddUTF8Font("DejaVu", "", regularFontName)
	pdf.AddUTF8Font("DejaVu", "B", boldFontName)

	pdf.SetAutoPageBreak(true, 30.0)
	pdf.AddPage()

	// Professional footer with confidential notice
	pdf.SetFooterFunc(func() {
		// Confidential notice at bottom
		pdf.SetY(-20)
		pdf.SetFont("DejaVu", "", 7)
		pdf.SetTextColor(100, 116, 139)
		confidentialText := "Este documento contiene informacion confidencial protegida por la Ley Federal de Proteccion de Datos Personales. Su distribucion no autorizada esta prohibida."
		pdf.MultiCell(180, 3, confidentialText, "", "C", false)

		// Page info above confidential notice
		pdf.SetY(-12)
		pdf.SetFont("DejaVu", "", 8)
		pdf.SetTextColor(128, 128, 128)
		footerText := fmt.Sprintf("Entorno35 | Reporte General NOM-035 | %s | Pagina %d", time.Now().Format("02/01/2006"), pdf.PageNo())
		pdf.CellFormat(0, 10, footerText, "", 0, "C", false, 0, "")
	})

	// ========== HEADER SECTION ==========
	// Modern gradient header - matching individual report format
	pdf.SetFillColor(15, 23, 42) // Slate-900
	pdf.Rect(0, 0, 210, 35, "F")

	// Accent stripe
	pdf.SetFillColor(59, 130, 246) // Blue-500
	pdf.Rect(0, 35, 210, 3, "F")

	// Header text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("DejaVu", "B", 18)
	pdf.SetXY(15, 10)
	pdf.Cell(0, 0, "Reporte General de Cumplimiento")

	pdf.SetFont("DejaVu", "", 11)
	pdf.SetXY(15, 20)
	periodText := "Todos los Periodos"
	if report.Period != nil {
		periodText = fmt.Sprintf("Periodo %d", *report.Period)
	}
	pdf.Cell(0, 0, fmt.Sprintf("NOM-035-STPS-2018 | %s", periodText))

	// Confidential badge - properly centered
	pdf.SetFillColor(239, 68, 68) // Red-500
	badgeX := 155.0
	badgeY := 8.0
	badgeWidth := 42.0
	badgeHeight := 8.0
	pdf.RoundedRect(badgeX, badgeY, badgeWidth, badgeHeight, 2, "1234", "F")
	pdf.SetFont("DejaVu", "B", 8)
	pdf.SetTextColor(255, 255, 255)
	// Use CellFormat with badgeHeight to properly center text vertically
	pdf.SetXY(badgeX, badgeY)
	pdf.CellFormat(badgeWidth, badgeHeight, "CONFIDENCIAL", "", 0, "CM", false, 0, "")

	// Date
	pdf.SetFont("DejaVu", "", 9)
	pdf.SetXY(155, 22)
	pdf.Cell(0, 0, time.Now().Format("02 de Enero, 2006"))

	pdf.SetY(48)

	// ========== COMPANY INFO CARD ==========
	// Matching individual report card format
	pdf.SetFillColor(248, 250, 252) // Slate-50
	pdf.RoundedRect(15, pdf.GetY(), 180, 32, 4, "1234", "F")

	infoY := pdf.GetY() + 6
	pdf.SetFont("DejaVu", "B", 12)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(22, infoY)
	pdf.Cell(0, 0, report.CompanyName)

	pdf.SetFont("DejaVu", "", 10)
	pdf.SetTextColor(100, 116, 139)

	// Row 1: Company ID
	pdf.SetXY(22, infoY+10)
	companyIDStr := report.CompanyID.String()
	if len(companyIDStr) > 12 {
		companyIDStr = companyIDStr[:12] + "..."
	}
	pdf.Cell(0, 0, fmt.Sprintf("ID: %s", companyIDStr))

	// Row 2: Date
	pdf.SetXY(22, infoY+18)
	pdf.Cell(0, 0, fmt.Sprintf("Fecha: %s", time.Now().Format("02 de Enero, 2006")))

	pdf.SetY(infoY + 35)

	// ========== EXECUTIVE SUMMARY ==========
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}

	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Resumen Ejecutivo")
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 10)
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
	pdf.SetFont("DejaVu", "B", 16)
	pdf.SetTextColor(59, 130, 246) // Blue
	pdf.SetX(15)
	pdf.Cell(0, 8, "Metricas Clave")
	pdf.Ln(10)

	metricsY := pdf.GetY()
	metricWidth := 55.0
	metricHeight := 45.0
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
		230, 240, 255, // Light blue background
		"") // No icon

	// Metric 2: Completed - Light green (#E6FFE6 = 230, 255, 230)
	drawEnhancedMetricCard(pdf, startX+metricWidth+spacing, metricsY, metricWidth, metricHeight,
		fmt.Sprintf("%d", report.CompletedAssessments), "Evaluaciones Completas",
		30, 41, 59, // Dark text
		230, 255, 230, // Light green background
		"") // No icon

	// Metric 3: Participation - Light yellow (#FFFBE6 = 255, 251, 230)
	drawEnhancedMetricCard(pdf, startX+2*(metricWidth+spacing), metricsY, metricWidth, metricHeight,
		fmt.Sprintf("%.1f%%", report.ParticipationRate), "Tasa de Participacion",
		30, 41, 59, // Dark text
		255, 251, 230, // Light yellow background
		"") // No icon

	pdf.SetY(metricsY + metricHeight + 15)

	// ========== RISK DISTRIBUTION ==========
	// Always start on a new page (page 2) to ensure the chart doesn't get split
	pdf.AddPage()

	// Subtitle
	pdf.SetFont("DejaVu", "B", 16)
	pdf.SetTextColor(59, 130, 246) // Blue
	pdf.SetX(15)
	pdf.Cell(0, 8, "Distribucion de Riesgo")
	pdf.Ln(10)

	// Interpretive analysis paragraph
	pdf.SetFont("DejaVu", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	analysisText := "La distribucion de riesgo muestra la cantidad de evaluaciones clasificadas en cada nivel segun los " +
		"criterios de la NOM-035-STPS-2018. Esta informacion es fundamental para identificar areas de atencion prioritaria " +
		"y desarrollar estrategias de prevencion efectivas."
	pdf.MultiCell(180, 5, analysisText, "", "L", false)
	pdf.Ln(6)

	if len(report.RiskDistribution) > 0 {
		riskOrder := []string{"nulo", "bajo", "medio", "alto", "muy_alto"}
		riskMap := make(map[string]int64)
		totalCount := int64(0)

		// Calculate total count
		for _, rd := range report.RiskDistribution {
			riskMap[string(rd.RiskLevel)] = rd.Count
			totalCount += rd.Count
		}

		// Draw pie chart with legend
		chartX := 15.0
		chartY := pdf.GetY()
		chartRadius := 40.0
		chartCenterX := chartX + chartRadius
		chartCenterY := chartY + chartRadius

		// Draw pie chart
		drawRiskPieChart(pdf, chartCenterX, chartCenterY, chartRadius, riskMap, totalCount, riskOrder)

		// Draw legend with NOM-035 thresholds
		legendX := chartX + chartRadius*2 + 20
		legendY := chartY
		drawRiskLegend(pdf, legendX, legendY, riskMap, totalCount, riskOrder)

		// Move Y position after chart
		pdf.SetY(chartY + chartRadius*2 + 15)
	} else {
		pdf.SetFont("DejaVu", "", 10)
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
		pdf.SetFont("DejaVu", "B", 16)
		pdf.SetTextColor(59, 130, 246) // Blue
		pdf.SetX(15)
		pdf.Cell(0, 8, "Mapa de Calor por Departamento")
		pdf.Ln(10)

		// Interpretive analysis paragraph
		pdf.SetFont("DejaVu", "", 10)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(15)
		heatmapAnalysis := "El mapa de calor por departamento permite identificar areas organizacionales con mayor concentracion " +
			"de riesgos psicosociales. Esta informacion es crucial para dirigir recursos y estrategias de prevencion de manera " +
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
		pdf.SetFont("DejaVu", "B", 9)
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
		pdf.SetFont("DejaVu", "", 9)
		fill := false
		riskLevels := []string{"nulo", "bajo", "medio", "alto", "muy_alto"}

		for _, dept := range departments {
			if pdf.GetY() > 250 {
				pdf.AddPage()
				// Redraw header with light blue background
				pdf.SetFillColor(230, 240, 255) // Light blue background
				pdf.SetFont("DejaVu", "B", 9)
				pdf.SetTextColor(30, 41, 59)
				pdf.SetX(15)
				pdf.CellFormat(55, headerHeight, "Departamento", "1", 0, "L", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Nulo", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Bajo", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Medio", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Alto", "1", 0, "R", true, 0, "")
				pdf.CellFormat(25, headerHeight, "Muy Alto", "1", 1, "R", true, 0, "")
				pdf.SetFont("DejaVu", "", 9)
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
			pdf.SetFont("DejaVu", "", 9)
			pdf.SetX(15)
			pdf.CellFormat(55, rowHeight, deptName, "LR", 0, "L", fill, 0, "")

			// Counts - right-aligned for numbers with heat map colors
			for _, level := range riskLevels {
				count := deptMap[dept][level]
				countStr := "-"
				cellFillColor := fill // Preserve row fill state

				if count > 0 {
					countStr = fmt.Sprintf("%d", count)
					// Apply heat map color gradients based on risk level (intensified colors)
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
						// Yellow (#FFC107 = 255, 193, 7) for medium risk
						pdf.SetFillColor(255, 193, 7)
						pdf.SetTextColor(30, 41, 59)
						cellFillColor = true
					case "alto":
						// Red (#DC3545 = 220, 53, 69) for high risk
						pdf.SetFillColor(220, 53, 69)
						pdf.SetTextColor(255, 255, 255)
						cellFillColor = true
					case "muy_alto":
						// Red (#DC3545 = 220, 53, 69) for very high risk
						pdf.SetFillColor(220, 53, 69)
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
		pdf.SetFont("DejaVu", "B", 16)
		pdf.SetTextColor(59, 130, 246) // Blue
		pdf.SetX(15)
		pdf.Cell(0, 8, "Analisis Demografico")
		pdf.Ln(10)

		// Interpretive analysis paragraph
		pdf.SetFont("DejaVu", "", 10)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(15)
		demoAnalysis := "El analisis demografico proporciona estadisticas sobre la composicion de la fuerza laboral, " +
			"permitiendo identificar patrones y desarrollar estrategias personalizadas de prevencion de riesgos psicosociales."
		pdf.MultiCell(180, 5, demoAnalysis, "", "L", false)
		pdf.Ln(6)

		// 2x2 grid layout - all 4 charts on one page
		chartWidth := 80.0   // Slightly smaller for 2-column layout
		chartHeight := 65.0  // Adjusted height to fit 2 rows comfortably
		chartSpacingX := 15.0 // Horizontal spacing between columns
		chartSpacingY := 15.0 // Vertical spacing between rows

		// Calculate positions for 2x2 grid
		leftX := 15.0
		rightX := leftX + chartWidth + chartSpacingX
		startY := pdf.GetY()

		// Top row: Age and Experience distributions
		if len(report.AgeDistribution) > 0 {
			drawDemoDonutChart(pdf, leftX, startY, chartWidth, chartHeight, "Distribucion por Edad", report.AgeDistribution)
		}

		if len(report.ExperienceDistribution) > 0 {
			drawDemoDonutChart(pdf, rightX, startY, chartWidth, chartHeight, "Experiencia Laboral", report.ExperienceDistribution)
		}

		// Bottom row: Shift and Marital Status distributions
		bottomRowY := startY + chartHeight + chartSpacingY

		if len(report.ShiftTypeDistribution) > 0 {
			drawDemoDonutChart(pdf, leftX, bottomRowY, chartWidth, chartHeight, "Distribucion por Turno", report.ShiftTypeDistribution)
		}

		if len(report.MaritalStatusDistribution) > 0 {
			drawDemoDonutChart(pdf, rightX, bottomRowY, chartWidth, chartHeight, "Estado Civil", report.MaritalStatusDistribution)
		}
	}

	// ========== DEMOGRAPHIC RISK ANALYSIS ==========
	hasRiskData := len(report.AgeRiskDistribution) > 0 ||
		len(report.ShiftRiskDistribution) > 0 ||
		len(report.ExperienceRiskDistribution) > 0 ||
		len(report.MaritalStatusRiskDistribution) > 0

	if hasRiskData {
		pdf.AddPage()

		// Subtitle
		pdf.SetFont("DejaVu", "B", 16)
		pdf.SetTextColor(59, 130, 246) // Blue
		pdf.SetX(15)
		pdf.Cell(0, 8, "Analisis de Riesgo Demografico")
		pdf.Ln(10)

		// Interpretive analysis paragraph
		pdf.SetFont("DejaVu", "", 10)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(15)
		riskAnalysisText := "El analisis de riesgo demografico proporciona una vista detallada de como se distribuyen los niveles de riesgo psicosocial " +
			"entre diferentes caracteristicas demograficas del personal. Esta informacion permite identificar patrones especificos y " +
			"desarrollar estrategias de prevencion mas efectivas y dirigidas a grupos particulares."
		pdf.MultiCell(180, 5, riskAnalysisText, "", "L", false)
		pdf.Ln(6)

		// Bar chart dimensions - optimized for 2 charts per page while preserving detail
		barChartWidth := 175.0
		barChartHeight := 85.0
		barChartSpacing := 35.0

		// Page 5a: Age and Shift Risk Distributions
		startY := pdf.GetY()

		if len(report.AgeRiskDistribution) > 0 {
			drawDemographicRiskBarChart(pdf, 15, startY, barChartWidth, barChartHeight, "Riesgo por Rango de Edad", report.AgeRiskDistribution)
		}

		if len(report.ShiftRiskDistribution) > 0 {
			secondChartY := startY + barChartHeight + barChartSpacing
			drawDemographicRiskBarChart(pdf, 15, secondChartY, barChartWidth, barChartHeight, "Riesgo por Tipo de Turno", report.ShiftRiskDistribution)
		}

		// Page 5b: Experience and Marital Status Risk Distributions
		pdf.AddPage()


		startY = pdf.GetY()

		if len(report.ExperienceRiskDistribution) > 0 {
			drawDemographicRiskBarChart(pdf, 15, startY, barChartWidth, barChartHeight, "Riesgo por Experiencia Laboral", report.ExperienceRiskDistribution)
		}

		if len(report.MaritalStatusRiskDistribution) > 0 {
			secondChartY := startY + barChartHeight + barChartSpacing
			drawDemographicRiskBarChart(pdf, 15, secondChartY, barChartWidth, barChartHeight, "Riesgo por Estado Civil", report.MaritalStatusRiskDistribution)
		}
	}

	// ========== PAGE 6: RECOMMENDATIONS AND ACTION PLAN ==========
	pdf.AddPage()
	pdf.SetFont("DejaVu", "B", 16)
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

	pdf.SetFont("DejaVu", "", 10)
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
	pdf.SetFont("DejaVu", "B", 16)
	pdf.SetTextColor(59, 130, 246)
	pdf.SetX(15)
	pdf.MultiCell(180, 8, "Declaracion de Politica de Prevencion de Riesgos Psicosociales", "", "L", false)
	pdf.Ln(2)

	pdf.SetFont("DejaVu", "", 10)
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

	// Add hyperlink to NOM-035 reference - use MultiCell to wrap long URL
	nom035URL := "https://www.gob.mx/stps/articulos/norma-oficial-mexicana-nom-035-stps-2018-factores-de-riesgo-psicosocial-en-el-trabajo-identificacion-analisis-y-prevencion"
	pdf.SetTextColor(59, 130, 246) // Blue color for link
	pdf.SetFont("DejaVu", "U", 10) // Underlined for link appearance
	pdf.SetX(15)
	// Use MultiCell to wrap the URL within page width
	pdf.MultiCell(180, 5, nom035URL, "", "L", false)
	pdf.SetFont("DejaVu", "", 10) // Reset font
	pdf.SetTextColor(71, 85, 105) // Reset text color

	// ========== LONG-TERM TRACKING NOTES ==========
	// Force this section to start on a new page
	pdf.AddPage()

	pdf.Ln(10)
	pdf.SetFont("DejaVu", "B", 16)
	pdf.SetTextColor(59, 130, 246)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Notas de Seguimiento a Largo Plazo")
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 10)
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
	pdf.SetFont("DejaVu", "B", 16)
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

	pdf.SetFont("DejaVu", "", 10)
	pdf.SetTextColor(71, 85, 105)
	for _, tip := range engagementTips {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		// Draw a small filled circle as bullet point to avoid encoding issues
		bulletY := pdf.GetY() + 3
		pdf.SetFillColor(71, 85, 105)
		pdf.Circle(20, bulletY, 1, "F")
		pdf.SetX(25)
		pdf.MultiCell(165, 5, tip, "", "L", false)
		pdf.Ln(2)
	}

	// ========== PAGE: OFFICIAL DOCUMENTATION APPENDIX ==========
	pdf.AddPage()

	// Header
	pdf.SetFillColor(15, 23, 42) // Dark blue #1A1F36 (26, 31, 54)
	pdf.Rect(0, 0, 210, 35, "F")

	// Accent stripe
	pdf.SetFillColor(59, 130, 246) // Blue-500
	pdf.Rect(0, 35, 210, 3, "F")

	// Header text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("DejaVu", "B", 18)
	pdf.SetXY(15, 10)
	pdf.Cell(0, 0, "Apéndice de Documentación Oficial")

	pdf.SetFont("DejaVu", "", 11)
	pdf.SetXY(15, 20)
	pdf.Cell(0, 0, "NOM-035-STPS-2018 | Referencias y Glosario")

	// Confidential badge
	pdf.SetFillColor(255, 77, 79) // Red #FF4D4F
	badgeX = 155.0
	badgeY = 8.0
	badgeWidth = 42.0
	badgeHeight = 8.0
	pdf.RoundedRect(badgeX, badgeY, badgeWidth, badgeHeight, 2, "1234", "F")
	pdf.SetFont("DejaVu", "B", 8)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(badgeX, badgeY)
	pdf.CellFormat(badgeWidth, badgeHeight, "CONFIDENCIAL", "", 0, "CM", false, 0, "")

	pdf.SetY(48)

	// ========== GLOSSARY ==========
	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Glosario de Términos")
	pdf.Ln(10)

	glossary := []struct {
		term       string
		definition string
	}{
		{
			term:       "Riesgo Psicosocial",
			definition: "Factores que pueden provocar trastornos de ansiedad, estrés, depresión u otros problemas de salud mental en los trabajadores.",
		},
		{
			term:       "NOM-035-STPS-2018",
			definition: "Norma Oficial Mexicana que establece los elementos para identificar, analizar y prevenir los factores de riesgo psicosocial en los centros de trabajo.",
		},
		{
			term:       "Nivel de Riesgo",
			definition: "Clasificación del riesgo psicosocial según los umbrales establecidos en la NOM-035: Nulo (<20), Bajo (20-<45), Medio (45-<70), Alto (70-<90), Muy Alto (≥90).",
		},
		{
			term:       "Factores de Riesgo Organizacional",
			definition: "Condiciones peligrosas o nocivas inherentes al trabajo, relacionadas con la organización del mismo.",
		},
		{
			term:       "Atención Médica Ocupacional",
			definition: "Servicios médicos especializados para trabajadores expuestos a factores de riesgo psicosocial.",
		},
	}

	pdf.SetFont("DejaVu", "", 10)
	for _, item := range glossary {
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
		pdf.SetFont("DejaVu", "B", 10)
		pdf.SetTextColor(30, 41, 59)
		pdf.SetX(20)
		pdf.Cell(0, 6, item.term+":")
		pdf.Ln(8)
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetX(25)
		pdf.MultiCell(160, 4, item.definition, "", "L", false)
		pdf.Ln(4)
	}

	pdf.Ln(8)

	// ========== OFFICIAL LINKS ==========
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}

	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetX(15)
	pdf.Cell(0, 8, "Enlaces a Documentación Oficial")
	pdf.Ln(10)

	// Single NOM-035 link
	nom035Title := "NOM-035-STPS-2018 - Texto Completo"
	nom035URL = "https://www.gob.mx/stps/articulos/norma-oficial-mexicana-nom-035-stps-2018-factores-de-riesgo-psicosocial-en-el-trabajo-identificacion-analisis-y-prevencion"

	pdf.SetFont("DejaVu", "", 10)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(20)
	pdf.Cell(5, 6, "1.")
	pdf.SetX(25)
	pdf.SetTextColor(59, 130, 246) // Blue for links
	pdf.SetFont("DejaVu", "U", 10) // Underlined
	pdf.MultiCell(165, 5, nom035Title, "", "L", false)
	pdf.SetFont("DejaVu", "", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetX(25)
	pdf.MultiCell(165, 4, nom035URL, "", "L", false)
	pdf.Ln(3)
	pdf.SetFont("DejaVu", "", 10)

	pdf.Ln(8)

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

	pdf.SetFont("DejaVu", "B", 22)
	pdf.SetTextColor(textR, textG, textB)
	valueWidth := pdf.GetStringWidth(value)
	pdf.SetXY(x+(w-valueWidth)/2, y+10)
	pdf.Cell(0, 0, value)

	pdf.SetFont("DejaVu", "", 9)
	pdf.SetTextColor(71, 85, 105)
	labelWidth := pdf.GetStringWidth(label)
	pdf.SetXY(x+(w-labelWidth)/2, y+h-10)
	pdf.Cell(0, 0, label)
}

// drawEnhancedMetricCard draws an enhanced metric card with proportional numbers and better centering
func drawEnhancedMetricCard(pdf *gofpdf.Fpdf, x, y, w, h float64, value, label string, textR, textG, textB, bgR, bgG, bgB int, icon string) {
	// Draw shadow (2px elevation - offset by 2mm)
	shadowOffset := 2.0
	pdf.SetFillColor(0, 0, 0)
	pdf.SetAlpha(0.1, "Normal")
	pdf.RoundedRect(x+shadowOffset, y+shadowOffset, w, h, 4, "1234", "F")
	pdf.SetAlpha(1.0, "Normal")

	// Draw card background
	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.RoundedRect(x, y, w, h, 4, "1234", "F")

	// Proportional bold number (28pt) - better proportioned and centered
	pdf.SetFont("DejaVu", "B", 28)
	pdf.SetTextColor(26, 31, 54) // #1A1F36
	valueWidth := pdf.GetStringWidth(value)
	// Center both horizontally and vertically - position number in upper 2/3 of card
	numberY := y + (h * 0.35) // Position number at 35% from top for better balance
	pdf.SetXY(x+(w-valueWidth)/2, numberY)
	pdf.Cell(0, 0, value)

	// Subtitle (9pt) - slightly smaller for better proportion
	pdf.SetFont("DejaVu", "", 9)
	pdf.SetTextColor(71, 85, 105)
	labelWidth := pdf.GetStringWidth(label)
	// Position label in lower 1/3 of card with better spacing
	labelY := y + (h * 0.75) // Position label at 75% from top
	pdf.SetXY(x+(w-labelWidth)/2, labelY)
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

	pdf.SetFont("DejaVu", "B", 9)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetX(15)
	pdf.Cell(50, 6, title)

	pdf.SetFont("DejaVu", "", 9)
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
	pdf.SetFont("DejaVu", "B", 9)
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
		pdf.SetFont("DejaVu", "", 7)
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

// drawRiskPieChart draws a pie chart for risk distribution
func drawRiskPieChart(pdf *gofpdf.Fpdf, centerX, centerY, radius float64, riskMap map[string]int64, totalCount int64, riskOrder []string) {
	if totalCount == 0 {
		return
	}

	// Risk level colors per NOM-035
	riskColors := map[string]struct{ r, g, b int }{
		"nulo":     {r: 220, g: 252, b: 231}, // Light green
		"bajo":     {r: 40, g: 167, b: 69},   // Green #28A745
		"medio":    {r: 255, g: 193, b: 7},   // Yellow #FFC107
		"alto":     {r: 255, g: 152, b: 0},   // Orange
		"muy_alto": {r: 220, g: 53, b: 69},   // Red #DC3545
	}

	startAngle := 90.0 // Start at top (12 o'clock)
	for _, level := range riskOrder {
		count, exists := riskMap[level]
		if !exists || count == 0 {
			continue
		}

		// Calculate angle for this slice
		percentage := float64(count) / float64(totalCount)
		angle := percentage * 360.0

		// Get color for this risk level
		color, hasColor := riskColors[level]
		if !hasColor {
			color = riskColors["bajo"] // Default to green
		}

		// Draw pie slice
		pdf.SetFillColor(color.r, color.g, color.b)
		pdf.SetDrawColor(255, 255, 255) // White border
		pdf.SetLineWidth(1.0)

		// Draw arc/pie slice
		endAngle := startAngle + angle
		drawPieSlice(pdf, centerX, centerY, radius, startAngle, endAngle)

		// Draw label outside slice if percentage > 5%
		if percentage > 0.05 {
			labelAngle := startAngle + angle/2
			// Position labels outside the pie slice with extra space for bottom labels
			baseOffset := 10.0
			// Add extra space for bottom labels to avoid legend overlap
			normalizedAngle := math.Mod(labelAngle, 360)
			if normalizedAngle > 225 && normalizedAngle < 315 {
				baseOffset = 15.0 // Extra space for bottom labels
			}
			labelRadius := radius + baseOffset
			labelX := centerX + labelRadius*cosDeg(labelAngle)
			labelY := centerY - labelRadius*sinDeg(labelAngle) // Negative because Y increases downward

			pdf.SetFont("DejaVu", "B", 8)
			pdf.SetTextColor(30, 41, 59) // Dark text for contrast on light background
			percentText := fmt.Sprintf("%.1f%%", percentage*100)

			// Get text width for centering
			textWidth := pdf.GetStringWidth(percentText)

			// Center the text on the label position
			pdf.SetXY(labelX-textWidth/2, labelY-2)
			pdf.Cell(0, 0, percentText)
		}

		startAngle = endAngle
	}
}

// drawPieSlice draws a pie slice using arcs
func drawPieSlice(pdf *gofpdf.Fpdf, cx, cy, radius, startAngle, endAngle float64) {
	// Calculate start point
	startX := cx + radius*cosDeg(startAngle)
	startY := cy - radius*sinDeg(startAngle) // Negative because Y increases downward

	// Draw path for pie slice
	pdf.MoveTo(cx, cy)
	pdf.LineTo(startX, startY)
	// Draw arc
	drawArc(pdf, cx, cy, radius, startAngle, endAngle)
	pdf.LineTo(cx, cy)
	pdf.ClosePath()
	pdf.DrawPath("F") // Fill
}

// drawArc draws an arc segment
func drawArc(pdf *gofpdf.Fpdf, cx, cy, radius, startAngle, endAngle float64) {
	// Move to starting point first
	startX := cx + radius*cosDeg(startAngle)
	startY := cy - radius*sinDeg(startAngle)
	pdf.MoveTo(startX, startY)

	// Approximate arc with small line segments
	angleStep := 2.0 // degrees per segment
	currentAngle := startAngle + angleStep
	for currentAngle < endAngle {
		x2 := cx + radius*cosDeg(currentAngle)
		y2 := cy - radius*sinDeg(currentAngle)

		pdf.LineTo(x2, y2)
		currentAngle += angleStep
	}

	// Ensure we reach the exact end point
	if currentAngle != endAngle {
		endX := cx + radius*cosDeg(endAngle)
		endY := cy - radius*sinDeg(endAngle)
		pdf.LineTo(endX, endY)
	}
}

// Helper functions for trigonometry
func cosDeg(angle float64) float64 {
	return math.Cos(angle * math.Pi / 180.0)
}

func sinDeg(angle float64) float64 {
	return math.Sin(angle * math.Pi / 180.0)
}

// drawRiskLegend draws a legend with NOM-035 thresholds (percentages shown on chart)
func drawRiskLegend(pdf *gofpdf.Fpdf, x, y float64, riskMap map[string]int64, totalCount int64, riskOrder []string) {
	riskColors := map[string]struct{ r, g, b int }{
		"nulo":     {r: 220, g: 252, b: 231},
		"bajo":     {r: 40, g: 167, b: 69},
		"medio":    {r: 255, g: 193, b: 7},
		"alto":     {r: 255, g: 152, b: 0},
		"muy_alto": {r: 220, g: 53, b: 69},
	}

	thresholds := map[string]string{
		"nulo":     "<20",
		"bajo":     "20-<45",
		"medio":    "45-<70",
		"alto":     "70-<90",
		"muy_alto": "≥90",
	}

	currentY := y
	itemHeight := 12.0
	boxSize := 8.0

	for _, level := range riskOrder {
		count, exists := riskMap[level]
		if !exists {
			continue
		}

		color := riskColors[level]
		threshold := thresholds[level]

		// Color box
		pdf.SetFillColor(color.r, color.g, color.b)
		pdf.Rect(x, currentY, boxSize, boxSize, "F")

		// Label with threshold (no percentages since they're on the chart)
		pdf.SetFont("DejaVu", "B", 9)
		pdf.SetTextColor(30, 41, 59)
		pdf.SetXY(x+boxSize+4, currentY)
		labelText := fmt.Sprintf("%s (%s): %d", formatRiskLevel(level), threshold, count)
		pdf.Cell(0, boxSize, labelText)

		currentY += itemHeight + 2
	}
}

// drawDemoDonutChart draws a segmented donut chart with proper proportions and legend
func drawDemoDonutChart(pdf *gofpdf.Fpdf, x, y, w, h float64, title string, data []domain.DemographicDistribution) {
	if len(data) == 0 {
		return
	}

	// Card background
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(x, y, w, h, 4, "1234", "F")

	// Title
	pdf.SetFont("DejaVu", "B", 10)
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

	// Pie chart dimensions - scale based on card height with better spacing
	centerX := x + w/2
	centerY := y + h*0.42 // Position center for optimal spacing
	radius := h * 0.20    // Smaller radius to prevent overlap with legend

	// Color palette for demographic segments
	colors := []struct{ r, g, b int }{
		{59, 130, 246}, // Blue
		{16, 185, 129}, // Emerald
		{245, 158, 11}, // Amber
		{239, 68, 68},  // Red
		{139, 92, 246}, // Violet
		{236, 72, 153}, // Pink
		{75, 85, 99},   // Gray
		{6, 182, 212},  // Cyan
	}

	// Draw segmented pie chart
	startAngle := 90.0 // Start at top (12 o'clock)
	for i, d := range data {
		if d.Count == 0 {
			continue
		}

		// Calculate angle for this segment
		percentage := float64(d.Count) / float64(total)
		angle := percentage * 360.0

		// Get color for this segment
		colorIndex := i % len(colors)
		color := colors[colorIndex]

		// Draw pie segment
		pdf.SetFillColor(color.r, color.g, color.b)
		pdf.SetDrawColor(255, 255, 255)
		pdf.SetLineWidth(0.5)

		endAngle := startAngle + angle
		drawPieSlice(pdf, centerX, centerY, radius, startAngle, endAngle)

		// Add percentage label outside the pie slice (only for segments > 3%)
		if percentage > 0.03 {
			labelAngle := startAngle + angle/2
			// Position all labels at consistent distance from chart
			baseOffset := 8.0
			labelRadius := radius + baseOffset
			labelX := centerX + labelRadius*cosDeg(labelAngle)
			labelY := centerY - labelRadius*sinDeg(labelAngle) // Negative because Y increases downward

			// Set font for percentage labels - outside positioning
			pdf.SetFont("DejaVu", "B", 8)
			pdf.SetTextColor(30, 41, 59) // Dark text for contrast on light background

			// Format percentage text
			percentText := fmt.Sprintf("%.1f%%", percentage*100)

			// Get text width for centering
			textWidth := pdf.GetStringWidth(percentText)

			// Center the text on the label position
			pdf.SetXY(labelX-textWidth/2, labelY-2)
			pdf.Cell(0, 0, percentText)
		}

		startAngle = endAngle
	}

	// Draw legend below the chart - position uniformly below all percentage labels
	legendX := x + 4
	itemsPerRow := 2
	legendItemWidth := (w - 8) / float64(itemsPerRow)

	// Calculate legend height
	numLegendItems := len(data)
	if numLegendItems > 6 {
		numLegendItems = 6 // Limit to 6 items
	}
	numRows := (numLegendItems + itemsPerRow - 1) / itemsPerRow // Ceiling division
	legendHeight := float64(numRows) * 5                        // 5 units per row

	// Position legend uniformly below chart - all percentage labels are at same distance
	// Chart center at 42% height, radius 20% height, labels at +8 units = ~54.5 units from top
	// Add 15 units clearance for legend
	legendY := y + h*0.42 + h*0.20 + 8.0 + 15.0 // Chart bottom + label offset + clearance

	// Ensure legend fits within card
	if legendY+legendHeight > y+h {
		legendY = y + h - legendHeight - 2 // Leave 2 units margin at bottom
	}

	pdf.SetFont("DejaVu", "", 7)
	for i, d := range data {
		if i >= 6 { // Limit to 6 items to fit within card space
			break
		}

		row := i / itemsPerRow
		col := i % itemsPerRow
		itemX := legendX + float64(col)*legendItemWidth
		itemY := legendY + float64(row)*5 // Reduced spacing between legend rows

		// Color box
		colorIndex := i % len(colors)
		color := colors[colorIndex]
		pdf.SetFillColor(color.r, color.g, color.b)
		pdf.Rect(itemX, itemY, 4, 4, "F")

		// Label with category name only (percentages now on chart)
		pdf.SetTextColor(71, 85, 105)
		pdf.SetXY(itemX+6, itemY-1)

		categoryName := d.Category
		if len(categoryName) > 16 {
			categoryName = categoryName[:13] + "..."
		}

		pdf.Cell(legendItemWidth-8, 6, categoryName)
	}
}

// drawDonutSlice draws a donut slice (ring segment)
func drawDonutSlice(pdf *gofpdf.Fpdf, cx, cy, outerRadius, innerRadius, startAngle, endAngle float64) {
	// Draw outer arc
	drawArc(pdf, cx, cy, outerRadius, startAngle, endAngle)
	// Draw connecting line at end angle
	endX := cx + outerRadius*cosDeg(endAngle)
	endY := cy - outerRadius*sinDeg(endAngle)
	pdf.LineTo(endX, endY)
	// Draw inner arc (reverse direction)
	drawArcReverse(pdf, cx, cy, innerRadius, endAngle, startAngle)
	// Draw connecting line at start angle
	startX := cx + innerRadius*cosDeg(startAngle)
	startY := cy - innerRadius*sinDeg(startAngle)
	pdf.LineTo(startX, startY)
	// Close path
	pdf.ClosePath()
	pdf.DrawPath("F")
}

// drawArcReverse draws an arc in reverse direction
func drawArcReverse(pdf *gofpdf.Fpdf, cx, cy, radius, startAngle, endAngle float64) {
	// Move to starting point first
	startX := cx + radius*cosDeg(startAngle)
	startY := cy - radius*sinDeg(startAngle)
	pdf.MoveTo(startX, startY)

	// Approximate arc with small line segments in reverse
	angleStep := 2.0
	currentAngle := startAngle - angleStep
	for currentAngle > endAngle {
		x2 := cx + radius*cosDeg(currentAngle)
		y2 := cy - radius*sinDeg(currentAngle)
		pdf.LineTo(x2, y2)
		currentAngle -= angleStep
	}

	// Ensure we reach the exact end point
	if currentAngle != endAngle {
		endX := cx + radius*cosDeg(endAngle)
		endY := cy - radius*sinDeg(endAngle)
		pdf.LineTo(endX, endY)
	}
}

// ChartLayout defines the layout structure for bar charts
type ChartLayout struct {
	CardX, CardY, CardW, CardH     float64 // Card boundaries
	TitleH                         float64 // Title height
	ChartX, ChartY, ChartW, ChartH float64 // Chart area
	LegendY, LegendH               float64 // Legend area
}

// createChartLayout creates a structured layout for the chart
func createChartLayout(x, y, w, h float64, hasAxisLines bool) ChartLayout {
	layout := ChartLayout{
		CardX: x, CardY: y, CardW: w, CardH: h,
		TitleH: 12, // Title height with more space
	}

	if hasAxisLines {
		// Adjust for axis lines and tick labels
		layout.ChartX = x + 25 // Space for Y-axis line and labels
		layout.ChartW = w - 35 // Reduced width for Y-axis space
		layout.ChartH = h - 45 // Reduced space since category labels moved to X-axis
	} else {
		layout.ChartX = x + 4
		layout.ChartW = w - 8
		layout.ChartH = h - 50
	}

	layout.ChartY = y + layout.TitleH + 6               // Space below title
	layout.LegendY = layout.ChartY + layout.ChartH + 25 // Space for X-axis labels and ticks
	layout.LegendH = 15

	return layout
}

// drawDemographicRiskBarChart draws stacked bar charts for demographic risk analysis
func drawDemographicRiskBarChart(pdf *gofpdf.Fpdf, x, y, w, h float64, title string, data []domain.DemographicRiskDistribution) {
	if len(data) == 0 {
		return
	}

	// Determine if this chart needs axis lines (all demographic risk charts)
	hasAxisLines := strings.Contains(title, "Riesgo por")

	// Create structured layout
	layout := createChartLayout(x, y, w, h, hasAxisLines)

	// Card background - extend to cover legend area
	pdf.SetFillColor(248, 250, 252)
	if hasAxisLines {
		// For charts with axis lines, extend background to include legend
		legendBottom := layout.LegendY + layout.LegendH + 5 // Extra padding
		totalHeight := legendBottom - layout.CardY
		pdf.RoundedRect(layout.CardX, layout.CardY, layout.CardW, totalHeight, 4, "1234", "F")
	} else {
		// Standard background for other charts
		pdf.RoundedRect(layout.CardX, layout.CardY, layout.CardW, layout.CardH, 4, "1234", "F")
	}

	// Title
	pdf.SetFont("DejaVu", "B", 10)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(layout.CardX+4, layout.CardY+4)
	pdf.Cell(layout.CardW-8, 5, title)

	// Calculate total for each category and group by category
	categoryTotals := make(map[string]int64)
	categoryRisks := make(map[string]map[string]int64)

	for _, item := range data {
		category := translateCategory(item.Category)
		riskLevel := string(item.RiskLevel)

		if categoryTotals[category] == 0 {
			categoryRisks[category] = make(map[string]int64)
		}
		categoryTotals[category] += item.Count
		categoryRisks[category][riskLevel] = item.Count
	}

	// Get sorted categories
	categories := make([]string, 0, len(categoryTotals))
	for category := range categoryTotals {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	// Use structured layout

	// Calculate bar dimensions
	numCategories := len(categories)
	if numCategories == 0 {
		return
	}

	barWidth := layout.ChartW / float64(numCategories)
	if barWidth > 25 {
		barWidth = 25 // Max bar width
	}
	barSpacing := (layout.ChartW - float64(numCategories)*barWidth) / float64(numCategories+1)

	// Risk level colors (matching the dashboard)
	riskColors := map[string]struct{ r, g, b int }{
		"nulo":     {r: 34, g: 197, b: 94},  // Green
		"bajo":     {r: 132, g: 204, b: 22}, // Light green
		"medio":    {r: 245, g: 158, b: 11}, // Yellow
		"alto":     {r: 249, g: 115, b: 22}, // Orange
		"muy_alto": {r: 239, g: 68, b: 68},  // Red
	}

	riskOrder := []string{"nulo", "bajo", "medio", "alto", "muy_alto"}

	// Draw professional axis lines for age chart
	if hasAxisLines {
		pdf.SetDrawColor(71, 85, 105) // Axis line color
		pdf.SetLineWidth(0.5)

		// Y-axis line (vertical)
		yAxisX := layout.ChartX - 5
		pdf.Line(yAxisX, layout.ChartY, yAxisX, layout.ChartY+layout.ChartH)

		// X-axis line (horizontal)
		xAxisY := layout.ChartY + layout.ChartH + 5
		pdf.Line(layout.ChartX, xAxisY, layout.ChartX+layout.ChartW, xAxisY)

		// Y-axis tick marks and labels
		pdf.SetFont("DejaVu", "", 6)
		pdf.SetTextColor(71, 85, 105)
		maxEmployees := int64(0)
		for _, category := range categories {
			if total := categoryTotals[category]; total > maxEmployees {
				maxEmployees = total
			}
		}

		// Draw 5 tick marks on Y-axis
		for i := 0; i <= 5; i++ {
			yValue := maxEmployees * int64(i) / 5
			yPos := layout.ChartY + layout.ChartH - (float64(yValue)/float64(maxEmployees))*layout.ChartH

			// Tick mark
			pdf.Line(yAxisX-2, yPos, yAxisX, yPos)

			// Label
			label := fmt.Sprintf("%d", yValue)
			labelWidth := pdf.GetStringWidth(label)
			pdf.SetXY(yAxisX-labelWidth-3, yPos-2)
			pdf.Cell(0, 0, label)
		}

		// X-axis category labels (age ranges)
		pdf.SetFont("DejaVu", "", 6)
		for i, category := range categories {
			categoryX := layout.ChartX + barSpacing + float64(i)*(barWidth+barSpacing) + barWidth/2

			// Tick mark
			pdf.Line(categoryX, xAxisY, categoryX, xAxisY+2)

			// Label below tick
			labelWidth := pdf.GetStringWidth(category)
			pdf.SetXY(categoryX-labelWidth/2, xAxisY+4)
			pdf.Cell(0, 0, category)
		}
	}

	// Calculate X-axis position for bars to touch the axis line
	xAxisY := layout.ChartY + layout.ChartH + 5 // X-axis line position

	// Draw bars for each category
	for i, category := range categories {
		categoryX := layout.ChartX + barSpacing + float64(i)*(barWidth+barSpacing)
		currentY := xAxisY // Start from X-axis line

		total := categoryTotals[category]
		if total == 0 {
			continue
		}

		// Draw stacked bars from bottom up (touching X-axis)
		for j := len(riskOrder) - 1; j >= 0; j-- {
			riskLevel := riskOrder[j]
			count := categoryRisks[category][riskLevel]
			if count == 0 {
				continue
			}

			// Calculate bar height proportional to count
			barHeight := (float64(count) / float64(total)) * layout.ChartH // Full chart height available

			// Draw the bar segment
			if color, exists := riskColors[riskLevel]; exists {
				pdf.SetFillColor(color.r, color.g, color.b)
				pdf.Rect(categoryX, currentY-barHeight, barWidth, barHeight, "F")
			}

			currentY -= barHeight
		}

		// Category labels now handled by X-axis tick marks (removed duplicate labels)
	}

	// Draw legend using structured layout
	legendY := layout.LegendY
	legendX := layout.CardX + 4

	// Legend title
	pdf.SetFont("DejaVu", "B", 8)
	pdf.SetTextColor(30, 41, 59) // Dark color for title
	pdf.SetXY(legendX, legendY-8)
	pdf.Cell(0, 0, "Niveles de Riesgo")

	// Separation line
	pdf.SetDrawColor(226, 232, 240) // Light gray line
	pdf.SetLineWidth(0.3)
	lineY := legendY - 4
	pdf.Line(legendX, lineY, legendX+layout.CardW-8, lineY)

	// Legend items
	itemsPerRow := 3
	legendItemWidth := (layout.CardW - 8) / float64(itemsPerRow)

	pdf.SetFont("DejaVu", "", 7)
	for i, riskLevel := range riskOrder {
		if i >= 5 { // Only show first 5 risk levels
			break
		}

		row := i / itemsPerRow
		col := i % itemsPerRow
		itemX := legendX + float64(col)*legendItemWidth
		itemY := legendY + float64(row)*8

		// Color box
		if color, exists := riskColors[riskLevel]; exists {
			pdf.SetFillColor(color.r, color.g, color.b)
			pdf.Rect(itemX, itemY, 4, 4, "F")
		}

		// Label
		pdf.SetTextColor(71, 85, 105)
		pdf.SetXY(itemX+6, itemY-1)

		riskLabel := formatRiskLevel(riskLevel)
		pdf.Cell(legendItemWidth-8, 6, riskLabel)
	}
}

// translateCategory translates demographic category values to readable Spanish labels
func translateCategory(category string) string {
	translations := map[string]string{
		"18-25":       "18-25 años",
		"26-35":       "26-35 años",
		"36-45":       "36-45 años",
		"46-55":       "46-55 años",
		"56+":         "56+ años",
		"diurno":      "Diurno",
		"nocturno":    "Nocturno",
		"mixto":       "Mixto",
		"0-2":         "0-2 años",
		"3-5":         "3-5 años",
		"6-10":        "6-10 años",
		"11-15":       "11-15 años",
		"16-20":       "16-20 años",
		"21+":         "21+ años",
		"soltero":     "Soltero/a",
		"casado":      "Casado/a",
		"divorciado":  "Divorciado/a",
		"viudo":       "Viudo/a",
		"union_libre": "Unión Libre",
		"separado":    "Separado/a",
	}
	if translation, exists := translations[category]; exists {
		return translation
	}
	return category
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
	pdf.SetFont("DejaVu", "B", 10)
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
	pdf.SetFont("DejaVu", "B", 9)
	pdf.SetTextColor(255, 255, 255)
	labelText := fmt.Sprintf("%s: 100%%", maxCategory)
	labelWidth := pdf.GetStringWidth(labelText)
	labelX := barX + (barWidth-labelWidth)/2
	pdf.SetXY(labelX, barY+6)
	pdf.Cell(0, 0, labelText)
}
