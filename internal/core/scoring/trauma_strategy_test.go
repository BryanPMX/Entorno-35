package scoring

import (
	"testing"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// Helper function to create a mock question for Guide I
func createGuideIQuestion(number int, section, subsection string) domain.Question {
	sectionPtr := &section
	subsectionPtr := &subsection
	return domain.Question{
		ID:             uint(number),
		QuestionNumber: number,
		GuideType:      domain.GuideTypeI,
		Type:           domain.QuestionTypeBinary,
		Text:           "Mock question",
		Section:        sectionPtr,
		Subsection:     subsectionPtr,
	}
}

// Helper function to create a mock response
func createResponse(questionID uint, selectedValue int, question domain.Question) domain.Response {
	return domain.Response{
		ID:             uint(questionID),
		QuestionID:     questionID,
		SelectedValue:  selectedValue, // 0 = No, 1 = Yes for binary
		CalculatedScore: selectedValue, // For binary, calculated = selected
		Question:       question,
	}
}

func TestTraumaStrategy_Calculate_SectionI_AllNo(t *testing.T) {
	// Test Case A: All "No" to Section I questions
	// Expected: Null Risk, RequiresMedicalAttention: false
	
	strategy := NewTraumaStrategy()
	
	// Create Section I questions (first 3 questions are typically Section I)
	questions := []domain.Question{
		createGuideIQuestion(1, "I", "Acontecimiento Traumático Severo"),
		createGuideIQuestion(2, "I", "Acontecimiento Traumático Severo"),
		createGuideIQuestion(3, "I", "Acontecimiento Traumático Severo"),
	}
	
	// All answers are "No" (0)
	responses := []domain.Response{
		createResponse(1, 0, questions[0]),
		createResponse(2, 0, questions[1]),
		createResponse(3, 0, questions[2]),
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result.RequiresMedicalAttention {
		t.Error("expected RequiresMedicalAttention to be false when all Section I answers are No")
	}
	
	if result.RiskLevel != domain.RiskLevelNulo {
		t.Errorf("expected RiskLevel to be Nulo, got %s", result.RiskLevel)
	}
	
	if result.TotalScore != 0 {
		t.Errorf("expected TotalScore to be 0, got %f", result.TotalScore)
	}
}

func TestTraumaStrategy_Calculate_SectionII_OneYes(t *testing.T) {
	// Test Case B: "Yes" to 1 Section II question
	// Expected: RequiresMedicalAttention: true
	
	strategy := NewTraumaStrategy()
	
	// Create Section I and Section II questions
	questions := []domain.Question{
		createGuideIQuestion(1, "I", "Acontecimiento Traumático Severo"),
		createGuideIQuestion(2, "I", "Acontecimiento Traumático Severo"),
		createGuideIQuestion(4, "II", "Recuerdos persistentes"),
		createGuideIQuestion(5, "II", "Recuerdos persistentes"),
	}
	
	// Section I: At least one "Yes" (triggers Section II)
	// Section II: One "Yes"
	responses := []domain.Response{
		createResponse(1, 1, questions[0]), // Yes to Section I
		createResponse(2, 0, questions[1]), // No
		createResponse(4, 1, questions[2]), // Yes to Section II
		createResponse(5, 0, questions[3]), // No
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.RequiresMedicalAttention {
		t.Error("expected RequiresMedicalAttention to be true when at least one Section II question is Yes")
	}
}

func TestTraumaStrategy_Calculate_SectionIII_ThreeYes(t *testing.T) {
	// Test Case C: "Yes" to 3 Section III questions
	// Expected: RequiresMedicalAttention: true
	
	strategy := NewTraumaStrategy()
	
	// Create questions for Section I, II, and III
	questions := []domain.Question{
		createGuideIQuestion(1, "I", "Acontecimiento Traumático Severo"),
		createGuideIQuestion(4, "II", "Recuerdos persistentes"),
		createGuideIQuestion(8, "III", "Esfuerzo por evitar"),
		createGuideIQuestion(9, "III", "Esfuerzo por evitar"),
		createGuideIQuestion(10, "III", "Esfuerzo por evitar"),
		createGuideIQuestion(11, "III", "Esfuerzo por evitar"),
	}
	
	// Section I: Yes (triggers Section II)
	// Section II: Yes (triggers Section III)
	// Section III: 3 Yes answers
	responses := []domain.Response{
		createResponse(1, 1, questions[0]),  // Yes to Section I
		createResponse(4, 1, questions[1]),  // Yes to Section II
		createResponse(8, 1, questions[2]),  // Yes to Section III
		createResponse(9, 1, questions[3]),  // Yes to Section III
		createResponse(10, 1, questions[4]), // Yes to Section III
		createResponse(11, 0, questions[5]), // No to Section III
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.RequiresMedicalAttention {
		t.Error("expected RequiresMedicalAttention to be true when 3+ Section III questions are Yes")
	}
}

func TestTraumaStrategy_Calculate_SectionIV_Triggers(t *testing.T) {
	// Additional test: Section IV triggers
	// Section IV (Afectación) is triggered if Section III has "Yes" answers
	
	strategy := NewTraumaStrategy()
	
	questions := []domain.Question{
		createGuideIQuestion(1, "I", "Acontecimiento Traumático Severo"),
		createGuideIQuestion(4, "II", "Recuerdos persistentes"),
		createGuideIQuestion(8, "III", "Esfuerzo por evitar"),
		createGuideIQuestion(13, "IV", "Afectación"),
	}
	
	// Trigger all sections
	responses := []domain.Response{
		createResponse(1, 1, questions[0]),  // Yes to Section I
		createResponse(4, 1, questions[1]),  // Yes to Section II
		createResponse(8, 1, questions[2]),  // Yes to Section III
		createResponse(13, 1, questions[3]), // Yes to Section IV
	}
	
	result, err := strategy.Calculate(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.RequiresMedicalAttention {
		t.Error("expected RequiresMedicalAttention to be true when Section IV has Yes answers")
	}
}

// Helper to avoid unused import warning
var _ = uuid.Nil

