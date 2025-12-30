//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/entorno35/backend/internal/core/jwt"
	"github.com/entorno35/backend/internal/database"
	"github.com/entorno35/backend/internal/domain"
	"gorm.io/gorm"
)

var (
	// Global test database connection
	TestDB *gorm.DB
	
	// Global JWT service for test authentication
	JWTService jwt.Service
	
	// Base URL for API tests (optional - can use direct DB access)
	BaseURL string
)

// TestMain sets up the test environment for all integration tests
func TestMain(m *testing.M) {
	// Get database URL from environment (required - no hardcoded defaults for security)
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("Error: DB_URL environment variable is required for integration tests")
		fmt.Println("Example: export DB_URL=postgres://user:password@localhost:5432/entorno35_test?sslmode=disable")
		fmt.Println("Never commit passwords to the repository - use environment variables only")
		os.Exit(1)
	}

	// Connect to database
	var err error
	TestDB, err = database.Connect(dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		fmt.Println("Make sure PostgreSQL is running and DB_URL is set correctly")
		fmt.Printf("Example: export DB_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable\n")
		os.Exit(1)
	}
	defer database.Close()

	// Run migrations (create tables)
	err = TestDB.AutoMigrate(
		&domain.Company{},
		&domain.Staff{},
		&domain.Category{},
		&domain.Domain{},
		&domain.Dimension{},
		&domain.Question{},
		&domain.Assessment{},
		&domain.Response{},
		&domain.AssessmentLink{},
	)
	if err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	// Initialize JWT service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret-key-for-integration-tests-only"
	}
	JWTService = jwt.NewService(jwtSecret)

	// Get base URL for API (optional)
	BaseURL = os.Getenv("API_BASE_URL")
	if BaseURL == "" {
		BaseURL = "http://localhost:8080"
	}

	// Run tests
	code := m.Run()

	// Global cleanup (if needed)
	// Individual tests should clean up their own data, but we can add global cleanup here if needed

	os.Exit(code)
}

// CreateAuthToken generates a JWT token for a company
func CreateAuthToken(companyID string) (string, error) {
	token, err := JWTService.GenerateToken(
		companyID,
		"", // No staff ID for company-level operations
		"test@example.com",
		"company",
		1*time.Hour,
	)
	return token, err
}

// CleanupTestCompany performs a hard delete of a test company and all related data
// This ensures cascading deletes work properly and test data is completely removed
func CleanupTestCompany(t *testing.T, companyID string) {
	// Use raw SQL for hard delete (bypass soft deletes)
	// Order matters due to foreign key constraints
	TestDB.Exec(`
		DELETE FROM responses 
		WHERE assessment_id IN (
			SELECT id FROM assessments WHERE company_id = ?
		)
	`, companyID)
	
	TestDB.Exec(`
		DELETE FROM assessment_links 
		WHERE staff_id IN (
			SELECT id FROM staff WHERE company_id = ?
		)
	`, companyID)
	
	TestDB.Exec(`
		DELETE FROM assessments 
		WHERE company_id = ?
	`, companyID)
	
	TestDB.Exec(`
		DELETE FROM staff 
		WHERE company_id = ?
	`, companyID)
	
	TestDB.Exec(`
		DELETE FROM companies 
		WHERE id = ?
	`, companyID)
}

