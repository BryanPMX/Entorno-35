package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	csvdetect "github.com/entorno35/backend/internal/core/csv"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// ImportResult represents the result of a CSV import operation
type ImportResult struct {
	TotalProcessed int      `json:"total_processed"`
	SuccessCount   int      `json:"success_count"`
	SkippedCount   int      `json:"skipped_count"`
	Errors         []string `json:"errors"`
}

// ImportPreview represents a preview of CSV import with validation
type ImportPreview struct {
	TotalRows     int            `json:"total_rows"`
	ValidRows     int            `json:"valid_rows"`
	InvalidRows   int            `json:"invalid_rows"`
	SampleRecords []ImportRecord `json:"sample_records"`
	Errors        []string       `json:"errors"`
}

// ImportRecord represents a single record from CSV import
type ImportRecord struct {
	RowNumber int               `json:"row_number"`
	Data      map[string]string `json:"data"`
	Errors    []string          `json:"errors"`
}

// StaffService handles staff-related business logic
type StaffService struct {
	staffRepo       ports.StaffRepository
	csvDetector     *csvdetect.DetectionService
}

// NewStaffService creates a new staff service
func NewStaffService(staffRepo ports.StaffRepository) *StaffService {
	return &StaffService{
		staffRepo:   staffRepo,
		csvDetector: csvdetect.NewDetectionService(),
	}
}

// CreateStaffRequest represents the request to create a staff member
type CreateStaffRequest struct {
	CURP         string                    `json:"curp,omitempty"` // Optional for staff without CURPs
	FullName     string                    `json:"full_name" binding:"required"`
	Email        string                    `json:"email,omitempty"`
	Demographics domain.DemographicsJSONB `json:"demographics,omitempty"`
}

// UpdateStaffRequest represents the request to update a staff member
type UpdateStaffRequest struct {
	FullName     string                    `json:"full_name,omitempty"`
	Email        string                    `json:"email,omitempty"`
	Demographics domain.DemographicsJSONB `json:"demographics,omitempty"`
}

// CreateStaff creates a new staff member
func (s *StaffService) CreateStaff(req CreateStaffRequest, companyID uuid.UUID) (*domain.Staff, error) {
	// Normalize input
	curp := strings.TrimSpace(req.CURP)
	fullName := strings.TrimSpace(req.FullName)
	email := strings.TrimSpace(req.Email)

	// Validate CURP if provided
	if curp != "" && len(curp) != 18 {
		return nil, fmt.Errorf("CURP must be exactly 18 characters, got %d", len(curp))
	}

	staff := &domain.Staff{
		CompanyID:    companyID,
		FullName:     fullName,
		Email:        email,
		Demographics: req.Demographics,
	}

	// Handle CURP vs Employee ID logic
	if curp != "" {
		// CURP provided - use it as primary identifier
		normalizedCURP := strings.ToUpper(curp)
		staff.CURP = &normalizedCURP
		staff.EmployeeID = nil // Explicitly set to nil for CURP users
	} else {
		// No CURP provided - auto-generate employee ID
		employeeID, err := s.generateEmployeeID(companyID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate employee ID: %w", err)
		}
		staff.EmployeeID = &employeeID // Set as pointer to generated string
		staff.CURP = nil // Explicitly set to nil
	}

	// Attempt to create staff, with retry logic for employee ID conflicts
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.staffRepo.Create(staff)
		if err == nil {
			// Success
			return staff, nil
		}

		// Check if this is a duplicate employee_id error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
		   strings.Contains(err.Error(), "uni_staff_employee_id") {
			// This is a duplicate employee ID error - regenerate and retry
			if curp == "" { // Only retry if we're using auto-generated employee IDs
				newEmployeeID, genErr := s.generateEmployeeID(companyID)
				if genErr != nil {
					return nil, fmt.Errorf("failed to regenerate employee ID: %w", genErr)
				}
				staff.EmployeeID = &newEmployeeID // Set as pointer
				continue // Retry with new employee ID
			}
		}

		// Not a duplicate employee ID error, or we've exhausted retries
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	// If we get here, we've exhausted all retries
	return nil, fmt.Errorf("failed to create staff after %d attempts due to employee ID conflicts", maxRetries)
}

