package scoring

import "testing"

func TestRiskThresholds_GetRiskLevel(t *testing.T) {
	thresholds := RiskThresholds{Ranges: [4]float64{20, 45, 70, 90}}

	tests := []struct {
		name         string
		score        float64
		expected     RiskLevel
		description  string
	}{
		// Nulo: Score < Limit1 (20)
		{
			name:        "nulo - below first threshold",
			score:       0,
			expected:    RiskLevelNulo,
			description: "Score 0 should be Nulo",
		},
		{
			name:        "nulo - just below first threshold",
			score:       19.9,
			expected:    RiskLevelNulo,
			description: "Score 19.9 should be Nulo",
		},
		// Bajo: Limit1 (20) <= Score < Limit2 (45)
		{
			name:        "bajo - at first threshold",
			score:       20,
			expected:    RiskLevelBajo,
			description: "Score 20 should be Bajo (boundary)",
		},
		{
			name:        "bajo - middle of range",
			score:       32.5,
			expected:    RiskLevelBajo,
			description: "Score 32.5 should be Bajo",
		},
		{
			name:        "bajo - just below second threshold",
			score:       44.9,
			expected:    RiskLevelBajo,
			description: "Score 44.9 should be Bajo",
		},
		// Medio: Limit2 (45) <= Score < Limit3 (70)
		{
			name:        "medio - at second threshold",
			score:       45,
			expected:    RiskLevelMedio,
			description: "Score 45 should be Medio (boundary)",
		},
		{
			name:        "medio - middle of range",
			score:       57.5,
			expected:    RiskLevelMedio,
			description: "Score 57.5 should be Medio",
		},
		{
			name:        "medio - just below third threshold",
			score:       69.9,
			expected:    RiskLevelMedio,
			description: "Score 69.9 should be Medio",
		},
		// Alto: Limit3 (70) <= Score < Limit4 (90)
		{
			name:        "alto - at third threshold",
			score:       70,
			expected:    RiskLevelAlto,
			description: "Score 70 should be Alto (boundary)",
		},
		{
			name:        "alto - middle of range",
			score:       80,
			expected:    RiskLevelAlto,
			description: "Score 80 should be Alto",
		},
		{
			name:        "alto - just below fourth threshold",
			score:       89.9,
			expected:    RiskLevelAlto,
			description: "Score 89.9 should be Alto",
		},
		// Muy Alto: Score >= Limit4 (90)
		{
			name:        "muy_alto - at fourth threshold",
			score:       90,
			expected:    RiskLevelMuyAlto,
			description: "Score 90 should be Muy Alto (boundary)",
		},
		{
			name:        "muy_alto - above threshold",
			score:       100,
			expected:    RiskLevelMuyAlto,
			description: "Score 100 should be Muy Alto",
		},
		{
			name:        "muy_alto - very high",
			score:       200,
			expected:    RiskLevelMuyAlto,
			description: "Score 200 should be Muy Alto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := thresholds.GetRiskLevel(tt.score)
			if result != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.description, tt.expected, result)
			}
		})
	}
}

func TestRiskThresholds_GuideIII(t *testing.T) {
	// Test Guide III thresholds (50, 75, 99, 140)
	thresholds := RiskThresholds{Ranges: [4]float64{50, 75, 99, 140}}

	tests := []struct {
		name     string
		score    float64
		expected RiskLevel
	}{
		{"nulo boundary", 49.9, RiskLevelNulo},
		{"bajo lower boundary", 50, RiskLevelBajo},
		{"bajo upper boundary", 74.9, RiskLevelBajo},
		{"medio lower boundary", 75, RiskLevelMedio},
		{"medio upper boundary", 98.9, RiskLevelMedio},
		{"alto lower boundary", 99, RiskLevelAlto},
		{"alto upper boundary", 139.9, RiskLevelAlto},
		{"muy_alto lower boundary", 140, RiskLevelMuyAlto},
		{"muy_alto high", 200, RiskLevelMuyAlto},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := thresholds.GetRiskLevel(tt.score)
			if result != tt.expected {
				t.Errorf("score %.1f: expected %s, got %s", tt.score, tt.expected, result)
			}
		})
	}
}

