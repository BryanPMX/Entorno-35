package http

import (
	"net/http"
	"strconv"

	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	if periodStr := c.Query("period"); periodStr != "" {
		periodInt, err := strconv.Atoi(periodStr)
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