// AnalyzeCSV analyzes a CSV file and returns column mapping suggestions
func (s *StaffService) AnalyzeCSV(reader io.Reader) (*csvdetect.DetectionResult, error) {
	return s.csvDetector.AnalyzeCSV(reader)
}

// PreviewCSVImport previews a CSV import with the given column mappings
func (s *StaffService) PreviewCSVImport(reader io.Reader, mappings []csvdetect.ColumnMapping, companyID uuid.UUID) (*ImportPreview, error) {
	// Reset reader to beginning
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	csvReader := csv.NewReader(reader)

	// Read headers first
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Create field mapping based on headers and mappings
	fieldMap := make(map[int]string) // column index -> field name
	for _, mapping := range mappings {
		// Find column index for this CSV header
		for colIndex, header := range headers {
			if strings.TrimSpace(strings.ToLower(header)) == strings.TrimSpace(strings.ToLower(mapping.CSVHeader)) {
				fieldMap[colIndex] = mapping.ExpectedField
				break
			}
		}
	}

	preview := &ImportPreview{
		TotalRows:     0,
		ValidRows:     0,
		InvalidRows:   0,
		SampleRecords: []ImportRecord{},
		Errors:        []string{},
	}

	// Process sample rows
	for i := 0; i < 5; i++ { // Preview first 5 rows
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		preview.TotalRows++
		record := s.parseCSVRow(row, fieldMap, preview.TotalRows)

		if len(record.Errors) == 0 {
			preview.ValidRows++
		} else {
			preview.InvalidRows++
		}

		preview.SampleRecords = append(preview.SampleRecords, record)
	}

	return preview, nil
}

// parseCSVRow parses a single CSV row according to the field mapping
func (s *StaffService) parseCSVRow(row []string, fieldMap map[int]string, rowNum int) ImportRecord {
	record := ImportRecord{
		RowNumber: rowNum,
		Data:      make(map[string]string),
		Errors:    []string{},
	}

	// Map columns to fields
	for colIndex, fieldName := range fieldMap {
		if colIndex < len(row) {
			record.Data[fieldName] = strings.TrimSpace(row[colIndex])
		}
	}

	// Validate the record
	s.validateImportRecord(&record)

	return record
}

// validateImportRecord validates a single import record
func (s *StaffService) validateImportRecord(record *ImportRecord) {
	// Check required fields
	if record.Data["name"] == "" {
		record.Errors = append(record.Errors, "name is required")
	}

	// Validate CURP if provided
	if curp := record.Data["curp"]; curp != "" && len(curp) != 18 {
		record.Errors = append(record.Errors, "CURP must be exactly 18 characters")
	}

	// Validate email if provided
	if email := record.Data["email"]; email != "" {
		if !strings.Contains(email, "@") {
			record.Errors = append(record.Errors, "invalid email format")
		}
	}
}

// generateEmployeeID generates a unique employee ID for staff without CURPs
// Format: CMP{COMPANY_ID}-{SEQUENCE} (e.g., CMP001-0001, CMP001-0002)
func (s *StaffService) generateEmployeeID(companyID uuid.UUID) (string, error) {
	// Get the short company ID (first 3 characters of UUID)
	companyShortID := companyID.String()[:3]

	// Find the highest existing employee ID for this company
	// This is a simplified approach - in production, you might want a dedicated sequence table
	prefix := fmt.Sprintf("CMP%s-", strings.ToUpper(companyShortID))

	// Query existing employee IDs with this prefix
	// For now, we'll use a simple approach - in production, consider using database sequences
	maxSequence := 0

	// Get all staff for this company and find the highest employee ID sequence
	staffList, _, err := s.staffRepo.ListByCompany(companyID, 1000, 0) // Get first 1000 records
	if err != nil {
		return "", fmt.Errorf("failed to query existing staff: %w", err)
	}

	// Find the highest sequence number for this company's prefix
	for _, staff := range staffList {
		if staff.EmployeeID != nil && strings.HasPrefix(*staff.EmployeeID, prefix) {
			// Extract sequence number from employee ID (format: CMPXXX-NNNN)
			parts := strings.Split(*staff.EmployeeID, "-")
			if len(parts) == 2 {
				var seq int
				if _, err := fmt.Sscanf(parts[1], "%d", &seq); err == nil {
					if seq > maxSequence {
						maxSequence = seq
					}
				}
			}
		}
	}

	// Generate next sequence number
	nextSequence := maxSequence + 1

	// Format as 4-digit zero-padded number
	return fmt.Sprintf("%s%04d", prefix, nextSequence), nil
}

