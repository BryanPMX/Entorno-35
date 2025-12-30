package scoring

import (
	"fmt"

	"github.com/entorno35/backend/internal/domain"
)

// riskStrategy implements scoring for Guide II and Guide III (Psychosocial Risk Assessment)
// Uses Likert scale (0-4) with polarity inversion and hierarchical aggregation
type riskStrategy struct {
	guideType  domain.GuideType
	rules      *ScoringRules
}

// NewRiskStrategy creates a new Guide II/III (Risk) scoring strategy
func NewRiskStrategy(guideType domain.GuideType, rules *ScoringRules) ScoringStrategy {
	return &riskStrategy{
		guideType: guideType,
		rules:     rules,
	}
}

// Calculate computes the assessment result for Guide II or Guide III
// Aggregation: Question → Dimension → Domain → Category → Total
func (s *riskStrategy) Calculate(responses []domain.Response) (*AssessmentResult, error) {
	if len(responses) == 0 {
		return &AssessmentResult{
			TotalScore:             0,
			RiskLevel:              domain.RiskLevelNulo,
			CategoryScores:         make(map[string]float64),
			DomainScores:           make(map[string]float64),
			RequiresMedicalAttention: false,
		}, nil
	}

	// Map to store scores by dimension (for aggregation)
	dimensionScores := make(map[string]float64)
	dimensionCounts := make(map[string]int)

	// Map to store scores by domain (for aggregation)
	domainScores := make(map[string]float64)
	domainCounts := make(map[string]int)

	// Map to store scores by category (for aggregation)
	categoryScores := make(map[string]float64)
	categoryCounts := make(map[string]int)

	totalScore := float64(0)

	// Process each response
	for _, response := range responses {
		if response.Question.ID == 0 {
			return nil, fmt.Errorf("question not loaded for response %d", response.ID)
		}

		question := response.Question

		// Apply polarity inversion to get calculated score
		if question.Polarity == nil {
			return nil, fmt.Errorf("question %d (Guide II/III) must have polarity", question.ID)
		}

		calculatedScore, err := ApplyPolarity(response.SelectedValue, *question.Polarity)
		if err != nil {
			return nil, fmt.Errorf("failed to apply polarity to question %d: %w", question.ID, err)
		}

		// Accumulate total score
		totalScore += float64(calculatedScore)

		// Aggregate by dimension (if dimension exists)
		if question.DimensionID != nil && question.Dimension != nil {
			dimensionName := question.Dimension.Name
			dimensionScores[dimensionName] += float64(calculatedScore)
			dimensionCounts[dimensionName]++
		}

		// Aggregate by domain
		if question.DomainID != nil {
			domainName := question.Domain.Name
			domainScores[domainName] += float64(calculatedScore)
			domainCounts[domainName]++
		}

		// Aggregate by category
		if question.CategoryID != nil {
			categoryName := question.Category.Name
			categoryScores[categoryName] += float64(calculatedScore)
			categoryCounts[categoryName]++
		}
	}

	// Calculate risk level based on total score
	var riskLevel domain.RiskLevel
	var thresholds RiskThresholds

	if s.guideType == domain.GuideTypeII {
		thresholds = s.rules.GuideII.TotalScore
	} else if s.guideType == domain.GuideTypeIII {
		thresholds = s.rules.GuideIII.TotalScore
	} else {
		return nil, fmt.Errorf("invalid guide type for risk strategy: %s", s.guideType)
	}

	riskLevel = thresholds.GetRiskLevel(totalScore)

	return &AssessmentResult{
		TotalScore:             totalScore,
		RiskLevel:              riskLevel,
		CategoryScores:         categoryScores,
		DomainScores:           domainScores,
		RequiresMedicalAttention: false, // Guide II/III don't use medical attention flag
	}, nil
}

