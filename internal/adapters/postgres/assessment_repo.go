package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/core/scoring"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAssessmentNotFound = errors.New("assessment not found")
)

// assessmentRepository implements the AssessmentRepository interface using GORM (High cohesion - persistence concerns only)
type assessmentRepository struct {
	db *gorm.DB
}

// NewAssessmentRepository creates a new Postgres implementation of AssessmentRepository
func NewAssessmentRepository(db *gorm.DB) ports.AssessmentRepository {
	return &assessmentRepository{db: db}
}

// GetByID retrieves an assessment by ID with all relationships loaded
func (r *assessmentRepository) GetByID(id uuid.UUID) (*domain.Assessment, error) {
	var assessment domain.Assessment

	result := r.db.
		Preload("Staff").
		Preload("Company").
		Where("id = ?", id).
		First(&assessment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAssessmentNotFound
		}
		return nil, fmt.Errorf("failed to fetch assessment: %w", result.Error)
	}

	return &assessment, nil
}

// GetByIDAndCompany retrieves an assessment by ID and company ID (for authorization)
func (r *assessmentRepository) GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Assessment, error) {
	var assessment domain.Assessment

	result := r.db.
		Preload("Staff").
		Preload("Company").
		Where("id = ? AND company_id = ?", id, companyID).
		First(&assessment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAssessmentNotFound
		}
		return nil, fmt.Errorf("failed to fetch assessment: %w", result.Error)
	}

	return &assessment, nil
}

// GetResponsesByAssessmentID retrieves all responses for an assessment with questions loaded
// Questions must have Category, Domain, and Dimension relationships loaded for scoring
func (r *assessmentRepository) GetResponsesByAssessmentID(assessmentID uuid.UUID) ([]domain.Response, error) {
	var responses []domain.Response

	result := r.db.
		Preload("Question").
		Preload("Question.Category").
		Preload("Question.Domain").
		Preload("Question.Dimension").
		Where("assessment_id = ?", assessmentID).
		Find(&responses)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch responses: %w", result.Error)
	}

	return responses, nil
}

// UpdateResult updates the assessment with scoring results
func (r *assessmentRepository) UpdateResult(assessment *domain.Assessment, result *scoring.AssessmentResult) error {
	now := time.Now()

	// Update assessment fields
	assessment.TotalScore = &result.TotalScore
	assessment.RiskLevel = &result.RiskLevel
	assessment.RequiresMedicalAttention = result.RequiresMedicalAttention
	
	// Set completed_at if not already set
	if assessment.CompletedAt == nil {
		assessment.CompletedAt = &now
	}

	// Update status to completed
	assessment.Status = domain.AssessmentStatusCompleted

	// Update in database
	updateResult := r.db.Model(assessment).Updates(map[string]interface{}{
		"total_score":               assessment.TotalScore,
		"risk_level":                assessment.RiskLevel,
		"requires_medical_attention": assessment.RequiresMedicalAttention,
		"status":                    assessment.Status,
		"completed_at":              assessment.CompletedAt,
		"updated_at":                now,
	})

	if updateResult.Error != nil {
		return fmt.Errorf("failed to update assessment result: %w", updateResult.Error)
	}

	return nil
}

