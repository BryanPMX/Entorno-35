package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/core/scoring"
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

	// Fetch responses with questions and relationships loaded for score calculation
	// Note: Category and domain scores are calculated deterministically from stored responses
	// This ensures consistency across report views for the same assessment
	responses, err := r.getResponsesWithQuestions(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch responses: %w", err)
	}

	// Calculate category and domain scores from stored responses
	// This is deterministic and ensures report consistency
	categoryScores, domainScores, err := r.calculateScoresFromResponses(responses, assessment.GuideType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate scores: %w", err)
	}

	// Calculate dynamic maximum scores based on questions answered
	categoryMaxScores, domainMaxScores, err := r.calculateDynamicMaxScores(responses, assessment.GuideType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate max scores: %w", err)
	}

	// Validate scores don't exceed maximums (data integrity check)
	categoryScores, domainScores = r.validateAndClampScores(categoryScores, domainScores, categoryMaxScores, domainMaxScores)

	// Calculate category and domain risk levels using NOM-035 thresholds
	categoryRiskLevels, domainRiskLevels, err := r.calculateRiskLevelsFromScores(categoryScores, domainScores, assessment.GuideType, categoryMaxScores, domainMaxScores)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate risk levels: %w", err)
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

	// Calculate total max score from questions answered
	totalMaxScore := float64(len(responses) * scoring.MaxScorePerQuestion)

	// Build question response details for detailed analysis
	questionResponses := make([]domain.QuestionResponseDetail, 0, len(responses))
	for _, response := range responses {
		if response.Question.ID == 0 {
			continue
		}

		question := response.Question
		
		// Calculate score with polarity
		calculatedScore := response.SelectedValue
		if question.Polarity != nil {
			calcScore, err := scoring.ApplyPolarity(response.SelectedValue, *question.Polarity)
			if err == nil {
				calculatedScore = calcScore
			}
		}

		// Get category and domain names
		categoryName := ""
		domainName := ""
		if question.Category != nil {
			categoryName = question.Category.Name
		}
		if question.Domain != nil {
			domainName = question.Domain.Name
		}

		detail := domain.QuestionResponseDetail{
			QuestionNumber:  question.QuestionNumber,
			QuestionText:    question.Text,
			Category:        categoryName,
			Domain:          domainName,
			SelectedValue:   response.SelectedValue,
			CalculatedScore: calculatedScore,
			MaxScore:        scoring.MaxScorePerQuestion,
		}
		questionResponses = append(questionResponses, detail)
	}

	// Build DTO
	dto := &domain.IndividualReportDTO{
		AssessmentID:       assessment.ID,
		Period:             assessment.Period,
		GuideType:          assessment.GuideType,
		StaffName:          assessment.Staff.FullName,
		Department:         department,
		Shift:              shift,
		TotalScore:         *assessment.TotalScore,
		TotalMaxScore:      totalMaxScore,
		RiskLevel:          *assessment.RiskLevel,
		CategoryScores:     categoryScores,
		CategoryRiskLevels: categoryRiskLevels,
		CategoryMaxScores:  categoryMaxScores,
		DomainScores:       domainScores,
		DomainRiskLevels:   domainRiskLevels,
		DomainMaxScores:    domainMaxScores,
		RequiresMedical:    assessment.RequiresMedicalAttention,
		CompletedAt:        assessment.CompletedAt,
		QuestionResponses:  questionResponses,
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
	// Also calculate average total_score per department for accurate score display
	var heatmapData []struct {
		Department string           `gorm:"column:department"`
		RiskLevel  domain.RiskLevel `gorm:"column:risk_level"`
		Count      int64            `gorm:"column:count"`
	}

	var avgScoreData []struct {
		Department string  `gorm:"column:department"`
		AvgScore   float64 `gorm:"column:avg_score"`
		Count      int64   `gorm:"column:count"`
	}

	heatmapQuery := r.db.Model(&domain.Assessment{}).
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

	// Calculate average total_score per department
	avgScoreQuery := r.db.Model(&domain.Assessment{}).
		Select("staff.demographics->>'department' as department, AVG(assessments.total_score) as avg_score, COUNT(*) as count").
		Joins("INNER JOIN staff ON assessments.staff_id = staff.id").
		Where("assessments.company_id = ? AND assessments.status = ? AND assessments.total_score IS NOT NULL", companyID, domain.AssessmentStatusCompleted).
		Where("staff.demographics->>'department' IS NOT NULL AND staff.demographics->>'department' != ''")

	if period != nil {
		avgScoreQuery = avgScoreQuery.Where("assessments.period = ?", *period)
	}

	if err := avgScoreQuery.
		Group("staff.demographics->>'department'").
		Scan(&avgScoreData).Error; err != nil {
		return nil, fmt.Errorf("failed to get department average scores: %w", err)
	}

	// Build a map of department -> average score
	avgScoreMap := make(map[string]float64)
	for _, avg := range avgScoreData {
		avgScoreMap[avg.Department] = avg.AvgScore
	}

	// Convert to DTO format
	departmentHeatmap := make([]domain.DepartmentRiskHeatmap, 0, len(heatmapData))
	for _, hd := range heatmapData {
		dept := domain.DepartmentRiskHeatmap{
			Department: hd.Department,
			RiskLevel:  hd.RiskLevel,
			Count:      hd.Count,
		}
		// Add average score if available
		if avgScore, exists := avgScoreMap[hd.Department]; exists {
			dept.AvgScore = &avgScore
		}
		departmentHeatmap = append(departmentHeatmap, dept)
	}

	// Aggregation 5: Age distribution (GROUP BY age_range) - based on staff with completed assessments
	ageDistribution := r.getDemographicDistribution(companyID, "age_range", period)

	// Aggregation 6: Marital status distribution - based on staff with completed assessments
	maritalStatusDistribution := r.getDemographicDistribution(companyID, "marital_status", period)

	// Aggregation 7: Shift type distribution - based on staff with completed assessments
	shiftTypeDistribution := r.getDemographicDistribution(companyID, "shift_type", period)

	// Aggregation 8: Experience distribution - based on staff with completed assessments
	experienceDistribution := r.getDemographicDistribution(companyID, "total_work_experience", period)

	// Aggregation 9: Age risk distribution (cross-analysis)
	ageRiskDistribution := r.getDemographicRiskDistribution(companyID, "age_range", period)

	// Aggregation 10: Shift risk distribution (cross-analysis)
	shiftRiskDistribution := r.getDemographicRiskDistribution(companyID, "shift_type", period)

	// Build DTO
	dto := &domain.GeneralReportDTO{
		CompanyID:                  company.ID,
		CompanyName:                company.Name,
		Period:                     period,
		TotalStaff:                 totalStaff,
		CompletedAssessments:       completedAssessments,
		ParticipationRate:          participationRate,
		RiskDistribution:           riskDistribution,
		DepartmentHeatmap:          departmentHeatmap,
		AgeDistribution:            ageDistribution,
		MaritalStatusDistribution:  maritalStatusDistribution,
		ShiftTypeDistribution:      shiftTypeDistribution,
		ExperienceDistribution:     experienceDistribution,
		AgeRiskDistribution:        ageRiskDistribution,
		ShiftRiskDistribution:      shiftRiskDistribution,
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
// IMPORTANT: Always recalculates from SelectedValue and Polarity to ensure accuracy
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

		// CRITICAL FIX: Recalculate score from SelectedValue and Polarity
		// DO NOT use stored CalculatedScore as it may be from old buggy logic
		calculatedScore, err := scoring.ApplyPolarity(response.SelectedValue, *question.Polarity)
		if err != nil {
			continue // Skip invalid responses
		}
		calculatedScoreFloat := float64(calculatedScore)

		// Aggregate by domain
		if question.DomainID != nil && question.Domain != nil {
			domainName := question.Domain.Name
			domainScores[domainName] += calculatedScoreFloat
		}

		// Aggregate by category
		if question.CategoryID != nil && question.Category != nil {
			categoryName := question.Category.Name
			categoryScores[categoryName] += calculatedScoreFloat
		}
	}

	return categoryScores, domainScores, nil
}

