package ports

import (
	"github.com/entorno35/backend/internal/core/scoring"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// AssessmentRepository defines the interface for assessment repository operations
type AssessmentRepository interface {
	// GetByID retrieves an assessment by ID with all relationships loaded
	GetByID(id uuid.UUID) (*domain.Assessment, error)
	
	// GetByIDAndCompany retrieves an assessment by ID and company ID (for authorization)
	GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Assessment, error)
	
	// GetResponsesByAssessmentID retrieves all responses for an assessment with questions loaded
	GetResponsesByAssessmentID(assessmentID uuid.UUID) ([]domain.Response, error)
	
	// UpdateResult updates the assessment with scoring results
	UpdateResult(assessment *domain.Assessment, result *scoring.AssessmentResult) error

	// Create creates a new assessment
	Create(assessment *domain.Assessment) error

	// ListByCompany retrieves all assessments for a company with optional filters
	ListByCompany(companyID uuid.UUID, staffID *uuid.UUID, period *int, status *domain.AssessmentStatus) ([]domain.Assessment, error)

	// CreateLink creates a new assessment link with secure token
	CreateLink(link *domain.AssessmentLink) error

	// GetLinkByToken retrieves an assessment link by token
	GetLinkByToken(token string) (*domain.AssessmentLink, error)

	// UpdateLinkAccess updates the link's accessed_at timestamp
	UpdateLinkAccess(linkID uuid.UUID) error

	// UpdateLinkAccessedAt updates the link's accessed_at timestamp (same as UpdateLinkAccess, alias for clarity)
	UpdateLinkAccessedAt(linkID uuid.UUID) error

	// GetQuestionsByGuideType retrieves all questions for a specific guide type
	GetQuestionsByGuideType(guideType domain.GuideType) ([]domain.Question, error)
}
