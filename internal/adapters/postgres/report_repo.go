package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrReportNotFound = errors.New("report not found")
)

// reportRepository implements the ReportRepository interface using GORM
type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new Postgres implementation of ReportRepository
func NewReportRepository(db *gorm.DB) ports.ReportRepository {
	return &reportRepository{db: db}
}

// GetIndividualReport retrieves data for an individual assessment report
func (r *reportRepository) GetIndividualReport(assessmentID uuid.UUID, companyID uuid.UUID) (*domain.IndividualReportDTO, error) {
	var assessment domain.Assessment

	// Fetch assessment with staff and company preloaded
	result := r.db.
		Preload("Staff").
		Preload("Company").
		Where("id = ? AND company_id = ?", assessmentID, companyID).
		First(&assessment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrReportNotFound
		}
		return nil, fmt.Errorf("failed to fetch assessment: %w", result.Error)
	}

	// Check if assessment has been scored
	if assessment.TotalScore == nil || assessment.RiskLevel == nil {
		return nil, fmt.Errorf("assessment has not been scored yet")
	}

	// Fetch responses with questions and relationships loaded for score recalculation
	responses, err := r.getResponsesWithQuestions(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch responses: %w", err)
	}

	// Recalculate category and domain scores from responses
	categoryScores, domainScores, err := r.calculateScoresFromResponses(responses, assessment.GuideType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate scores: %w", err)
	}

	// Extract demographics
	department := ""
	shift := ""
	if assessment.Staff.Demographics.Department != "" {
		department = assessment.Staff.Demographics.Department
	}
	if assessment.Staff.Demographics.ShiftType != "" {
		shift = assessment.Staff.Demographics.ShiftType
	}

	// Build DTO
	dto := &domain.IndividualReportDTO{
		AssessmentID:    assessment.ID,
		Period:          assessment.Period,
		GuideType:       assessment.GuideType,
		StaffName:       assessment.Staff.FullName,
		Department:      department,
		Shift:           shift,
		TotalScore:      *assessment.TotalScore,
		RiskLevel:       *assessment.RiskLevel,
		CategoryScores:  categoryScores,
		DomainScores:    domainScores,
		RequiresMedical: assessment.RequiresMedicalAttention,
		CompletedAt:    assessment.CompletedAt,
	}

	return dto, nil
}

