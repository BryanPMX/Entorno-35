package http

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/domain"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	reportService *services.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetIndividualReport retrieves an individual assessment report
// GET /api/v1/reports/individual/:assessment_id
func (h *ReportHandler) GetIndividualReport(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return // Already aborted with error
	}

	// Extract assessment ID from URL
	assessmentIDStr := c.Param("assessment_id")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment ID"})
		return
	}

	// Generate individual report
	report, err := h.reportService.GenerateIndividualReport(assessmentID, companyID)
	if err != nil {
		// Check if report not found
		if err.Error() == "failed to get individual report: report not found" ||
			err.Error() == "failed to get individual report: assessment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "assessment report not found"})
			return
		}
		if err.Error() == "failed to get individual report: assessment has not been scored yet" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "assessment has not been scored yet"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetIndividualReportPDF generates and returns a PDF report for an individual assessment
// GET /api/v1/reports/individual/:assessment_id/pdf
func (h *ReportHandler) GetIndividualReportPDF(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return // Already aborted with error
	}

	// Extract assessment ID from URL
	assessmentIDStr := c.Param("assessment_id")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment ID"})
		return
	}

	// Generate individual report
	report, err := h.reportService.GenerateIndividualReport(assessmentID, companyID)
	if err != nil {
		// Check if report not found
		if err.Error() == "failed to get individual report: report not found" ||
			err.Error() == "failed to get individual report: assessment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "assessment report not found"})
			return
		}
		if err.Error() == "failed to get individual report: assessment has not been scored yet" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "assessment has not been scored yet"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF
	pdf, err := h.generateIndividualReportPDF(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PDF"})
		return
	}

	// Set headers for PDF download
	filename := fmt.Sprintf("NOM035_Report_%s_%s.pdf", report.StaffName, time.Now().Format("20060102"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdf)))

	// Send PDF data
	c.Data(http.StatusOK, "application/pdf", pdf)
}

// generateIndividualReportPDF creates a professional PDF report
func (h *ReportHandler) generateIndividualReportPDF(report *domain.IndividualReportDTO) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Helvetica", "", 12)

	lineHeight := 7.0

	// Header
	pdf.SetFont("Helvetica", "B", 20)
	pdf.Cell(0, lineHeight*2, "NOM-035 Psychosocial Risk Assessment Report")
	pdf.Ln(lineHeight * 2.5)

	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(0, lineHeight, "Guía II - Identificación y análisis de los factores de riesgo psicosocial")
	pdf.Ln(lineHeight * 1.5)

	// Assessment Info Section
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, lineHeight, "Assessment Information")
	pdf.Ln(lineHeight * 1.2)

	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, lineHeight, fmt.Sprintf("Staff Member: %s", report.StaffName))
	pdf.Ln(lineHeight)

	if report.Department != "" {
		pdf.Cell(0, lineHeight, fmt.Sprintf("Department: %s", report.Department))
		pdf.Ln(lineHeight)
	}

	if report.Shift != "" {
		pdf.Cell(0, lineHeight, fmt.Sprintf("Shift: %s", report.Shift))
		pdf.Ln(lineHeight)
	}

	pdf.Cell(0, lineHeight, fmt.Sprintf("Period: %d", report.Period))
	pdf.Ln(lineHeight)

	if report.CompletedAt != nil {
		pdf.Cell(0, lineHeight, fmt.Sprintf("Completed: %s", report.CompletedAt.Format("2006-01-02 15:04:05")))
		pdf.Ln(lineHeight)
	}

	pdf.Ln(lineHeight)

	// Risk Assessment Section
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, lineHeight, "Risk Assessment Results")
	pdf.Ln(lineHeight * 1.2)

	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(0, lineHeight, fmt.Sprintf("Total Score: %.0f/100", report.TotalScore))
	pdf.Ln(lineHeight)

	pdf.Cell(0, lineHeight, fmt.Sprintf("Risk Level: %s", report.RiskLevel))
	pdf.Ln(lineHeight)

	if report.RequiresMedical {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(255, 0, 0) // Red color
		pdf.Cell(0, lineHeight, "⚠️  MEDICAL ATTENTION REQUIRED")
		pdf.SetTextColor(0, 0, 0) // Reset to black
		pdf.Ln(lineHeight)
	}

	pdf.Ln(lineHeight)

	// Category Scores Section
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, lineHeight, "Category Scores")
	pdf.Ln(lineHeight * 1.2)

	pdf.SetFont("Helvetica", "", 10)
	for category, score := range report.CategoryScores {
		pdf.Cell(0, lineHeight, fmt.Sprintf("%s: %.1f/20", category, score))
		pdf.Ln(lineHeight)
	}

	pdf.Ln(lineHeight)

	// Domain Scores Section
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, lineHeight, "Domain Scores")
	pdf.Ln(lineHeight * 1.2)

	pdf.SetFont("Helvetica", "", 10)
	for domain, score := range report.DomainScores {
		pdf.Cell(0, lineHeight, fmt.Sprintf("%s: %.1f/15", domain, score))
		pdf.Ln(lineHeight)
	}

	pdf.Ln(lineHeight)

	// Recommendations Section
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, lineHeight, "Recommendations")
	pdf.Ln(lineHeight * 1.2)

	pdf.SetFont("Helvetica", "", 10)
	for i, rec := range report.Recommendations {
		pdf.MultiCell(0, lineHeight, fmt.Sprintf("%d. %s", i+1, rec), "", "", false)
		pdf.Ln(lineHeight * 0.5)
	}

	// Footer
	pdf.SetY(-30)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.Cell(0, lineHeight, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(lineHeight)
	pdf.Cell(0, lineHeight, "NOM-035-STPS-2018 Compliance Report - Entorno35 Platform")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetGeneralReportRequest represents query parameters for general report
type GetGeneralReportRequest struct {
	Period *int `form:"period"` // Optional: filter by period (e.g., 2025)
}

// GetGeneralReport retrieves a company-wide general report
// GET /api/v1/reports/general?period=2025
func (h *ReportHandler) GetGeneralReport(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return // Already aborted with error
	}

	// Parse query parameters
	var req GetGeneralReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	// Parse period from query string if provided
	var period *int
	var periodInt int
	if periodStr := c.Query("period"); periodStr != "" {
		var err error
		periodInt, err = strconv.Atoi(periodStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period format"})
			return
		}
		period = &periodInt
	}

	// Generate general report
	report, err := h.reportService.GenerateGeneralReport(companyID, period)
	if err != nil {
		if err.Error() == "failed to get general report: company not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