// calculateDynamicMaxScores calculates the maximum possible scores for categories and domains
// based on the questions that were actually answered, not hardcoded NOM-035 specifications
func (r *reportRepository) calculateDynamicMaxScores(responses []domain.Response, guideType domain.GuideType) (map[string]float64, map[string]float64, error) {
	categoryMaxScores := make(map[string]float64)
	domainMaxScores := make(map[string]float64)

	// Only calculate for Guide II/III (Guide I doesn't have category/domain structure)
	if guideType == domain.GuideTypeI {
		return categoryMaxScores, domainMaxScores, nil
	}

	// Track unique questions per category/domain to avoid double-counting
	categoryQuestionCounts := make(map[string]int)
	domainQuestionCounts := make(map[string]int)

	// Process each response to count questions per category/domain
	for _, response := range responses {
		if response.Question.ID == 0 {
			continue
		}

		question := response.Question

		// Skip questions without polarity (shouldn't happen in Guide II/III)
		if question.Polarity == nil {
			continue
		}

		// Count by domain
		if question.DomainID != nil && question.Domain != nil {
			domainName := question.Domain.Name
			if _, exists := domainQuestionCounts[domainName]; !exists {
				domainQuestionCounts[domainName] = 0
			}
			domainQuestionCounts[domainName]++
		}

		// Count by category
		if question.CategoryID != nil && question.Category != nil {
			categoryName := question.Category.Name
			if _, exists := categoryQuestionCounts[categoryName]; !exists {
				categoryQuestionCounts[categoryName] = 0
			}
			categoryQuestionCounts[categoryName]++
		}
	}

	// Calculate maximum scores based on question counts
	// Each question can contribute a maximum of MaxScorePerQuestion points (after polarity adjustment)
	for categoryName, questionCount := range categoryQuestionCounts {
		categoryMaxScores[categoryName] = float64(questionCount * scoring.MaxScorePerQuestion)
	}

	for domainName, questionCount := range domainQuestionCounts {
		domainMaxScores[domainName] = float64(questionCount * scoring.MaxScorePerQuestion)
	}

	return categoryMaxScores, domainMaxScores, nil
}

