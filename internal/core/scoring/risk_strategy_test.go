package scoring

import (
	"testing"

	"github.com/entorno35/backend/internal/domain"
)

// Helper function to create a mock question for Guide II/III with polarity
func createLikertQuestion(id uint, number int, guideType domain.GuideType, category, domainName, dimension string, polarity domain.QuestionPolarity) domain.Question {
	categoryID := uint(1) // Mock category ID
	domainID := uint(1)   // Mock domain ID
	dimensionID := uint(1) // Mock dimension ID
	polarityPtr := &polarity
	
	return domain.Question{
		ID:             id,
		QuestionNumber: number,
		GuideType:      guideType,
		Type:           domain.QuestionTypeLikert,
		Text:           "Mock question",
		Polarity:       polarityPtr,
		CategoryID:     &categoryID,
		DomainID:       &domainID,
		DimensionID:    &dimensionID,
		Category: domain.Category{
			ID:   categoryID,
			Name: category,
		},
		Domain: domain.Domain{
			ID:   domainID,
			Name: domainName,
		},
		Dimension: &domain.Dimension{
			ID:   dimensionID,
			Name: dimension,
		},
	}
}

// Helper function to create a response with calculated score
func createLikertResponse(id uint, questionID uint, selectedValue int, question domain.Question, calculatedScore int) domain.Response {
	return domain.Response{
		ID:              id,
		QuestionID:      questionID,
		SelectedValue:   selectedValue,
		CalculatedScore: calculatedScore,
		Question:        question,
	}
}

