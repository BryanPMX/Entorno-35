package services

import (
	"fmt"
	"time"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// AssessmentService handles assessment business logic
type AssessmentService struct {
	assessmentRepo ports.AssessmentRepository
	companyRepo    ports.CompanyRepository
	staffRepo      ports.StaffRepository
	responseRepo   ports.ResponseRepository
	scoringService *ScoringService
}

// NewAssessmentService creates a new assessment service
func NewAssessmentService(
	assessmentRepo ports.AssessmentRepository,
	companyRepo ports.CompanyRepository,
	staffRepo ports.StaffRepository,
	responseRepo ports.ResponseRepository,
	scoringService *ScoringService,
) *AssessmentService {
	return &AssessmentService{
		assessmentRepo: assessmentRepo,
		companyRepo:    companyRepo,
		staffRepo:      staffRepo,
		responseRepo:   responseRepo,
		scoringService: scoringService,
	}
}

// CreateAssessmentRequest represents the request to create an assessment
type CreateAssessmentRequest struct {
	StaffID uuid.UUID
	Period  int // e.g., 2025
}

// CreateAssessment creates a new assessment for a staff member
// Determines guide type based on company employee count:
// - Guide I: Always for trauma assessment (optional)
// - Guide II: 16-50 employees
// - Guide III: >50 employees
func (s *AssessmentService) CreateAssessment(req CreateAssessmentRequest, companyID uuid.UUID) (*domain.Assessment, error) {
	// Verify staff belongs to company
	_, err := s.staffRepo.GetByIDAndCompany(req.StaffID, companyID)
	if err != nil {
		return nil, fmt.Errorf("staff not found or does not belong to company: %w", err)
	}

	// Get company to determine guide type
	company, err := s.companyRepo.GetByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch company: %w", err)
	}

	// Determine guide type based on employee count
	var guideType domain.GuideType
	if company.EmployeeCount >= 16 && company.EmployeeCount <= 50 {
		guideType = domain.GuideTypeII
	} else if company.EmployeeCount > 50 {
		guideType = domain.GuideTypeIII
	} else {
		// Companies with <16 employees typically use Guide II
		guideType = domain.GuideTypeII
	}

	// Create assessment
	assessment := &domain.Assessment{
		StaffID:                  req.StaffID,
		CompanyID:                companyID,
		Period:                   req.Period,
		GuideType:                guideType,
		Status:                   domain.AssessmentStatusPending,
		RequiresMedicalAttention: false,
	}

	err = s.assessmentRepo.Create(assessment)
	if err != nil {
		return nil, fmt.Errorf("failed to create assessment: %w", err)
	}

	return assessment, nil
}

// CreateAssessmentLinkRequest represents the request to create an assessment link
type CreateAssessmentLinkRequest struct {
	AssessmentID uuid.UUID
	ExpiresIn    time.Duration // e.g., 7 * 24 * time.Hour for 7 days
}

// CreateAssessmentLink creates a secure link for staff to access an assessment
func (s *AssessmentService) CreateAssessmentLink(req CreateAssessmentLinkRequest, staffID uuid.UUID) (*domain.AssessmentLink, error) {
	// Verify assessment exists and belongs to the staff member
	assessment, err := s.assessmentRepo.GetByID(req.AssessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assessment: %w", err)
	}

	if assessment.StaffID != staffID {
		return nil, fmt.Errorf("assessment does not belong to staff member")
	}

	// Generate secure token (UUID)
	token := uuid.New().String()

	// Create link
	link := &domain.AssessmentLink{
		Token:        token,
		StaffID:      staffID,
		AssessmentID: &req.AssessmentID,
		ExpiresAt:    time.Now().Add(req.ExpiresIn),
	}

	err = s.assessmentRepo.CreateLink(link)
	if err != nil {
		return nil, fmt.Errorf("failed to create assessment link: %w", err)
	}

	return link, nil
}

// GetAssessment retrieves an assessment by ID with company authorization
func (s *AssessmentService) GetAssessment(assessmentID uuid.UUID, companyID uuid.UUID) (*domain.Assessment, error) {
	assessment, err := s.assessmentRepo.GetByIDAndCompany(assessmentID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assessment: %w", err)
	}

	return assessment, nil
}

// ListAssessments retrieves assessments for a company with optional filters
func (s *AssessmentService) ListAssessments(companyID uuid.UUID, staffID *uuid.UUID, period *int, status *domain.AssessmentStatus) ([]domain.Assessment, error) {
	assessments, err := s.assessmentRepo.ListByCompany(companyID, staffID, period, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list assessments: %w", err)
	}

	return assessments, nil
}

// ResponseDTO represents a response in the submission request
type ResponseDTO struct {
	QuestionID uint `json:"question_id" binding:"required"`
	Value      int  `json:"value" binding:"required,min=0,max=4"`
}

// SubmitAssessmentRequest represents the request to submit assessment responses
type SubmitAssessmentRequest struct {
	Responses []ResponseDTO `json:"responses" binding:"required,min=1"`
}

// SubmitAssessment submits responses for an assessment via public token
// This is the public endpoint that staff use to submit their answers
func (s *AssessmentService) SubmitAssessment(token string, req SubmitAssessmentRequest) error {
	// 1. Validate Link: Fetch AssessmentLink by token
	link, err := s.assessmentRepo.GetLinkByToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired assessment link: %w", err)
	}

	// Check if link is expired
	if time.Now().After(link.ExpiresAt) {
		return fmt.Errorf("assessment link has expired")
	}

	// Check if link has already been used (accessed_at is set)
	if link.AccessedAt != nil {
		return fmt.Errorf("assessment link has already been used")
	}

	// Ensure assessment ID is set
	if link.AssessmentID == nil {
		return fmt.Errorf("assessment link is not associated with an assessment")
	}
	assessmentID := *link.AssessmentID

	// Fetch assessment to verify it exists and is in pending status
	assessment, err := s.assessmentRepo.GetByID(assessmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch assessment: %w", err)
	}

	if assessment.Status != domain.AssessmentStatusPending {
		return fmt.Errorf("assessment is not in pending status (current status: %s)", assessment.Status)
	}

	// 2. Map DTOs to Domain models
	responses := make([]domain.Response, 0, len(req.Responses))
	now := time.Now()

	for _, dto := range req.Responses {
		response := domain.Response{
			AssessmentID:  assessmentID,
			QuestionID:    dto.QuestionID,
			SelectedValue: dto.Value,
			AnsweredAt:    &now,
		}

		// Note: CalculatedScore will be set during scoring calculation
		// For now, we'll set it to SelectedValue (scoring service will recalculate with polarity)
		response.CalculatedScore = dto.Value

		responses = append(responses, response)
	}

	// 3. Save Answers in transaction
	err = s.responseRepo.SaveResponses(responses)
	if err != nil {
		return fmt.Errorf("failed to save responses: %w", err)
	}

	// 4. Mark Link as Used (update accessed_at)
	err = s.assessmentRepo.UpdateLinkAccessedAt(link.ID)
	if err != nil {
		return fmt.Errorf("failed to mark assessment link as used: %w", err)
	}

	// 5. Trigger Scoring
	err = s.scoringService.CalculateAssessmentWithCompany(assessmentID, assessment.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to calculate assessment score: %w", err)
	}

	// Note: Scoring service already updates assessment status to COMPLETED
	// and sets completed_at timestamp, so no need to do it here

	return nil
}