// GetGeneralReport retrieves aggregated data for a company-wide general report
func (r *reportRepository) GetGeneralReport(companyID uuid.UUID, period *int) (*domain.GeneralReportDTO, error) {
	// Fetch company
	var company domain.Company
	if err := r.db.Where("id = ?", companyID).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("failed to fetch company: %w", err)
	}

	// Aggregation 1: Count total staff
	var totalStaff int64
	staffQuery := r.db.Model(&domain.Staff{}).Where("company_id = ?", companyID)
	if err := staffQuery.Count(&totalStaff).Error; err != nil {
		return nil, fmt.Errorf("failed to count staff: %w", err)
	}

	// Aggregation 2: Count completed assessments
	var completedAssessments int64
	assessmentQuery := r.db.Model(&domain.Assessment{}).
		Where("company_id = ? AND status = ?", companyID, domain.AssessmentStatusCompleted)
	
	if period != nil {
		assessmentQuery = assessmentQuery.Where("period = ?", *period)
	}

	if err := assessmentQuery.Count(&completedAssessments).Error; err != nil {
		return nil, fmt.Errorf("failed to count assessments: %w", err)
	}

	// Calculate participation rate
	participationRate := float64(0)
	if totalStaff > 0 {
		participationRate = (float64(completedAssessments) / float64(totalStaff)) * 100
	}

	// Aggregation 3: Risk distribution (GROUP BY risk_level)
	var riskDistributions []struct {
		RiskLevel domain.RiskLevel `gorm:"column:risk_level"`
		Count     int64            `gorm:"column:count"`
	}

	riskQuery := r.db.Model(&domain.Assessment{}).
		Select("risk_level, COUNT(*) as count").
		Where("company_id = ? AND status = ? AND risk_level IS NOT NULL", companyID, domain.AssessmentStatusCompleted)
	
	if period != nil {
		riskQuery = riskQuery.Where("period = ?", *period)
	}

	if err := riskQuery.
		Group("risk_level").
		Scan(&riskDistributions).Error; err != nil {
		return nil, fmt.Errorf("failed to get risk distribution: %w", err)
	}

	// Convert to DTO format
	riskDistribution := make([]domain.RiskDistribution, 0, len(riskDistributions))
	for _, rd := range riskDistributions {
		riskDistribution = append(riskDistribution, domain.RiskDistribution{
			RiskLevel: rd.RiskLevel,
			Count:     rd.Count,
		})
	}

	// Aggregation 4: Department heatmap (JOIN assessments + staff, GROUP BY department AND risk_level)
	var heatmapData []struct {
		Department string          `gorm:"column:department"`
		RiskLevel  domain.RiskLevel `gorm:"column:risk_level"`
		Count      int64            `gorm:"column:count"`
	}

	heatmapQuery := r.db.Table("assessments").
		Select("staff.demographics->>'department' as department, assessments.risk_level, COUNT(*) as count").
		Joins("INNER JOIN staff ON assessments.staff_id = staff.id").
		Where("assessments.company_id = ? AND assessments.status = ? AND assessments.risk_level IS NOT NULL", companyID, domain.AssessmentStatusCompleted).
		Where("staff.demographics->>'department' IS NOT NULL AND staff.demographics->>'department' != ''")
	
	if period != nil {
		heatmapQuery = heatmapQuery.Where("assessments.period = ?", *period)
	}

	if err := heatmapQuery.
		Group("staff.demographics->>'department', assessments.risk_level").
		Scan(&heatmapData).Error; err != nil {
		return nil, fmt.Errorf("failed to get department heatmap: %w", err)
	}

	// Convert to DTO format
	departmentHeatmap := make([]domain.DepartmentRiskHeatmap, 0, len(heatmapData))
	for _, hd := range heatmapData {
		departmentHeatmap = append(departmentHeatmap, domain.DepartmentRiskHeatmap{
			Department: hd.Department,
			RiskLevel:  hd.RiskLevel,
			Count:      hd.Count,
		})
	}

	// Build DTO
	dto := &domain.GeneralReportDTO{
		CompanyID:             company.ID,
		CompanyName:           company.Name,
		Period:                period,
		TotalStaff:            totalStaff,
		CompletedAssessments: completedAssessments,
		ParticipationRate:     participationRate,
		RiskDistribution:      riskDistribution,
		DepartmentHeatmap:     departmentHeatmap,
	}

	return dto, nil
}

// getResponsesWithQuestions fetches responses with questions and relationships loaded
func (r *reportRepository) getResponsesWithQuestions(assessmentID uuid.UUID) ([]domain.Response, error) {
	var responses []domain.Response

	result := r.db.
		Preload("Question").
		Preload("Question.Category").
		Preload("Question.Domain").
		Preload("Question.Dimension").
		Where("assessment_id = ?", assessmentID).
		Find(&responses)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch responses: %w", result.Error)
	}

	return responses, nil
}

// calculateScoresFromResponses recalculates category and domain scores from responses
// This mirrors the logic in risk_strategy.go but only calculates category/domain scores
func (r *reportRepository) calculateScoresFromResponses(responses []domain.Response, guideType domain.GuideType) (map[string]float64, map[string]float64, error) {
	categoryScores := make(map[string]float64)
	domainScores := make(map[string]float64)

	// Only calculate for Guide II/III (Guide I doesn't have category/domain structure)
	if guideType == domain.GuideTypeI {
		return categoryScores, domainScores, nil
	}

	// Process each response
	for _, response := range responses {
		if response.Question.ID == 0 {
			continue
		}

		question := response.Question

		// Apply polarity inversion (same logic as scoring strategy)
		if question.Polarity == nil {
			continue
		}

		calculatedScore := float64(response.CalculatedScore)

		// Aggregate by domain
		if question.DomainID != nil && question.Domain != nil {
			domainName := question.Domain.Name
			domainScores[domainName] += calculatedScore
		}

		// Aggregate by category
		if question.CategoryID != nil && question.Category != nil {
			categoryName := question.Category.Name
			categoryScores[categoryName] += calculatedScore
		}
	}

	return categoryScores, domainScores, nil
}

