//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAuditData creates a complete test scenario for "The Compliance Audit"
// Returns: companyID, staffAID, staffBID, assessmentAID, assessmentBID
func setupAuditData(t *testing.T) (uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID) {
	// Create a unique test company
	companyID := uuid.New()
	company := &domain.Company{
		ID:                 companyID,
		RFC:                fmt.Sprintf("TEST%08d", time.Now().Unix()%100000000),
		Name:               "Test Compliance Company",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:      100, // Guide III
	}
	err := TestDB.Create(company).Error
	require.NoError(t, err, "Failed to create test company")

	// Create Staff A: IT Dept, High Risk
	staffAID := uuid.New()
	staffA := &domain.Staff{
		ID:        staffAID,
		CompanyID: companyID,
		CURP:      "TESTA123456HIJK01", // 18 characters (valid CURP length)
		FullName:  "Juan Pérez - IT Department",
		Email:     "juan.perez@test.com",
		Demographics: domain.DemographicsJSONB{
			Department: "IT",
			ShiftType:  "Diurno",
			Gender:     "Masculino",
		},
	}
	err = TestDB.Create(staffA).Error
	require.NoError(t, err, "Failed to create Staff A")

	// Create Staff B: HR Dept, Low Risk
	staffBID := uuid.New()
	staffB := &domain.Staff{
		ID:        staffBID,
		CompanyID: companyID,
		CURP:      "TESTB123456HIJK02", // 18 characters (valid CURP length)
		FullName:  "María González - HR Department",
		Email:     "maria.gonzalez@test.com",
		Demographics: domain.DemographicsJSONB{
			Department: "HR",
			ShiftType:  "Diurno",
			Gender:     "Femenino",
		},
	}
	err = TestDB.Create(staffB).Error
	require.NoError(t, err, "Failed to create Staff B")

	// Ensure questions exist (seed minimal Guide III questions if needed)
	seedMinimalQuestions(t)

	// Create Assessment A for Staff A (Guide III - High Risk)
	assessmentAID := uuid.New()
	assessmentA := &domain.Assessment{
		ID:        assessmentAID,
		StaffID:   staffAID,
		CompanyID: companyID,
		Period:    2025,
		GuideType: domain.GuideTypeIII,
		Status:    domain.AssessmentStatusPending,
	}
	err = TestDB.Create(assessmentA).Error
	require.NoError(t, err, "Failed to create Assessment A")

	// Create Assessment B for Staff B (Guide III - Low Risk)
	assessmentBID := uuid.New()
	assessmentB := &domain.Assessment{
		ID:        assessmentBID,
		StaffID:   staffBID,
		CompanyID: companyID,
		Period:    2025,
		GuideType: domain.GuideTypeIII,
		Status:    domain.AssessmentStatusPending,
	}
	err = TestDB.Create(assessmentB).Error
	require.NoError(t, err, "Failed to create Assessment B")

	// Get Guide III questions for responses
	var questions []domain.Question
	err = TestDB.Where("guide_type = ?", domain.GuideTypeIII).
		Preload("Category").
		Preload("Domain").
		Preload("Dimension").
		Order("question_number").
		Find(&questions).Error
	require.NoError(t, err, "Failed to fetch Guide III questions")
	require.GreaterOrEqual(t, len(questions), 25, "Need at least 25 Guide III questions to generate high risk scores (Guide III threshold for 'alto' is 99)")

	// Create responses for Staff A (High Risk - select high values for negative polarity questions)
	responsesA := createHighRiskResponses(t, assessmentAID, questions)
	err = TestDB.CreateInBatches(responsesA, 100).Error
	require.NoError(t, err, "Failed to create responses for Staff A")

	// Create responses for Staff B (Low Risk - select low values)
	responsesB := createLowRiskResponses(t, assessmentBID, questions)
	err = TestDB.CreateInBatches(responsesB, 100).Error
	require.NoError(t, err, "Failed to create responses for Staff B")

	// Calculate scores using the scoring service
	assessmentRepo := postgres.NewAssessmentRepository(TestDB)
	scoringService := services.NewScoringService(assessmentRepo)

	// Score Assessment A (High Risk)
	err = scoringService.CalculateAssessmentWithCompany(assessmentAID, companyID)
	require.NoError(t, err, "Failed to calculate score for Assessment A")

	// Score Assessment B (Low Risk)
	err = scoringService.CalculateAssessmentWithCompany(assessmentBID, companyID)
	require.NoError(t, err, "Failed to calculate score for Assessment B")

	// Verify assessments are marked as completed
	var updatedA domain.Assessment
	TestDB.First(&updatedA, assessmentAID)
	require.Equal(t, domain.AssessmentStatusCompleted, updatedA.Status, "Assessment A should be completed")
	require.NotNil(t, updatedA.RiskLevel, "Assessment A should have risk level")
	require.NotNil(t, updatedA.TotalScore, "Assessment A should have total score")

	var updatedB domain.Assessment
	TestDB.First(&updatedB, assessmentBID)
	require.Equal(t, domain.AssessmentStatusCompleted, updatedB.Status, "Assessment B should be completed")
	require.NotNil(t, updatedB.RiskLevel, "Assessment B should have risk level")
	require.NotNil(t, updatedB.TotalScore, "Assessment B should have total score")

	return companyID, staffAID, staffBID, assessmentAID, assessmentBID
}

