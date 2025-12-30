package scoring

import "github.com/entorno35/backend/internal/domain"

// AssessmentResult represents the result of a scoring calculation
type AssessmentResult struct {
	TotalScore               float64            `json:"total_score"`
	RiskLevel                domain.RiskLevel   `json:"risk_level"`
	CategoryScores           map[string]float64 `json:"category_scores"`
	DomainScores             map[string]float64 `json:"domain_scores"`
	RequiresMedicalAttention bool               `json:"requires_medical_attention"`
}

// ScoringStrategy defines the interface for scoring strategies (Strategy Pattern)
// Different guides (I, II, III) implement different scoring logic
type ScoringStrategy interface {
	// Calculate computes the assessment result from responses
	// Responses should already be linked to questions with proper relationships loaded
	Calculate(responses []domain.Response) (*AssessmentResult, error)
}