// UpdateStaff updates an existing staff member
func (s *StaffService) UpdateStaff(id uuid.UUID, companyID uuid.UUID, req UpdateStaffRequest) (*domain.Staff, error) {
	// Fetch existing staff
	staff, err := s.staffRepo.GetByIDAndCompany(id, companyID)
	if err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}

	// Update fields if provided
	if req.FullName != "" {
		staff.FullName = req.FullName
	}
	if req.Email != "" {
		staff.Email = req.Email
	}
	// Update demographics if provided (merge with existing)
	if req.Demographics.Gender != "" {
		staff.Demographics.Gender = req.Demographics.Gender
	}
	if req.Demographics.AgeRange != "" {
		staff.Demographics.AgeRange = req.Demographics.AgeRange
	}
	if req.Demographics.MaritalStatus != "" {
		staff.Demographics.MaritalStatus = req.Demographics.MaritalStatus
	}
	if req.Demographics.EducationLevel != "" {
		staff.Demographics.EducationLevel = req.Demographics.EducationLevel
	}
	if req.Demographics.TimeInPosition != "" {
		staff.Demographics.TimeInPosition = req.Demographics.TimeInPosition
	}
	if req.Demographics.ShiftType != "" {
		staff.Demographics.ShiftType = req.Demographics.ShiftType
	}
	if req.Demographics.ShiftRotation != "" {
		staff.Demographics.ShiftRotation = req.Demographics.ShiftRotation
	}
	if req.Demographics.TotalWorkExperience != "" {
		staff.Demographics.TotalWorkExperience = req.Demographics.TotalWorkExperience
	}
	if req.Demographics.Department != "" {
		staff.Demographics.Department = req.Demographics.Department
	}
	if req.Demographics.Role != "" {
		staff.Demographics.Role = req.Demographics.Role
	}

	if err := s.staffRepo.Update(staff); err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return staff, nil
}

// DeleteStaff soft deletes a staff member
func (s *StaffService) DeleteStaff(id uuid.UUID, companyID uuid.UUID) error {
	return s.staffRepo.Delete(id, companyID)
}

// ListStaff retrieves staff members for a company with pagination
func (s *StaffService) ListStaff(companyID uuid.UUID, limit, offset int) ([]domain.Staff, int64, error) {
	return s.staffRepo.ListByCompany(companyID, limit, offset)
}

// GetStaff retrieves a staff member by ID
func (s *StaffService) GetStaff(id uuid.UUID, companyID uuid.UUID) (*domain.Staff, error) {
	return s.staffRepo.GetByIDAndCompany(id, companyID)
}

