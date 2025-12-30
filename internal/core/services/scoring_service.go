package services

import (
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/core/scoring"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// ScoringService handles assessment scoring operations
type ScoringService struct {
	assessmentRepo ports.AssessmentRepository
	rules          *scoring.ScoringRules
}

// NewScoringService creates a new scoring service
func NewScoringService(assessmentRepo ports.AssessmentRepository) *ScoringService {
	return &ScoringService{
		assessmentRepo: assessmentRepo,
		rules:          scoring.LoadScoringRules(),
	}
}

// CalculateAssessment calculates the score for an assessment
// Factory logic: Selects appropriate strategy based on GuideType
func (s *ScoringService) CalculateAssessment(assessmentID uuid.UUID) error {
	// Fetch assessment with relationships
	assessment, err := s.assessmentRepo.GetByID(assessmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch assessment: %w", err)
	}

	// Fetch all responses with questions loaded
	responses, err := s.assessmentRepo.GetResponsesByAssessmentID(assessmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch responses: %w", err)
	}

	// Select strategy based on guide type (Factory Pattern)
	var strategy scoring.ScoringStrategy
	switch assessment.GuideType {
	case domain.GuideTypeI:
		strategy = scoring.NewTraumaStrategy()
	case domain.GuideTypeII:
		strategy = scoring.NewRiskStrategy(domain.GuideTypeII, s.rules)
	case domain.GuideTypeIII:
		strategy = scoring.NewRiskStrategy(domain.GuideTypeIII, s.rules)
	default:
		return fmt.Errorf("invalid guide type: %s", assessment.GuideType)
	}

	// Calculate result using selected strategy
	result, err := strategy.Calculate(responses)
	if err != nil {
		return fmt.Errorf("failed to calculate assessment score: %w", err)
	}

	// Update assessment with results
	err = s.assessmentRepo.UpdateResult(assessment, result)
	if err != nil {
		return fmt.Errorf("failed to update assessment result: %w", err)
	}

	return nil
}