// seedMinimalQuestions ensures we have at least some Guide III questions for testing
func seedMinimalQuestions(t *testing.T) {
	// Check if we have enough questions (need at least 25 for high risk: 25*4=100 points)
	var count int64
	TestDB.Model(&domain.Question{}).Where("guide_type = ?", domain.GuideTypeIII).Count(&count)
	if count >= 25 {
		return // We have enough questions
	}

	// Create minimal category, domain, and questions for Guide III
	category := domain.Category{Name: "Factores propios de la actividad"}
	TestDB.FirstOrCreate(&category, category)

	domainObj := domain.Domain{CategoryID: category.ID, Name: "Carga de trabajo"}
	TestDB.FirstOrCreate(&domainObj, domainObj)

	// Create enough questions to generate high risk scores
	// Guide III threshold for "alto" is 99, so we need at least 25 questions with score 4 each (25*4=100)
	for i := 1; i <= 30; i++ {
		polarity := domain.QuestionPolarityNegative
		question := domain.Question{
			QuestionNumber: i,
			GuideType:      domain.GuideTypeIII,
			Type:           domain.QuestionTypeLikert,
			Text:           fmt.Sprintf("Test Question %d", i),
			Polarity:       &polarity,
			CategoryID:     &category.ID,
			DomainID:       &domainObj.ID,
			OrderIndex:     i - 1,
		}
		TestDB.FirstOrCreate(&question, domain.Question{
			QuestionNumber: i,
			GuideType:      domain.GuideTypeIII,
		})
	}
}

// createHighRiskResponses creates responses that will result in high risk scores
// For high risk, we need high calculated scores:
// - Positive polarity: select 4 -> calculated score 4
// - Negative polarity: select 0 -> calculated score 4 (after inversion: 4-0=4)
func createHighRiskResponses(t *testing.T, assessmentID uuid.UUID, questions []domain.Question) []domain.Response {
	responses := make([]domain.Response, 0, len(questions))
	
	for _, question := range questions {
		var selectedValue int
		var calculatedScore int

		if question.Polarity == nil {
			// Binary question or no polarity - use high value
			selectedValue = 4
			calculatedScore = 4
		} else if *question.Polarity == domain.QuestionPolarityPositive {
			// Positive: high selected value = high calculated score
			selectedValue = 4
			calculatedScore = 4
		} else {
			// Negative: low selected value (0) = high calculated score (4) after inversion
			selectedValue = 0
			calculatedScore = 4 // After inversion: 4 - 0 = 4
		}

		response := domain.Response{
			AssessmentID:    assessmentID,
			QuestionID:      question.ID,
			SelectedValue:   selectedValue,
			CalculatedScore: calculatedScore,
			AnsweredAt:      timePtr(time.Now()),
		}
		responses = append(responses, response)
	}

	return responses
}

