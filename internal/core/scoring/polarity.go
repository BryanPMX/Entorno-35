package scoring

import (
	"fmt"

	"github.com/entorno35/backend/internal/domain"
)

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