// ImportFromCSV imports staff members from a CSV file with batch processing
// CSV Format: Name, CURP, Email, Area, Job, Shift, Gender
// Processes in batches of BatchSize to avoid OOM errors with large files
func (s *StaffService) ImportFromCSV(r io.Reader, companyID uuid.UUID) (*ImportResult, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header format (case-insensitive)
	expectedHeaders := []string{"name", "curp", "email", "area", "job", "shift", "gender"}
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Check if all required headers are present
	for _, expected := range expectedHeaders {
		if _, exists := headerMap[expected]; !exists {
			return nil, fmt.Errorf("missing required header: %s", expected)
		}
	}

	result := &ImportResult{
		Errors: make([]string, 0),
	}

	// Batch processing: accumulate valid staff and flush in chunks
	validStaff := make([]*domain.Staff, 0, BatchSize)
	rowNum := 1 // Start at 1 (header is row 0)

	// Helper function to flush batch to database
	flushBatch := func(batch []*domain.Staff) error {
		if len(batch) == 0 {
			return nil
		}

		insertedCount, err := s.staffRepo.BulkCreate(batch)
		if err != nil {
			return fmt.Errorf("failed to bulk create staff: %w", err)
		}

		result.SuccessCount += insertedCount
		result.SkippedCount += len(batch) - insertedCount // Records skipped due to duplicates
		return nil
	}

	// Read data rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			// Flush remaining batch
			if err := flushBatch(validStaff); err != nil {
				return nil, err
			}
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: failed to read CSV row: %v", rowNum, err))
			result.TotalProcessed++
			rowNum++
			continue
		}

		result.TotalProcessed++
		rowNum++

		// Skip empty rows
		if len(row) == 0 || (len(row) == 1 && strings.TrimSpace(row[0]) == "") {
			continue
		}

		// Extract values from row
		nameIdx := headerMap["name"]
		curpIdx := headerMap["curp"]
		emailIdx := headerMap["email"]
		areaIdx := headerMap["area"]
		jobIdx := headerMap["job"]
		shiftIdx := headerMap["shift"]
		genderIdx := headerMap["gender"]

		// Validate row has enough columns
		if len(row) <= max(nameIdx, curpIdx, emailIdx, areaIdx, jobIdx, shiftIdx, genderIdx) {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: insufficient columns", rowNum-1))
			result.SkippedCount++
			continue
		}

		name := strings.TrimSpace(row[nameIdx])
		curp := strings.TrimSpace(row[curpIdx])
		email := strings.TrimSpace(row[emailIdx])
		area := strings.TrimSpace(row[areaIdx])
		job := strings.TrimSpace(row[jobIdx])
		shift := strings.TrimSpace(row[shiftIdx])
		gender := strings.TrimSpace(row[genderIdx])

		// Strict validation using validation functions
		if err := ValidateName(name); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %s", rowNum-1, err.Error()))
			result.SkippedCount++
			continue
		}

		if curp == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: CURP is required", rowNum-1))
			result.SkippedCount++
			continue
		}

		// Strict CURP format validation
		curp = strings.ToUpper(curp)
		if err := ValidateCURPFormat(curp); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %s", rowNum-1, err.Error()))
			result.SkippedCount++
			continue
		}

		// Validate email format (optional field)
		if err := ValidateEmailFormat(email); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %s", rowNum-1, err.Error()))
			result.SkippedCount++
			continue
		}

		// Map demographics (CSV columns: Area -> Department, Job -> Role, Shift -> ShiftType)
		demographics := domain.DemographicsJSONB{
			Department: area,
			Role:       job,
			ShiftType:  shift,
			Gender:     gender,
		}

		// Create staff record
		staff := &domain.Staff{
			CompanyID:    companyID,
			FullName:     name,
			Email:        email,
			Demographics: demographics,
		}

		// Handle CURP assignment (now a *string)
		if curp != "" {
			staff.CURP = &curp
		} else {
			staff.CURP = nil
		}

		validStaff = append(validStaff, staff)

		// Flush batch when it reaches BatchSize
		if len(validStaff) >= BatchSize {
			if err := flushBatch(validStaff); err != nil {
				return nil, err
			}
			// Reset batch
			validStaff = validStaff[:0] // Clear slice but keep capacity
		}
	}

	return result, nil
}

// max returns the maximum of the given integers
func max(nums ...int) int {
	maxNum := nums[0]
	for _, n := range nums[1:] {
		if n > maxNum {
			maxNum = n
		}
	}
	return maxNum
}