// createLowRiskResponses creates responses that will result in low risk scores
// For low risk, we need low calculated scores:
// - Positive polarity: select 0 -> calculated score 0
// - Negative polarity: select 4 -> calculated score 0 (after inversion: 4-4=0)
func createLowRiskResponses(t *testing.T, assessmentID uuid.UUID, questions []domain.Question) []domain.Response {
	responses := make([]domain.Response, 0, len(questions))
	
	for _, question := range questions {
		var selectedValue int
		var calculatedScore int

		if question.Polarity == nil {
			// Binary question or no polarity - use low value
			selectedValue = 0
			calculatedScore = 0
		} else if *question.Polarity == domain.QuestionPolarityPositive {
			// Positive: low selected value = low calculated score
			selectedValue = 0
			calculatedScore = 0
		} else {
			// Negative: high selected value (4) = low calculated score (0) after inversion
			selectedValue = 4
			calculatedScore = 0 // After inversion: 4 - 4 = 0
		}

		response := domain.Response{
			AssessmentID:    assessmentID,
			QuestionID:      question.ID,
			SelectedValue:   selectedValue,
			CalculatedScore: calculatedScore,
			AnsweredAt:      timePtr(time.Now()),
		}
		responses = append(responses, response)
	}

	return responses
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// TestReportE2E_TheComplianceAudit is the main E2E test for the reporting engine
func TestReportE2E_TheComplianceAudit(t *testing.T) {
	// Setup: Create test data
	companyID, _, _, assessmentAID, _ := setupAuditData(t)
	
	// Teardown: Always clean up, even if test fails
	defer func() {
		CleanupTestCompany(t, companyID.String())
	}()

	// Initialize report service
	reportRepo := postgres.NewReportRepository(TestDB)
	reportService := services.NewReportService(reportRepo)

	// ============================================================================
	// Test A: Individual Report for Staff A (High Risk)
	// ============================================================================
	t.Run("Individual_Report_High_Risk", func(t *testing.T) {
		report, err := reportService.GenerateIndividualReport(assessmentAID, companyID)
		require.NoError(t, err, "Should generate individual report successfully")
		require.NotNil(t, report, "Report should not be nil")

		// Assert basic fields
		assert.Equal(t, assessmentAID, report.AssessmentID, "Assessment ID should match")
		assert.Equal(t, 2025, report.Period, "Period should be 2025")
		assert.Equal(t, domain.GuideTypeIII, report.GuideType, "Guide type should be III")
		assert.Equal(t, "Juan Pérez - IT Department", report.StaffName, "Staff name should match")
		assert.Equal(t, "IT", report.Department, "Department should be IT")
		assert.Equal(t, "Diurno", report.Shift, "Shift should be Diurno")

		// Assert scoring fields
		assert.Greater(t, report.TotalScore, float64(0), "Total score should be greater than 0")
		assert.NotNil(t, report.RiskLevel, "Risk level should be set")
		
		// For high risk test, we expect at least "alto" risk level
		// Guide III thresholds: nulo < 50, bajo 50-74, medio 75-98, alto 99-139, muy_alto >= 140
		// With 30 questions * 4 points = 120 points, we should get "alto" risk
		assert.GreaterOrEqual(t, report.TotalScore, float64(99), "Total score should be >= 99 for 'alto' risk level in Guide III")
		assert.True(t, report.RiskLevel == domain.RiskLevelAlto || report.RiskLevel == domain.RiskLevelMuyAlto, 
			"Risk level should be 'alto' or 'muy_alto' for high risk test (got: %s, score: %.2f)", report.RiskLevel, report.TotalScore)
		
		// Critical NOM-035 requirement: Recommendations must be present for high risk
		assert.NotEmpty(t, report.Recommendations, "Recommendations array should NOT be empty for high risk (Critical NOM-035 requirement)")
		assert.Greater(t, len(report.Recommendations), 0, "Should have at least one recommendation")

		// Assert domain scores exist
		assert.NotNil(t, report.DomainScores, "Domain scores should not be nil")
		if len(report.DomainScores) > 0 {
			// Verify domain scores are populated
			hasScores := false
			for _, score := range report.DomainScores {
				if score > 0 {
					hasScores = true
					break
				}
			}
			assert.True(t, hasScores, "At least one domain should have a score > 0")
		}

		t.Logf("Individual Report - Staff A:")
		t.Logf("  Total Score: %.2f", report.TotalScore)
		t.Logf("  Risk Level: %s", report.RiskLevel)
		t.Logf("  Recommendations: %d", len(report.Recommendations))
	})

	// ============================================================================
	// Test B: General Report (The Aggregation)
	// ============================================================================
	t.Run("General_Report_Aggregation", func(t *testing.T) {
		period := 2025
		report, err := reportService.GenerateGeneralReport(companyID, &period)
		require.NoError(t, err, "Should generate general report successfully")
		require.NotNil(t, report, "Report should not be nil")

		// Assert basic fields
		assert.Equal(t, companyID, report.CompanyID, "Company ID should match")
		assert.Equal(t, "Test Compliance Company", report.CompanyName, "Company name should match")
		assert.Equal(t, 2025, *report.Period, "Period should be 2025")

		// Assert participation metrics
		assert.Equal(t, int64(2), report.TotalStaff, "Should have 2 staff members")
		assert.Equal(t, int64(2), report.CompletedAssessments, "Should have 2 completed assessments")
		assert.Equal(t, 100.0, report.ParticipationRate, "Participation rate should be 100% (2/2)")

		// Assert risk distribution
		assert.NotNil(t, report.RiskDistribution, "Risk distribution should not be nil")
		assert.Greater(t, len(report.RiskDistribution), 0, "Should have at least one risk level")

		// Verify risk distribution counts match our seeded data
		totalCount := int64(0)
		hasHighRisk := false
		hasLowRisk := false
		for _, dist := range report.RiskDistribution {
			totalCount += dist.Count
			if dist.RiskLevel == domain.RiskLevelAlto || dist.RiskLevel == domain.RiskLevelMuyAlto {
				hasHighRisk = true
			}
			if dist.RiskLevel == domain.RiskLevelBajo || dist.RiskLevel == domain.RiskLevelNulo {
				hasLowRisk = true
			}
		}
		assert.Equal(t, int64(2), totalCount, "Total assessments in risk distribution should be 2")
		assert.True(t, hasHighRisk || hasLowRisk, "Should have at least one risk level represented")

		// Assert department heatmap
		assert.NotNil(t, report.DepartmentHeatmap, "Department heatmap should not be nil")
		
		// Verify heatmap contains both departments
		departments := make(map[string]bool)
		for _, heatmap := range report.DepartmentHeatmap {
			departments[heatmap.Department] = true
		}
		
		assert.True(t, departments["IT"], "Heatmap should contain IT department")
		assert.True(t, departments["HR"], "Heatmap should contain HR department")

		// Verify heatmap has risk levels for each department
		itRiskLevels := make(map[domain.RiskLevel]bool)
		hrRiskLevels := make(map[domain.RiskLevel]bool)
		for _, heatmap := range report.DepartmentHeatmap {
			if heatmap.Department == "IT" {
				itRiskLevels[heatmap.RiskLevel] = true
			}
			if heatmap.Department == "HR" {
				hrRiskLevels[heatmap.RiskLevel] = true
			}
		}
		assert.True(t, len(itRiskLevels) > 0, "IT department should have at least one risk level")
		assert.True(t, len(hrRiskLevels) > 0, "HR department should have at least one risk level")

		t.Logf("General Report:")
		t.Logf("  Total Staff: %d", report.TotalStaff)
		t.Logf("  Completed Assessments: %d", report.CompletedAssessments)
		t.Logf("  Participation Rate: %.2f%%", report.ParticipationRate)
		t.Logf("  Risk Distribution Entries: %d", len(report.RiskDistribution))
		t.Logf("  Department Heatmap Entries: %d", len(report.DepartmentHeatmap))
	})

	// ============================================================================
	// Test C: HTTP Endpoint Test (if API is running)
	// ============================================================================
	t.Run("HTTP_Endpoints", func(t *testing.T) {
		// Skip if API is not running
		if BaseURL == "" || BaseURL == "http://localhost:8080" {
			// Try to ping the health endpoint
			resp, err := http.Get(BaseURL + "/health")
			if err != nil {
				t.Skip("API server not running, skipping HTTP endpoint tests")
				return
			}
			resp.Body.Close()
		}

		// Get auth token
		token, err := CreateAuthToken(companyID.String())
		require.NoError(t, err, "Should create auth token")

		// Test Individual Report Endpoint
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reports/individual/%s", BaseURL, assessmentAID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Individual report endpoint should return 200")

		var report domain.IndividualReportDTO
		err = json.NewDecoder(resp.Body).Decode(&report)
		require.NoError(t, err, "Should decode individual report JSON")
		assert.Equal(t, assessmentAID, report.AssessmentID, "Assessment ID should match")

		// Test General Report Endpoint
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reports/general?period=2025", BaseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "General report endpoint should return 200")

		var generalReport domain.GeneralReportDTO
		err = json.NewDecoder(resp.Body).Decode(&generalReport)
		require.NoError(t, err, "Should decode general report JSON")
		assert.Equal(t, companyID, generalReport.CompanyID, "Company ID should match")
	})
}

