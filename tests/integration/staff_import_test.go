//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testCompanyAID uuid.UUID
	testCompanyBID uuid.UUID
)

// initTestCompanies creates test companies for staff import tests
func initTestCompanies() {
	testCompanyAID = uuid.New()
	testCompanyBID = uuid.New()

	testCompanyA := &domain.Company{
		ID:                testCompanyAID,
		RFC:               "TESTCOMPA01",
		Name:              "Test Company A",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:     100,
	}
	TestDB.Create(testCompanyA)

	testCompanyB := &domain.Company{
		ID:                testCompanyBID,
		RFC:               "TESTCOMPB01",
		Name:              "Test Company B",
		SubscriptionStatus: domain.SubscriptionStatusActive,
		EmployeeCount:     50,
	}
	TestDB.Create(testCompanyB)
}

// cleanupTestCompanies removes test companies
func cleanupTestCompanies() {
	TestDB.Exec("DELETE FROM responses WHERE assessment_id IN (SELECT id FROM assessments WHERE company_id IN (?, ?))", testCompanyAID, testCompanyBID)
	TestDB.Exec("DELETE FROM assessment_links WHERE staff_id IN (SELECT id FROM staff WHERE company_id IN (?, ?))", testCompanyAID, testCompanyBID)
	TestDB.Exec("DELETE FROM assessments WHERE company_id IN (?, ?)", testCompanyAID, testCompanyBID)
	TestDB.Exec("DELETE FROM staff WHERE company_id IN (?, ?)", testCompanyAID, testCompanyBID)
	TestDB.Exec("DELETE FROM companies WHERE id IN (?, ?)", testCompanyAID, testCompanyBID)
}

// createAuthToken generates a JWT token for a company (uses shared JWTService)
func createAuthToken(companyID uuid.UUID) (string, error) {
	return CreateAuthToken(companyID.String())
}

