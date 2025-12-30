package scoring

import (
	"fmt"

	"github.com/entorno35/backend/internal/domain"
)

// traumaStrategy implements scoring for Guide I (Trauma Assessment)
// Guide I uses binary Yes/No questions with conditional section logic
type traumaStrategy struct{}

// NewTraumaStrategy creates a new Guide I (Trauma) scoring strategy
func NewTraumaStrategy() ScoringStrategy {
	return &traumaStrategy{}
}

// Calculate computes the assessment result for Guide I (Trauma)
// Logic:
// - Section I: If all "No", assessment ends (no medical attention needed)
// - Section II: Triggered if Section I has at least one "Yes"
// - Section III: Triggered if Section II has at least one "Yes"
// - Section IV: Triggered if Section III has at least one "Yes"
// - Medical attention required if: Section II, III, or IV has at least one "Yes"
func (s *traumaStrategy) Calculate(responses []domain.Response) (*AssessmentResult, error) {
	if len(responses) == 0 {
		return &AssessmentResult{
			TotalScore:             0,
			RiskLevel:              domain.RiskLevelNulo,
			CategoryScores:         make(map[string]float64),
			DomainScores:           make(map[string]float64),
			RequiresMedicalAttention: false,
		}, nil
	}

	// Group responses by section
	sectionI := []domain.Response{}
	sectionII := []domain.Response{}
	sectionIII := []domain.Response{}
	sectionIV := []domain.Response{}

	for _, response := range responses {
		// Load question if not already loaded
		if response.Question.ID == 0 {
			return nil, fmt.Errorf("question not loaded for response %d", response.ID)
		}

		section := response.Question.Section
		if section == nil {
			// Skip responses without section (should not happen for Guide I)
			continue
		}

		switch *section {
		case "I":
			sectionI = append(sectionI, response)
		case "II":
			sectionII = append(sectionII, response)
		case "III":
			sectionIII = append(sectionIII, response)
		case "IV":
			sectionIV = append(sectionIV, response)
		}
	}

	// Check Section I: If all "No", no medical attention needed
	sectionIHasYes := false
	for _, response := range sectionI {
		if response.SelectedValue == 1 { // 1 = Yes for binary
			sectionIHasYes = true
			break
		}
	}

	// If Section I has no "Yes" answers, assessment ends (no medical attention)
	if !sectionIHasYes {
		return &AssessmentResult{
			TotalScore:             0,
			RiskLevel:              domain.RiskLevelNulo,
			CategoryScores:         make(map[string]float64),
			DomainScores:           make(map[string]float64),
			RequiresMedicalAttention: false,
		}, nil
	}

	// Check if any subsequent section has "Yes" answers (triggers medical attention)
	requiresMedicalAttention := false

	// Check Section II
	for _, response := range sectionII {
		if response.SelectedValue == 1 {
			requiresMedicalAttention = true
			break
		}
	}

	// Check Section III (if not already triggered)
	if !requiresMedicalAttention {
		for _, response := range sectionIII {
			if response.SelectedValue == 1 {
				requiresMedicalAttention = true
				break
			}
		}
	}

	// Check Section IV (if not already triggered)
	if !requiresMedicalAttention {
		for _, response := range sectionIV {
			if response.SelectedValue == 1 {
				requiresMedicalAttention = true
				break
			}
		}
	}

	// Guide I doesn't calculate scores (binary Yes/No only)
	return &AssessmentResult{
		TotalScore:             0,
		RiskLevel:              domain.RiskLevelNulo, // Guide I doesn't use risk levels
		CategoryScores:         make(map[string]float64),
		DomainScores:           make(map[string]float64),
		RequiresMedicalAttention: requiresMedicalAttention,
	}, nil
}