func TestRiskStrategy_Calculate_GuideII_FullAssessment(t *testing.T) {
	// Test Case A: Calculate a full assessment score for Guide II
	// Test Domain: "Carga de trabajo" with mixed polarity questions
	
	strategy := NewRiskStrategy(domain.GuideTypeII, LoadScoringRules())
	
	// Create questions for "Carga de trabajo" domain
	// According to official document, "Carga de trabajo" has negative polarity questions
	questions := []domain.Question{
		createLikertQuestion(1, 1, domain.GuideTypeII, "Factores propios de la actividad", "Carga de trabajo", "Dimension 1", domain.QuestionPolarityNegative),
		createLikertQuestion(2, 2, domain.GuideTypeII, "Factores propios de la actividad", "Carga de trabajo", "Dimension 1", domain.QuestionPolarityNegative),
		createLikertQuestion(3, 18, domain.GuideTypeII, "Factores propios de la actividad", "Carga de trabajo", "Dimension 1", domain.QuestionPolarityPositive),
	}
	
	// Create responses:
	// Q1: Selected 0 (Siempre) → Negative polarity → Calculated 4
	// Q2: Selected 2 (Algunas Veces) → Negative polarity → Calculated 2
	// Q3: Selected 4 (Nunca) → Positive polarity → Calculated 4
	responses := []domain.Response{
		createLikertResponse(1, 1, 0, questions[0], 4), // Negative: 0→4
		createLikertResponse(2, 2, 2, questions[1], 2), // Negative: 2→2
		createLikertResponse(3, 3, 4, questions[2], 4), // Positive: 4→4
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Expected total: 4 + 2 + 4 = 10
	expectedTotal := float64(10)
	if result.TotalScore != expectedTotal {
		t.Errorf("expected TotalScore to be %f, got %f", expectedTotal, result.TotalScore)
	}
	
	// Check risk level (Guide II total score threshold: 20, 45, 70, 90)
	// Score 10 < 20 → Nulo
	if result.RiskLevel != domain.RiskLevelNulo {
		t.Errorf("expected RiskLevel to be Nulo (score %f < 20), got %s", result.TotalScore, result.RiskLevel)
	}
	
	// Check domain score exists
	if len(result.DomainScores) == 0 {
		t.Error("expected DomainScores to be populated")
	}
}

func TestRiskStrategy_Calculate_GuideIII_MediumRisk(t *testing.T) {
	// Test Guide III with a score that falls in "Medio" range
	// Guide III thresholds: [50, 75, 99, 140]
	// Target: Score 80 → Medio
	
	strategy := NewRiskStrategy(domain.GuideTypeIII, LoadScoringRules())
	
	// Create enough questions to reach score ~80
	// Use "Carga de trabajo" domain (thresholds: [15, 21, 27, 37])
	questions := []domain.Question{}
	responses := []domain.Response{}
	
	// Create 20 questions with average score of 4 each = 80 total
	for i := 1; i <= 20; i++ {
		q := createLikertQuestion(uint(i), i, domain.GuideTypeIII, "Factores propios de la actividad", "Carga de trabajo", "Dimension 1", domain.QuestionPolarityNegative)
		questions = append(questions, q)
		
		// Selected value 0 (Siempre) → Negative polarity → Calculated 4
		r := createLikertResponse(uint(i), uint(i), 0, q, 4)
		responses = append(responses, r)
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Expected total: 20 * 4 = 80
	expectedTotal := float64(80)
	if result.TotalScore != expectedTotal {
		t.Errorf("expected TotalScore to be %f, got %f", expectedTotal, result.TotalScore)
	}
	
	// Score 80: 75 <= 80 < 99 → Medio
	if result.RiskLevel != domain.RiskLevelMedio {
		t.Errorf("expected RiskLevel to be Medio (score %f in range [75, 99)), got %s", result.TotalScore, result.RiskLevel)
	}
	
	// Check domain score for "Carga de trabajo"
	// Expected: 80 (all 20 questions in same domain)
	if domainScore, exists := result.DomainScores["Carga de trabajo"]; !exists {
		t.Error("expected DomainScores to contain 'Carga de trabajo'")
	} else if domainScore != expectedTotal {
		t.Errorf("expected domain score for 'Carga de trabajo' to be %f, got %f", expectedTotal, domainScore)
	}
}

func TestRiskStrategy_Calculate_PolarityInversion(t *testing.T) {
	// Test that polarity inversion is applied correctly
	// Positive polarity: 0→0, 1→1, 2→2, 3→3, 4→4
	// Negative polarity: 0→4, 1→3, 2→2, 3→1, 4→0
	
	strategy := NewRiskStrategy(domain.GuideTypeII, LoadScoringRules())
	
	// Create one positive and one negative question
	posQ := createLikertQuestion(1, 18, domain.GuideTypeII, "Factores propios de la actividad", "Carga de trabajo", "Dimension 1", domain.QuestionPolarityPositive)
	negQ := createLikertQuestion(2, 1, domain.GuideTypeII, "Factores propios de la actividad", "Carga de trabajo", "Dimension 1", domain.QuestionPolarityNegative)
	
	// Test selected value 0:
	// Positive: 0 → 0
	// Negative: 0 → 4
	responses := []domain.Response{
		createLikertResponse(1, 1, 0, posQ, 0), // Positive: 0→0
		createLikertResponse(2, 2, 0, negQ, 4), // Negative: 0→4
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Expected total: 0 + 4 = 4
	expectedTotal := float64(4)
	if result.TotalScore != expectedTotal {
		t.Errorf("expected TotalScore to be %f (0+4), got %f", expectedTotal, result.TotalScore)
	}
	
	// Test selected value 4:
	// Positive: 4 → 4
	// Negative: 4 → 0
	responses2 := []domain.Response{
		createLikertResponse(1, 1, 4, posQ, 4), // Positive: 4→4
		createLikertResponse(2, 2, 4, negQ, 0), // Negative: 4→0
	}
	
	result2, err := strategy.Calculate(responses2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Expected total: 4 + 0 = 4
	if result2.TotalScore != expectedTotal {
		t.Errorf("expected TotalScore to be %f (4+0), got %f", expectedTotal, result2.TotalScore)
	}
}

