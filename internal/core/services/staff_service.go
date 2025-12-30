package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

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

// StaffService handles staff-related business logic
type StaffService struct {
	staffRepo ports.StaffRepository
}

// NewStaffService creates a new staff service
func NewStaffService(staffRepo ports.StaffRepository) *StaffService {
	return &StaffService{
		staffRepo: staffRepo,
	}
}

// CreateStaffRequest represents the request to create a staff member
type CreateStaffRequest struct {
	CURP         string                    `json:"curp" binding:"required"`
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
	// Validate CURP format (must be 18 characters)
	if len(req.CURP) != 18 {
		return nil, fmt.Errorf("CURP must be exactly 18 characters, got %d", len(req.CURP))
	}

	staff := &domain.Staff{
		CompanyID:    companyID,
		CURP:         strings.ToUpper(strings.TrimSpace(req.CURP)), // Normalize to uppercase and trim
		FullName:     strings.TrimSpace(req.FullName),
		Email:        strings.TrimSpace(req.Email),
		Demographics: req.Demographics,
	}

	if err := s.staffRepo.Create(staff); err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return staff, nil
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
			CURP:         curp,
			FullName:     name,
			Email:        email,
			Demographics: demographics,
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

