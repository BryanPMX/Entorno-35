package domain

import (
	"time"

	"github.com/google/uuid"
)

// QuestionJSON represents the structure of questions in the JSON file
type QuestionJSON struct {
	Number      int    `json:"number"`
	Text        string `json:"text"`
	Type        string `json:"type"`
	Section     string `json:"section,omitempty"`      // Guide I only
	Subsection  string `json:"subsection,omitempty"`   // Guide I only
	Category    string `json:"category,omitempty"`     // Guide II/III only
	Domain      string `json:"domain,omitempty"`       // Guide II/III only
	Dimension   string `json:"dimension,omitempty"`    // Guide II/III only
	Polarity    string `json:"polarity,omitempty"`     // Guide II/III only
}

// GuideJSON represents a guide structure in the JSON file
type GuideJSON struct {
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	TotalQuestions int            `json:"total_questions"`
	Questions      []QuestionJSON `json:"questions"`
}

// QuestionsDataJSON represents the root structure of nom035_questions.json
type QuestionsDataJSON struct {
	GuideI   GuideJSON `json:"guide_i"`
	GuideII  GuideJSON `json:"guide_ii"`
	GuideIII GuideJSON `json:"guide_iii"`
}

// IndividualReportDTO represents the data structure for an individual assessment report
type IndividualReportDTO struct {
	AssessmentID    uuid.UUID            `json:"assessment_id"`
	Period          int                  `json:"period"`
	GuideType       GuideType            `json:"guide_type"`
	StaffName       string               `json:"staff_name"`
	Department      string               `json:"department,omitempty"`
	Shift           string               `json:"shift,omitempty"`
	TotalScore      float64              `json:"total_score"`
	RiskLevel       RiskLevel            `json:"risk_level"`
	CategoryScores  map[string]float64   `json:"category_scores"`
	DomainScores    map[string]float64   `json:"domain_scores"`
	RequiresMedical bool                 `json:"requires_medical_attention"`
	CompletedAt     *time.Time           `json:"completed_at,omitempty"`
	Recommendations []string             `json:"recommendations,omitempty"`
}

// RiskDistribution represents the count of assessments by risk level
type RiskDistribution struct {
	RiskLevel RiskLevel `json:"risk_level"`
	Count     int64     `json:"count"`
}

// DepartmentRiskHeatmap represents risk distribution by department
type DepartmentRiskHeatmap struct {
	Department string     `json:"department"`
	RiskLevel  RiskLevel  `json:"risk_level"`
	Count      int64      `json:"count"`
}

// GeneralReportDTO represents the data structure for a company-wide general report
type GeneralReportDTO struct {
	CompanyID            uuid.UUID              `json:"company_id"`
	CompanyName          string                  `json:"company_name"`
	Period               *int                    `json:"period,omitempty"`
	TotalStaff            int64                   `json:"total_staff"`
	CompletedAssessments int64                   `json:"completed_assessments"`
	ParticipationRate    float64                 `json:"participation_rate"` // percentage
	RiskDistribution     []RiskDistribution      `json:"risk_distribution"`
	DepartmentHeatmap    []DepartmentRiskHeatmap `json:"department_heatmap"`
}

