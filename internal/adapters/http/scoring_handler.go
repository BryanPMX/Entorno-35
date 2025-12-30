package http

import (
	"net/http"

	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ScoringHandler handles scoring-related HTTP requests (High cohesion - HTTP concerns only)
type ScoringHandler struct {
	scoringService *services.ScoringService
}

// NewScoringHandler creates a new scoring handler
func NewScoringHandler(scoringService *services.ScoringService) *ScoringHandler {
	return &ScoringHandler{
		scoringService: scoringService,
	}
}

// CalculateAssessment calculates the score for an assessment
// POST /assessments/:id/calculate
func (h *ScoringHandler) CalculateAssessment(c *gin.Context) {
	// Extract assessment ID from URL
	assessmentIDStr := c.Param("id")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment ID"})
		return
	}

	// Verify company access (must belong to authenticated company)
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return // Already aborted with error
	}

	// Calculate assessment score (with company authorization)
	err = h.scoringService.CalculateAssessmentWithCompany(assessmentID, companyID)
	if err != nil {
		// Check if assessment not found
		if err.Error() == "failed to fetch assessment: assessment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Assessment scored successfully",
		"assessment_id": assessmentID,
	})
}

