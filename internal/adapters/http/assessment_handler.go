package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/domain"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssessmentHandler handles assessment-related HTTP requests
type AssessmentHandler struct {
	assessmentService *services.AssessmentService
}

// NewAssessmentHandler creates a new assessment handler
func NewAssessmentHandler(assessmentService *services.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService: assessmentService,
	}
}

// CreateAssessmentRequest represents the request body for creating an assessment
type CreateAssessmentRequest struct {
	StaffID uuid.UUID `json:"staff_id" binding:"required"`
	Period  int       `json:"period" binding:"required"`
}

// CreateAssessmentResponse represents the response for creating an assessment
type CreateAssessmentResponse struct {
	ID        uuid.UUID               `json:"id"`
	StaffID   uuid.UUID               `json:"staff_id"`
	CompanyID uuid.UUID               `json:"company_id"`
	Period    int                     `json:"period"`
	GuideType domain.GuideType        `json:"guide_type"`
	Status    domain.AssessmentStatus `json:"status"`
	CreatedAt string                  `json:"created_at"`
}

// CreateAssessment creates a new assessment
// POST /api/v1/assessments
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return // Already aborted with error
	}

	var req CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify staff belongs to company (should be handled by service/repository)
	assessment, err := h.assessmentService.CreateAssessment(
		services.CreateAssessmentRequest{
			StaffID: req.StaffID,
			Period:  req.Period,
		},
		companyID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAssessmentResponse{
		ID:        assessment.ID,
		StaffID:   assessment.StaffID,
		CompanyID: assessment.CompanyID,
		Period:    assessment.Period,
		GuideType: assessment.GuideType,
		Status:    assessment.Status,
		CreatedAt: assessment.CreatedAt.Format(time.RFC3339),
	})
}

// GetAssessment retrieves an assessment by ID
// GET /api/v1/assessments/:id
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	assessmentIDStr := c.Param("id")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment ID"})
		return
	}

	assessment, err := h.assessmentService.GetAssessment(assessmentID, companyID)
	if err != nil {
		if err.Error() == "failed to fetch assessment: assessment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assessment)
}

// ListAssessmentsRequest represents query parameters for listing assessments
type ListAssessmentsRequest struct {
	StaffID *uuid.UUID               `form:"staff_id"`
	Period  *int                     `form:"period"`
	Status  *domain.AssessmentStatus `form:"status"`
	Limit   int                      `form:"limit,default=50"`
	Offset  int                      `form:"offset,default=0"`
}

// ListAssessmentsResponse represents the paginated response for listing assessments
type ListAssessmentsResponse struct {
	Data   []domain.Assessment `json:"data"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

// ListAssessments retrieves assessments with optional filters
// GET /api/v1/assessments
func (h *AssessmentHandler) ListAssessments(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	var req ListAssessmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessments, total, err := h.assessmentService.ListAssessments(companyID, req.StaffID, req.Period, req.Status, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ListAssessmentsResponse{
		Data:   assessments,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	c.JSON(http.StatusOK, response)
}

// CreateAssessmentLinkRequest represents the request body for creating an assessment link
type CreateAssessmentLinkRequest struct {
	ExpiresInDays int `json:"expires_in_days" binding:"required,min=1,max=365"`
}

// CreateAssessmentLinkResponse represents the response for creating an assessment link
type CreateAssessmentLinkResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Link      string    `json:"link"` // Full URL (to be generated by frontend or service)
}

// CreateAssessmentLink creates a secure link for an assessment
// POST /api/v1/assessments/:id/links
func (h *AssessmentHandler) CreateAssessmentLink(c *gin.Context) {
	assessmentIDStr := c.Param("id")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment ID"})
		return
	}

	// Allow both company admins and staff to create assessment links
	authCtx, ok := middleware.RequireAuth(c)
	if !ok {
		return
	}

	// For staff users, use their own ID; for company users, verify they can create links for their staff
	var requesterID uuid.UUID
	if authCtx.StaffID != nil {
		// Staff user creating link for themselves
		requesterID = *authCtx.StaffID
	} else {
		// Company admin creating link - use their own ID for authorization
		requesterID = authCtx.CompanyID
	}

	var req CreateAssessmentLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use assessment ID from URL
	expiresIn := time.Duration(req.ExpiresInDays) * 24 * time.Hour

	// Prepare request with company ID if it's a company admin
	linkReq := services.CreateAssessmentLinkRequest{
		AssessmentID: assessmentID,
		ExpiresIn:    expiresIn,
	}

	if authCtx.StaffID == nil {
		// Company admin - pass company ID for authorization
		linkReq.CompanyID = &authCtx.CompanyID
	}

	link, err := h.assessmentService.CreateAssessmentLink(linkReq, requesterID)
	if err != nil {
		if err.Error() == "assessment does not belong to staff member" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAssessmentLinkResponse{
		Token:     link.Token,
		ExpiresAt: link.ExpiresAt,
		Link:      "", // Frontend will construct the full URL
	})
}

// GetPublicAssessmentResponse represents the response for getting public assessment details
type GetPublicAssessmentResponse struct {
	Assessment *domain.Assessment `json:"assessment"`
	Questions  []domain.Question  `json:"questions"`
}

// GetPublicAssessment retrieves assessment details for public access via token
// GET /api/v1/assessments/public/:token
// This is a public endpoint (no authentication required)
func (h *AssessmentHandler) GetPublicAssessment(c *gin.Context) {
	// Extract token from URL
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	// Get assessment and questions
	assessment, questions, err := h.assessmentService.GetPublicAssessment(token)
	if err != nil {
		errMsg := err.Error()

		// Check for specific error types (client errors - 400)
		if errMsg == "invalid or expired assessment link: assessment link not found" ||
			errMsg == "assessment link has expired" ||
			errMsg == "assessment link is not associated with an assessment" ||
			errMsg == "assessment is not in pending status" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		// All other errors are server errors - 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, GetPublicAssessmentResponse{
		Assessment: assessment,
		Questions:  questions,
	})
}

// SubmitAssessmentRequest represents the request body for submitting assessment responses
type SubmitAssessmentRequest struct {
	Responses []services.ResponseDTO `json:"responses" binding:"required,min=1"`
}

// SubmitAssessmentResponse represents the response for submitting assessment responses
type SubmitAssessmentResponse struct {
	Message     string    `json:"message"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// SubmitAssessment submits responses for an assessment via public token
// POST /api/v1/assessments/public/:token/submit
// This is a public endpoint (no authentication required)
func (h *AssessmentHandler) SubmitAssessment(c *gin.Context) {
	// Extract token from URL
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	// Parse request body
	var req SubmitAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to service request
	serviceReq := services.SubmitAssessmentRequest{
		Responses: req.Responses,
	}

	// Submit assessment (this will save responses, mark link as used, and trigger scoring)
	err := h.assessmentService.SubmitAssessment(token, serviceReq)
	if err != nil {
		errMsg := err.Error()

		// Check for specific error types (client errors - 400)
		if errMsg == "invalid or expired assessment link: assessment link not found" ||
			errMsg == "assessment link has expired" ||
			errMsg == "assessment link has already been used" ||
			errMsg == "assessment link is not associated with an assessment" ||
			strings.Contains(errMsg, "assessment is not in pending status") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		// All other errors are server errors - 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// Return success response
	// Note: Assessment ID can be retrieved if needed via link lookup, but for public endpoint
	// we keep it simple and just return success
	c.JSON(http.StatusOK, SubmitAssessmentResponse{
		Message:     "Assessment submitted successfully",
		SubmittedAt: time.Now(),
	})
}
