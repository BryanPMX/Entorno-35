package scoring

import (
	"testing"

	"github.com/entorno35/backend/internal/domain"
)

func TestApplyPolarity(t *testing.T) {
	tests := []struct {
		name           string
		selectedValue  int
		polarity       domain.QuestionPolarity
		expectedScore  int
		expectedError  bool
	}{
		// Positive polarity (0→4 scale: Siempre=0, Nunca=4)
		{
			name:          "positive polarity - siempre (0)",
			selectedValue: 0,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 0,
			expectedError: false,
		},
		{
			name:          "positive polarity - casi siempre (1)",
			selectedValue: 1,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 1,
			expectedError: false,
		},
		{
			name:          "positive polarity - algunas veces (2)",
			selectedValue: 2,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 2,
			expectedError: false,
		},
		{
			name:          "positive polarity - casi nunca (3)",
			selectedValue: 3,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 3,
			expectedError: false,
		},
		{
			name:          "positive polarity - nunca (4)",
			selectedValue: 4,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 4,
			expectedError: false,
		},
		// Negative polarity (4→0 scale: Siempre=4, Nunca=0) - INVERTED
		{
			name:          "negative polarity - siempre (0→4)",
			selectedValue: 0,
			polarity:      domain.QuestionPolarityNegative,
			expectedScore: 4,
			expectedError: false,
		},
		{
			name:          "negative polarity - casi siempre (1→3)",
			selectedValue: 1,
			polarity:      domain.QuestionPolarityNegative,
			expectedScore: 3,
			expectedError: false,
		},
		{
			name:          "negative polarity - algunas veces (2→2)",
			selectedValue: 2,
			polarity:      domain.QuestionPolarityNegative,
			expectedScore: 2,
			expectedError: false,
		},
		{
			name:          "negative polarity - casi nunca (3→1)",
			selectedValue: 3,
			polarity:      domain.QuestionPolarityNegative,
			expectedScore: 1,
			expectedError: false,
		},
		{
			name:          "negative polarity - nunca (4→0)",
			selectedValue: 4,
			polarity:      domain.QuestionPolarityNegative,
			expectedScore: 0,
			expectedError: false,
		},
		// Invalid values
		{
			name:          "invalid value - negative",
			selectedValue: -1,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 0,
			expectedError: true,
		},
		{
			name:          "invalid value - too high",
			selectedValue: 5,
			polarity:      domain.QuestionPolarityPositive,
			expectedScore: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, err := ApplyPolarity(tt.selectedValue, tt.polarity)
			
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if score != tt.expectedScore {
				t.Errorf("expected score %d, got %d", tt.expectedScore, score)
			}
		})
	}
}

