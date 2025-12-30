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
}

// NewAssessmentService creates a new assessment service
func NewAssessmentService(assessmentRepo ports.AssessmentRepository, companyRepo ports.CompanyRepository, staffRepo ports.StaffRepository) *AssessmentService {
	return &AssessmentService{
		assessmentRepo: assessmentRepo,
		companyRepo:    companyRepo,
		staffRepo:      staffRepo,
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

