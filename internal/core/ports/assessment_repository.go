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
}