// validateAndClampScores ensures scores don't exceed their calculated maximums
// This prevents impossible scores like 40/38 and provides data integrity
func (r *reportRepository) validateAndClampScores(
	categoryScores, domainScores, categoryMaxScores, domainMaxScores map[string]float64,
) (map[string]float64, map[string]float64) {
	// Validate category scores
	for category, score := range categoryScores {
		if maxScore, exists := categoryMaxScores[category]; exists {
			if score > maxScore {
				// Log warning but don't expose sensitive data (no staff/assessment IDs)
				fmt.Printf("WARNING: Category '%s' score %.1f exceeds maximum %.1f, clamping to maximum\n", category, score, maxScore)
				categoryScores[category] = maxScore
			}
		}
	}

	// Validate domain scores
	for domain, score := range domainScores {
		if maxScore, exists := domainMaxScores[domain]; exists {
			if score > maxScore {
				// Log warning but don't expose sensitive data
				fmt.Printf("WARNING: Domain '%s' score %.1f exceeds maximum %.1f, clamping to maximum\n", domain, score, maxScore)
				domainScores[domain] = maxScore
			}
		}
	}

	return categoryScores, domainScores
}

// calculateRiskLevelsFromScores calculates risk levels for categories and domains using NOM-035 thresholds
func (r *reportRepository) calculateRiskLevelsFromScores(categoryScores, domainScores map[string]float64, guideType domain.GuideType, categoryMaxScores, domainMaxScores map[string]float64) (map[string]string, map[string]string, error) {
	categoryRiskLevels := make(map[string]string)
	domainRiskLevels := make(map[string]string)

	// Load scoring rules
	rules := scoring.LoadScoringRules()
	var guideRules interface{}

	if guideType == domain.GuideTypeII {
		guideRules = rules.GuideII
	} else if guideType == domain.GuideTypeIII {
		guideRules = rules.GuideIII
	} else {
		// Guide I doesn't have category/domain structure
		return categoryRiskLevels, domainRiskLevels, nil
	}

	// Calculate category risk levels
	if guideII, ok := guideRules.(scoring.GuideIIScoringRules); ok {
		for category, score := range categoryScores {
			// Try exact match first
			if thresholds, exists := guideII.Categories[category]; exists {
				categoryRiskLevels[category] = string(thresholds.GetRiskLevel(score))
			} else {
				// Fallback: always calculate based on score magnitude
				// Since NOM-035 scores are accumulated, use score-based logic
				categoryRiskLevels[category] = r.calculateRiskFromScore(score)
			}
		}
	} else if guideIII, ok := guideRules.(scoring.GuideIIIScoringRules); ok {
		for category, score := range categoryScores {
			if thresholds, exists := guideIII.Categories[category]; exists {
				categoryRiskLevels[category] = string(thresholds.GetRiskLevel(score))
			} else {
				categoryRiskLevels[category] = r.calculateRiskFromScore(score)
			}
		}
	} else {
		// Fallback for any guide type
		for category, score := range categoryScores {
			categoryRiskLevels[category] = r.calculateRiskFromScore(score)
		}
	}

	// Calculate domain risk levels
	if guideII, ok := guideRules.(scoring.GuideIIScoringRules); ok {
		for domain, score := range domainScores {
			if thresholds, exists := guideII.Domains[domain]; exists {
				domainRiskLevels[domain] = string(thresholds.GetRiskLevel(score))
			} else {
				// Fallback to percentage-based calculation using actual max score
				maxScore := domainMaxScores[domain]
				if maxScore == 0 {
					maxScore = 15 // Safety fallback if max score is unavailable
				}
				domainRiskLevels[domain] = string(r.getRiskLevelFromPercentage(score, maxScore))
			}
		}
	} else if guideIII, ok := guideRules.(scoring.GuideIIIScoringRules); ok {
		for domain, score := range domainScores {
			if thresholds, exists := guideIII.Domains[domain]; exists {
				domainRiskLevels[domain] = string(thresholds.GetRiskLevel(score))
			} else {
				// Fallback to percentage-based calculation using actual max score
				maxScore := domainMaxScores[domain]
				if maxScore == 0 {
					maxScore = 15 // Safety fallback if max score is unavailable
				}
				domainRiskLevels[domain] = string(r.getRiskLevelFromPercentage(score, maxScore))
			}
		}
	}

	return categoryRiskLevels, domainRiskLevels, nil
}