// createMultipartRequest creates a multipart form request with CSV file
func createMultipartRequest(url string, csvContent string, token string) (*http.Request, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	fileWriter, err := writer.CreateFormFile("file", "staff.csv")
	if err != nil {
		return nil, err
	}
	_, err = fileWriter.Write([]byte(csvContent))
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

// generateValidCURP generates a valid 18-character CURP for testing
// Format: 4 letters + 6 digits + 2 letters + 6 alphanumeric = 18 characters total
// Pattern: ^[A-Z]{4}[0-9]{6}[A-Z]{2}[A-Z0-9]{6}$
func generateValidCURP(index int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	
	curp := ""
	// First 4: letters
	for i := 0; i < 4; i++ {
		curp += string(letters[(index+i)%26])
	}
	// Next 6: digits
	for i := 0; i < 6; i++ {
		curp += string(digits[(index+i)%10])
	}
	// Next 2: letters
	curp += string(letters[index%26])
	curp += string(letters[(index+1)%26])
	// Next 6: alphanumeric (mix of letters and digits)
	for i := 0; i < 6; i++ {
		if (index+i)%2 == 0 {
			curp += string(digits[(index+i)%10])
		} else {
			curp += string(letters[(index+i)%26])
		}
	}
	
	return curp
}

// TestStaffImport_HappyPath tests successful import of 2,500 valid rows (forces multiple batches)
func TestStaffImport_HappyPath(t *testing.T) {
	// Initialize test companies
	initTestCompanies()
	defer cleanupTestCompanies()
	
	// Clean up before test
	TestDB.Exec("DELETE FROM staff WHERE company_id = ?", testCompanyAID)

	// Create auth token
	token, err := createAuthToken(testCompanyAID)
	require.NoError(t, err)

	// Generate CSV with 2,500 valid rows (forces batch processing with batch size 1000)
	var csvBuilder strings.Builder
	csvBuilder.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	for i := 1; i <= 2500; i++ {
		curp := generateValidCURP(i)
		csvBuilder.WriteString(fmt.Sprintf("Test User %d,%s,user%d@example.com,Production,Operator,Day,Male\n", i, curp, i))
	}

	// Create request
	url := BaseURL + "/api/v1/staff/import"
	req, err := createMultipartRequest(url, csvBuilder.String(), token)
	require.NoError(t, err)

	// Send request (longer timeout for large file)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var result struct {
		TotalProcessed int      `json:"total_processed"`
		SuccessCount   int      `json:"success_count"`
		SkippedCount   int      `json:"skipped_count"`
		Errors         []string `json:"errors,omitempty"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Assert results
	assert.Equal(t, 2500, result.TotalProcessed, "Should process 2500 rows")
	assert.Equal(t, 2500, result.SuccessCount, "Should successfully import all 2500 rows")
	assert.Equal(t, 0, result.SkippedCount, "Should not skip any rows")
	assert.Empty(t, result.Errors, "Should have no errors")

	// Verify database - query DB to confirm persistence
	var count int64
	TestDB.Model(&domain.Staff{}).Where("company_id = ?", testCompanyAID).Count(&count)
	assert.Equal(t, int64(2500), count, "Database should contain exactly 2500 staff records")

	// Verify sample records have correct data
	var staff domain.Staff
	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Test User 1").First(&staff)
	assert.Equal(t, "Test User 1", staff.FullName)
	assert.Equal(t, generateValidCURP(1), staff.CURP)
	assert.Equal(t, "user1@example.com", staff.Email)

	// Verify a record from the last batch
	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Test User 2500").First(&staff)
	assert.Equal(t, "Test User 2500", staff.FullName)
	assert.Equal(t, generateValidCURP(2500), staff.CURP)
	assert.Equal(t, "user2500@example.com", staff.Email)
}

// TestStaffImport_SpanishCharacters tests UTF-8 encoding with Spanish characters
func TestStaffImport_SpanishCharacters(t *testing.T) {
	// Initialize test companies
	initTestCompanies()
	defer cleanupTestCompanies()
	
	// Clean up before test
	TestDB.Exec("DELETE FROM staff WHERE company_id = ?", testCompanyAID)

	// Create auth token
	token, err := createAuthToken(testCompanyAID)
	require.NoError(t, err)

	// CSV with Spanish characters (UTF-8 encoding test)
	var csvBuilder strings.Builder
	csvBuilder.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	
	names := []string{
		"José Nuñez López",
		"María González Pérez",
		"Carlos Ñoño Martínez",
		"Ana O'Brien Sánchez",
	}
	
	for i, name := range names {
		curp := generateValidCURP(30000 + i)
		csvBuilder.WriteString(fmt.Sprintf("%s,%s,test%d@example.com,Producción,Operador,Diurno,Masculino\n", 
			name, curp, i+1))
	}
	csvContent := csvBuilder.String()

	// Create request
	url := BaseURL + "/api/v1/staff/import"
	req, err := createMultipartRequest(url, csvContent, token)
	require.NoError(t, err)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var result struct {
		TotalProcessed int      `json:"total_processed"`
		SuccessCount   int      `json:"success_count"`
		SkippedCount   int      `json:"skipped_count"`
		Errors         []string `json:"errors,omitempty"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Assert results
	assert.Equal(t, 4, result.TotalProcessed)
	assert.Equal(t, 4, result.SuccessCount)
	assert.Empty(t, result.Errors)

	// Verify UTF-8 encoding in database (encoding check)
	var staff domain.Staff
	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "José Nuñez López").First(&staff)
	assert.Equal(t, "José Nuñez López", staff.FullName, "Database should store José Nuñez López correctly")

	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "María González Pérez").First(&staff)
	assert.Equal(t, "María González Pérez", staff.FullName, "Database should store María González Pérez correctly")

	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Carlos Ñoño Martínez").First(&staff)
	assert.Equal(t, "Carlos Ñoño Martínez", staff.FullName, "Database should store Carlos Ñoño Martínez correctly")

	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Ana O'Brien Sánchez").First(&staff)
	assert.Equal(t, "Ana O'Brien Sánchez", staff.FullName, "Database should store Ana O'Brien Sánchez correctly")
}

