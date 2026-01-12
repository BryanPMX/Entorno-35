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
	ErrLinkNotFound       = errors.New("assessment link not found")
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

// Create creates a new assessment
func (r *assessmentRepository) Create(assessment *domain.Assessment) error {
	result := r.db.Create(assessment)
	if result.Error != nil {
		return fmt.Errorf("failed to create assessment: %w", result.Error)
	}

	return nil
}

// ListByCompany retrieves all assessments for a company with optional filters
func (r *assessmentRepository) ListByCompany(companyID uuid.UUID, staffID *uuid.UUID, period *int, status *domain.AssessmentStatus) ([]domain.Assessment, error) {
	var assessments []domain.Assessment

	query := r.db.Where("company_id = ?", companyID)

	if staffID != nil {
		query = query.Where("staff_id = ?", *staffID)
	}

	if period != nil {
		query = query.Where("period = ?", *period)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	result := query.
		Preload("Staff").
		Order("created_at DESC").
		Find(&assessments)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list assessments: %w", result.Error)
	}

	return assessments, nil
}

// ListByCompanyPaginated retrieves assessments for a company with pagination and filters
func (r *assessmentRepository) ListByCompanyPaginated(companyID uuid.UUID, staffID *uuid.UUID, period *int, status *domain.AssessmentStatus, limit int, offset int) ([]domain.Assessment, int, error) {
	var assessments []domain.Assessment
	var total int64

	query := r.db.Model(&domain.Assessment{}).Where("company_id = ?", companyID)

	if staffID != nil {
		query = query.Where("staff_id = ?", *staffID)
	}

	if period != nil {
		query = query.Where("period = ?", *period)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count assessments: %w", err)
	}

	// Get paginated results
	result := query.
		Preload("Staff").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&assessments)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list assessments: %w", result.Error)
	}

	return assessments, int(total), nil
}

// CreateLink creates a new assessment link with secure token
func (r *assessmentRepository) CreateLink(link *domain.AssessmentLink) error {
	result := r.db.Create(link)
	if result.Error != nil {
		return fmt.Errorf("failed to create assessment link: %w", result.Error)
	}

	return nil
}

// GetLinkByToken retrieves an assessment link by token
func (r *assessmentRepository) GetLinkByToken(token string) (*domain.AssessmentLink, error) {
	var link domain.AssessmentLink

	result := r.db.
		Preload("Staff").
		Preload("Assessment").
		Where("token = ?", token).
		First(&link)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, fmt.Errorf("failed to fetch assessment link: %w", result.Error)
	}

	return &link, nil
}

// UpdateLinkAccess updates the link's accessed_at timestamp
func (r *assessmentRepository) UpdateLinkAccess(linkID uuid.UUID) error {
	return r.UpdateLinkAccessedAt(linkID)
}

// UpdateLinkAccessedAt updates the link's accessed_at timestamp
func (r *assessmentRepository) UpdateLinkAccessedAt(linkID uuid.UUID) error {
	now := time.Now()

	result := r.db.Model(&domain.AssessmentLink{}).
		Where("id = ?", linkID).
		Update("accessed_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to update link accessed_at: %w", result.Error)
	}

	return nil
}

// GetQuestionsByGuideType retrieves all questions for a specific guide type
func (r *assessmentRepository) GetQuestionsByGuideType(guideType domain.GuideType) ([]domain.Question, error) {
	var questions []domain.Question

	result := r.db.
		Where("guide_type = ?", guideType).
		Order("order_index ASC").
		Find(&questions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch questions for guide type %s: %w", guideType, result.Error)
	}

	return questions, nil
}

// GetQuestionsByGuideTypeWithRelations retrieves all questions for a specific guide type with relationships preloaded
func (r *assessmentRepository) GetQuestionsByGuideTypeWithRelations(guideType domain.GuideType) ([]domain.Question, error) {
	var questions []domain.Question

	result := r.db.
		Preload("Category").
		Preload("Domain").
		Preload("Dimension").
		Where("guide_type = ?", guideType).
		Order("order_index ASC").
		Find(&questions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch questions with relations for guide type %s: %w", guideType, result.Error)
	}

	return questions, nil
}

// Delete removes an assessment and its associated data (links and responses)
// Uses a transaction to ensure data integrity
func (r *assessmentRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, delete associated assessment links
		if err := tx.Where("assessment_id = ?", id).Delete(&domain.AssessmentLink{}).Error; err != nil {
			return fmt.Errorf("failed to delete assessment links: %w", err)
		}
		// Then, delete associated responses
		if err := tx.Where("assessment_id = ?", id).Delete(&domain.Response{}).Error; err != nil {
			return fmt.Errorf("failed to delete responses: %w", err)
		}
		// Finally, delete the assessment itself
		if err := tx.Delete(&domain.Assessment{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete assessment: %w", err)
		}
		return nil
	})
}