// calculateRiskFromScore provides intelligent risk level calculation based on score magnitude
// This is used when specific NOM-035 thresholds are not available
func (r *reportRepository) calculateRiskFromScore(score float64) string {
	// For NOM-035, higher accumulated scores indicate higher risk
	// Use reasonable thresholds based on typical score ranges
	if score <= 5 {
		return "nulo"
	}
	if score <= 15 {
		return "bajo"
	}
	if score <= 30 {
		return "medio"
	}
	if score <= 50 {
		return "alto"
	}
	return "muy_alto"
}

// getRiskLevelFromPercentage provides fallback risk level calculation based on percentage
func (r *reportRepository) getRiskLevelFromPercentage(score, maxScore float64) string {
	percentage := (score / maxScore) * 100
	if percentage < 25 {
		return "nulo"
	}
	if percentage < 50 {
		return "bajo"
	}
	if percentage < 75 {
		return "medio"
	}
	if percentage < 90 {
		return "alto"
	}
	return "muy_alto"
}

// getDemographicDistribution retrieves the distribution of staff by a demographic field
// Filters by completed assessments to ensure data reflects actual assessment participants
// Uses Model to ensure soft deletes are properly handled
func (r *reportRepository) getDemographicDistribution(companyID uuid.UUID, demographicField string, period *int) []domain.DemographicDistribution {
	var results []struct {
		Category string `gorm:"column:category"`
		Count    int64  `gorm:"column:count"`
	}

	demographicQuery := fmt.Sprintf("staff.demographics->>'%s'", demographicField)
	query := r.db.Model(&domain.Assessment{}).
		Select(demographicQuery + " as category, COUNT(DISTINCT assessments.staff_id) as count").
		Joins("INNER JOIN staff ON assessments.staff_id = staff.id").
		Where("assessments.company_id = ? AND assessments.status = ?", companyID, domain.AssessmentStatusCompleted).
		Where(demographicQuery + " IS NOT NULL AND " + demographicQuery + " != ''")

	if period != nil {
		query = query.Where("assessments.period = ?", *period)
	}

	if err := query.
		Group(demographicQuery).
		Order("count DESC").
		Scan(&results).Error; err != nil {
		// Return empty distribution on error rather than failing the entire report
		return []domain.DemographicDistribution{}
	}

	distribution := make([]domain.DemographicDistribution, 0, len(results))
	for _, res := range results {
		distribution = append(distribution, domain.DemographicDistribution{
			Category: res.Category,
			Count:    res.Count,
		})
	}

	return distribution
}

// getDemographicRiskDistribution retrieves risk distribution cross-referenced with a demographic field
func (r *reportRepository) getDemographicRiskDistribution(companyID uuid.UUID, demographicField string, period *int) []domain.DemographicRiskDistribution {
	var results []struct {
		Category  string           `gorm:"column:category"`
		RiskLevel domain.RiskLevel `gorm:"column:risk_level"`
		Count     int64            `gorm:"column:count"`
	}

	demographicQuery := fmt.Sprintf("staff.demographics->>'%s'", demographicField)
	query := r.db.Model(&domain.Assessment{}).
		Select(demographicQuery+" as category, assessments.risk_level, COUNT(*) as count").
		Joins("INNER JOIN staff ON assessments.staff_id = staff.id").
		Where("assessments.company_id = ? AND assessments.status = ? AND assessments.risk_level IS NOT NULL",
			companyID, domain.AssessmentStatusCompleted).
		Where(demographicQuery + " IS NOT NULL AND " + demographicQuery + " != ''")

	if period != nil {
		query = query.Where("assessments.period = ?", *period)
	}

	query.Group(demographicQuery + ", assessments.risk_level").
		Order("category, risk_level").
		Scan(&results)

	distribution := make([]domain.DemographicRiskDistribution, 0, len(results))
	for _, res := range results {
		distribution = append(distribution, domain.DemographicRiskDistribution{
			Category:  res.Category,
			RiskLevel: res.RiskLevel,
			Count:     res.Count,
		})
	}

	return distribution
}