// TestStaffImport_DuplicateHandling tests idempotency (same CSV uploaded twice)
func TestStaffImport_DuplicateHandling(t *testing.T) {
	// Clean up before test
	TestDB.Exec("DELETE FROM staff WHERE company_id = ?", testCompanyAID)

	// Create auth token
	token, err := createAuthToken(testCompanyAID)
	require.NoError(t, err)

	// Generate CSV with 100 rows
	var csvBuilder strings.Builder
	csvBuilder.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	for i := 1; i <= 100; i++ {
		curp := generateValidCURP(1000 + i) // Use different range to avoid conflicts
		csvBuilder.WriteString(fmt.Sprintf("Duplicate Test User %d,%s,duptest%d@example.com,Area,Job,Shift,Gender\n", i, curp, i))
	}
	csvContent := csvBuilder.String()

	// First upload
	url := BaseURL + "/api/v1/staff/import"
	req1, err := createMultipartRequest(url, csvContent, token)
	require.NoError(t, err)

	client := &http.Client{Timeout: 30 * time.Second}
	resp1, err := client.Do(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()

	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	var result1 struct {
		TotalProcessed int      `json:"total_processed"`
		SuccessCount   int      `json:"success_count"`
		SkippedCount   int      `json:"skipped_count"`
		Errors         []string `json:"errors,omitempty"`
	}
	err = json.NewDecoder(resp1.Body).Decode(&result1)
	require.NoError(t, err)

	assert.Equal(t, 100, result1.TotalProcessed)
	assert.Equal(t, 100, result1.SuccessCount, "First upload should succeed for all rows")

	// Verify database count
	var countAfterFirst int64
	TestDB.Model(&domain.Staff{}).Where("company_id = ?", testCompanyAID).Count(&countAfterFirst)
	assert.Equal(t, int64(100), countAfterFirst)

	// Second upload (same CSV)
	req2, err := createMultipartRequest(url, csvContent, token)
	require.NoError(t, err)

	resp2, err := client.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var result2 struct {
		TotalProcessed int      `json:"total_processed"`
		SuccessCount   int      `json:"success_count"`
		SkippedCount   int      `json:"skipped_count"`
		Errors         []string `json:"errors,omitempty"`
	}
	err = json.NewDecoder(resp2.Body).Decode(&result2)
	require.NoError(t, err)

	// Second upload should have 0 success (all duplicates)
	assert.Equal(t, 100, result2.TotalProcessed)
	assert.Equal(t, 0, result2.SuccessCount, "Second upload should have 0 success (all duplicates)")
	assert.Equal(t, 100, result2.SkippedCount, "All rows should be skipped as duplicates")

	// Verify database count hasn't doubled
	var countAfterSecond int64
	TestDB.Model(&domain.Staff{}).Where("company_id = ?", testCompanyAID).Count(&countAfterSecond)
	assert.Equal(t, int64(100), countAfterSecond, "Database count should not double")
}

// TestStaffImport_MixedValidInvalid tests partial success with error reporting
func TestStaffImport_MixedValidInvalid(t *testing.T) {
	// Initialize test companies
	initTestCompanies()
	defer cleanupTestCompanies()
	// Clean up before test
	TestDB.Exec("DELETE FROM staff WHERE company_id = ?", testCompanyAID)

	// Create auth token
	token, err := createAuthToken(testCompanyAID)
	require.NoError(t, err)

	// CSV with mixed valid and invalid rows (10 valid, 1 invalid)
	// Use generateValidCURP to ensure valid format
	var csvBuilder strings.Builder
	csvBuilder.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	
	// Add 10 valid rows
	for i := 1; i <= 10; i++ {
		curp := generateValidCURP(5000 + i) // Use high index to avoid conflicts
		csvBuilder.WriteString(fmt.Sprintf("Valid User %d,%s,valid%d@example.com,Area,Job,Shift,Gender\n", i, curp, i))
	}
	
	// Add 1 invalid row (bad CURP format)
	csvBuilder.WriteString("Invalid CURP,SHORT,invalid@example.com,Area,Job,Shift,Gender\n")
	
	csvContent := csvBuilder.String()

	// Create request
	url := BaseURL + "/api/v1/staff/import"
	req, err := createMultipartRequest(url, csvContent, token)
	require.NoError(t, err)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var result struct {
		TotalProcessed int      `json:"total_processed"`
		SuccessCount   int      `json:"success_count"`
		SkippedCount   int      `json:"skipped_count"`
		Errors         []string `json:"errors,omitempty"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Assert results
	assert.Equal(t, 11, result.TotalProcessed, "Should process 11 rows (10 valid + 1 invalid)")
	assert.Equal(t, 10, result.SuccessCount, "Should have 10 valid rows successfully imported")
	assert.Equal(t, 1, result.SkippedCount, "Should skip 1 invalid row")
	assert.Len(t, result.Errors, 1, "Should have 1 error message for invalid CURP")
	assert.Contains(t, strings.Join(result.Errors, " "), "CURP", "Error should mention CURP")

	// Verify only valid records were inserted
	var count int64
	TestDB.Model(&domain.Staff{}).Where("company_id = ?", testCompanyAID).Count(&count)
	assert.Equal(t, int64(10), count, "Database should contain exactly 10 valid staff records")

	// Verify valid records exist
	var staff domain.Staff
	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Valid User 1").First(&staff)
	assert.Equal(t, "Valid User 1", staff.FullName)
	assert.Equal(t, generateValidCURP(5001), staff.CURP)

	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Valid User 10").First(&staff)
	assert.Equal(t, "Valid User 10", staff.FullName)
	assert.Equal(t, generateValidCURP(5010), staff.CURP)
}

// TestStaffImport_TenantIsolation tests multi-tenant isolation
func TestStaffImport_TenantIsolation(t *testing.T) {
	// Initialize test companies
	initTestCompanies()
	defer cleanupTestCompanies()
	// Clean up before test
	TestDB.Exec("DELETE FROM staff WHERE company_id IN (?, ?)", testCompanyAID, testCompanyBID)

	// Create auth tokens for both companies
	tokenA, err := createAuthToken(testCompanyAID)
	require.NoError(t, err)

	tokenB, err := createAuthToken(testCompanyBID)
	require.NoError(t, err)

	// Company A imports staff
	var csvBuilderA strings.Builder
	csvBuilderA.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	for i := 1; i <= 2; i++ {
		curp := generateValidCURP(10000 + i)
		csvBuilderA.WriteString(fmt.Sprintf("Company A User %d,%s,companya%d@example.com,Area,Job,Shift,Gender\n", i, curp, i))
	}
	csvA := csvBuilderA.String()

	url := BaseURL + "/api/v1/staff/import"
	reqA, err := createMultipartRequest(url, csvA, tokenA)
	require.NoError(t, err)

	client := &http.Client{Timeout: 10 * time.Second}
	respA, err := client.Do(reqA)
	require.NoError(t, err)
	defer respA.Body.Close()

	assert.Equal(t, http.StatusOK, respA.StatusCode)

	// Company B imports staff
	var csvBuilderB strings.Builder
	csvBuilderB.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	for i := 1; i <= 2; i++ {
		curp := generateValidCURP(20000 + i)
		csvBuilderB.WriteString(fmt.Sprintf("Company B User %d,%s,companyb%d@example.com,Area,Job,Shift,Gender\n", i, curp, i))
	}
	csvB := csvBuilderB.String()

	reqB, err := createMultipartRequest(url, csvB, tokenB)
	require.NoError(t, err)

	respB, err := client.Do(reqB)
	require.NoError(t, err)
	defer respB.Body.Close()

	assert.Equal(t, http.StatusOK, respB.StatusCode)

	// Verify Company A can only see its own staff
	var countA int64
	TestDB.Model(&domain.Staff{}).Where("company_id = ?", testCompanyAID).Count(&countA)
	assert.Equal(t, int64(2), countA)

	var staffA domain.Staff
	TestDB.Where("company_id = ? AND full_name = ?", testCompanyAID, "Company A User 1").First(&staffA)
	assert.Equal(t, "Company A User 1", staffA.FullName)
	assert.Equal(t, testCompanyAID, staffA.CompanyID)

	// Verify Company B can only see its own staff
	var countB int64
	TestDB.Model(&domain.Staff{}).Where("company_id = ?", testCompanyBID).Count(&countB)
	assert.Equal(t, int64(2), countB)

	var staffB domain.Staff
	TestDB.Where("company_id = ? AND full_name = ?", testCompanyBID, "Company B User 1").First(&staffB)
	assert.Equal(t, "Company B User 1", staffB.FullName)
	assert.Equal(t, testCompanyBID, staffB.CompanyID)

	// Verify cross-tenant access doesn't work
	var crossStaff domain.Staff
	result := TestDB.Where("company_id = ? AND full_name = ?", testCompanyBID, "Company A User 1").First(&crossStaff)
	assert.Error(t, result.Error)
	require.Error(t, result.Error, "Should not find cross-tenant staff")
}

// TestStaffImport_UnauthorizedAccess tests that unauthorized requests are rejected
func TestStaffImport_UnauthorizedAccess(t *testing.T) {
	// CSV content with valid CURP
	curp := generateValidCURP(99999)
	csvContent := fmt.Sprintf(`Name,CURP,Email,Area,Job,Shift,Gender
Test User,%s,test@example.com,Area,Job,Shift,Gender
`, curp)

	// Create request without token
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, err := writer.CreateFormFile("file", "staff.csv")
	require.NoError(t, err)
	_, err = fileWriter.Write([]byte(csvContent))
	require.NoError(t, err)
	writer.Close()

	url := BaseURL + "/api/v1/staff/import"
	req, err := http.NewRequest("POST", url, &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// No Authorization header

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should be unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

