package scoring

import (
	"fmt"

	"github.com/entorno35/backend/internal/domain"
)

// Scoring constants for maintainability
const (
	// MaxScorePerQuestion is the maximum score any single question can contribute
	// After polarity adjustment, each question contributes 0-4 points
	MaxScorePerQuestion = 4
)

// ScoreCalculation represents the result of score calculations for categories and domains
type ScoreCalculation struct {
	Scores     map[string]float64 // Actual scores achieved
	MaxScores  map[string]float64 // Maximum possible scores
	RiskLevels map[string]string  // Risk levels for each category/domain
}

// Validate ensures scores don't exceed maximums and provides data integrity
func (sc *ScoreCalculation) Validate() []string {
	var warnings []string

	for name, score := range sc.Scores {
		if maxScore, exists := sc.MaxScores[name]; exists {
			if score > maxScore {
				warnings = append(warnings, fmt.Sprintf("Score %.1f exceeds maximum %.1f for %s", score, maxScore, name))
				// Clamp the score to maximum
				sc.Scores[name] = maxScore
			}
		}
	}

	return warnings
}

// ApplyPolarity applies polarity inversion to a selected value
// Positive polarity: 0→0, 1→1, 2→2, 3→3, 4→4 (no inversion)
// Negative polarity: 0→4, 1→3, 2→2, 3→1, 4→0 (inverted)
func ApplyPolarity(selectedValue int, polarity domain.QuestionPolarity) (int, error) {
	if selectedValue < 0 || selectedValue > 4 {
		return 0, fmt.Errorf("selected value must be between 0 and 4, got %d", selectedValue)
	}

	if polarity == domain.QuestionPolarityPositive {
		// Positive: no inversion (0→0, 1→1, 2→2, 3→3, 4→4)
		return selectedValue, nil
	}

	if polarity == domain.QuestionPolarityNegative {
		// Negative: invert (0→4, 1→3, 2→2, 3→1, 4→0)
		return 4 - selectedValue, nil
	}

	return 0, fmt.Errorf("invalid polarity: %s", polarity)
}

