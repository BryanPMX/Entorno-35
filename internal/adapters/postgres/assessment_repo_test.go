package postgres

import (
	"testing"

	"github.com/entorno35/backend/internal/core/scoring"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
// Note: This test suite currently uses SQLite instead of PostgreSQL for simplicity.
// JSONB fields are handled as TEXT in SQLite for testing purposes.
// In production, PostgreSQL with proper JSONB support is used.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Create tables manually to handle JSONB fields in SQLite
	// In production PostgreSQL, these would be proper JSONB columns
	err = db.Exec(`
		CREATE TABLE companies (
			id TEXT PRIMARY KEY,
			rfc TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			address TEXT,
			subscription_status TEXT NOT NULL DEFAULT 'inactive',
			employee_count INTEGER NOT NULL DEFAULT 0,
			subscription_start_date DATE,
			subscription_end_date DATE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		);

		CREATE TABLE staff (
			id TEXT PRIMARY KEY,
			company_id TEXT NOT NULL,
			curp TEXT NOT NULL,
			full_name TEXT NOT NULL,
			email TEXT,
			password_hash TEXT, -- Bcrypt hashed password for staff authentication
			demographics TEXT, -- JSON stored as TEXT in SQLite
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY (company_id) REFERENCES companies(id)
		);

		CREATE TABLE assessments (
			id TEXT PRIMARY KEY,
			staff_id TEXT NOT NULL,
			company_id TEXT NOT NULL,
			period INTEGER NOT NULL,
			guide_type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			total_score REAL,
			risk_level TEXT,
			requires_medical_attention BOOLEAN DEFAULT FALSE,
			completed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY (staff_id) REFERENCES staff(id),
			FOREIGN KEY (company_id) REFERENCES companies(id)
		);

		CREATE TABLE responses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			assessment_id TEXT NOT NULL,
			question_id INTEGER NOT NULL,
			selected_value INTEGER NOT NULL CHECK (selected_value >= 0 AND selected_value <= 4),
			calculated_score INTEGER NOT NULL,
			answered_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (assessment_id) REFERENCES assessments(id)
		);

		CREATE TABLE assessment_links (
			id TEXT PRIMARY KEY,
			token TEXT UNIQUE NOT NULL,
			staff_id TEXT NOT NULL,
			assessment_id TEXT,
			expires_at DATETIME NOT NULL,
			accessed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (staff_id) REFERENCES staff(id),
			FOREIGN KEY (assessment_id) REFERENCES assessments(id)
		);
	`).Error
	require.NoError(t, err)

	return db
}

func TestAssessmentRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAssessmentRepository(db)

	// Create test company and staff
	companyID := uuid.New()
	staffID := uuid.New()
	company := domain.Company{
		ID:                companyID,
		RFC:               "TEST123456789",
		Name:              "Test Company",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:     10,
	}
	db.Create(&company)

	staff := domain.Staff{
		ID:        staffID,
		CompanyID: companyID,
		CURP:      "TESTCURP12345678901",
		FullName:  "Test Staff",
	}
	db.Create(&staff)

	// Create test assessment
	assessmentID := uuid.New()
	assessment := domain.Assessment{
		ID:        assessmentID,
		StaffID:   staffID,
		CompanyID: companyID,
		Period:    2025,
		GuideType: domain.GuideTypeII,
		Status:    domain.AssessmentStatusPending,
	}
	db.Create(&assessment)

	// Test GetByID
	retrieved, err := repo.GetByID(assessmentID)
	require.NoError(t, err)
	assert.Equal(t, assessmentID, retrieved.ID)
	assert.Equal(t, companyID, retrieved.CompanyID)
	assert.Equal(t, staffID, retrieved.StaffID)
}

func TestAssessmentRepository_GetByIDAndCompany(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAssessmentRepository(db)

	// Create test company and staff
	companyID := uuid.New()
	otherCompanyID := uuid.New()
	staffID := uuid.New()

	company := domain.Company{
		ID:                companyID,
		RFC:               "TEST123456789",
		Name:              "Test Company",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:     10,
	}
	db.Create(&company)

	otherCompany := domain.Company{
		ID:                otherCompanyID,
		RFC:               "OTHER123456789",
		Name:              "Other Company",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:     10,
	}
	db.Create(&otherCompany)

	staff := domain.Staff{
		ID:        staffID,
		CompanyID: companyID,
		CURP:      "TESTCURP12345678901",
		FullName:  "Test Staff",
	}
	db.Create(&staff)

	// Create test assessment
	assessmentID := uuid.New()
	assessment := domain.Assessment{
		ID:        assessmentID,
		StaffID:   staffID,
		CompanyID: companyID,
		Period:    2025,
		GuideType: domain.GuideTypeII,
		Status:    domain.AssessmentStatusPending,
	}
	db.Create(&assessment)

	// Test GetByIDAndCompany with correct company
	retrieved, err := repo.GetByIDAndCompany(assessmentID, companyID)
	require.NoError(t, err)
	assert.Equal(t, assessmentID, retrieved.ID)

	// Test GetByIDAndCompany with wrong company (should not found)
	_, err = repo.GetByIDAndCompany(assessmentID, otherCompanyID)
	assert.Error(t, err)
	assert.Equal(t, ErrAssessmentNotFound, err)
}

func TestAssessmentRepository_UpdateResult(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAssessmentRepository(db)

	// Create test data
	companyID := uuid.New()
	staffID := uuid.New()
	assessmentID := uuid.New()

	company := domain.Company{
		ID:                companyID,
		RFC:               "TEST123456789",
		Name:              "Test Company",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:     10,
	}
	db.Create(&company)

	staff := domain.Staff{
		ID:        staffID,
		CompanyID: companyID,
		CURP:      "TESTCURP12345678901",
		FullName:  "Test Staff",
	}
	db.Create(&staff)

	assessment := domain.Assessment{
		ID:        assessmentID,
		StaffID:   staffID,
		CompanyID: companyID,
		Period:    2025,
		GuideType: domain.GuideTypeII,
		Status:    domain.AssessmentStatusPending,
	}
	db.Create(&assessment)

	// Create result
	result := &scoring.AssessmentResult{
		TotalScore:               85.5,
		RiskLevel:                domain.RiskLevelAlto,
		CategoryScores:           make(map[string]float64),
		DomainScores:             make(map[string]float64),
		RequiresMedicalAttention: false,
	}

	// Update result
	err := repo.UpdateResult(&assessment, result)
	require.NoError(t, err)

	// Verify update
	var updated domain.Assessment
	db.First(&updated, assessmentID)
	assert.NotNil(t, updated.TotalScore)
	assert.Equal(t, 85.5, *updated.TotalScore)
	assert.NotNil(t, updated.RiskLevel)
	assert.Equal(t, domain.RiskLevelAlto, *updated.RiskLevel)
	assert.False(t, updated.RequiresMedicalAttention)
	assert.Equal(t, domain.AssessmentStatusCompleted, updated.Status)
	assert.NotNil(t, updated.CompletedAt)
}

