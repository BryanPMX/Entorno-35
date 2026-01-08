package http

import (
	"net/http"
	"strings"

	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/domain"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StaffHandler handles staff-related HTTP requests
type StaffHandler struct {
	staffService *services.StaffService
}

// NewStaffHandler creates a new staff handler
func NewStaffHandler(staffService *services.StaffService) *StaffHandler {
	return &StaffHandler{
		staffService: staffService,
	}
}

// CreateStaffRequest represents the request body for creating a staff member
type CreateStaffRequest struct {
	CURP         string `json:"curp,omitempty"` // Optional for staff without CURPs
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email,omitempty"`
	Demographics struct {
		Gender              string `json:"gender,omitempty"`
		AgeRange            string `json:"age_range,omitempty"`
		MaritalStatus       string `json:"marital_status,omitempty"`
		EducationLevel      string `json:"education_level,omitempty"`
		TimeInPosition      string `json:"time_in_position,omitempty"`
		ShiftType           string `json:"shift_type,omitempty"`
		ShiftRotation       string `json:"shift_rotation,omitempty"`
		TotalWorkExperience string `json:"total_work_experience,omitempty"`
		Department          string `json:"department,omitempty"`
		Role                string `json:"role,omitempty"`
	} `json:"demographics,omitempty"`
}

// UpdateStaffRequest represents the request body for updating a staff member
type UpdateStaffRequest struct {
	FullName     string `json:"full_name,omitempty"`
	Email        string `json:"email,omitempty"`
	Demographics struct {
		Gender              string `json:"gender,omitempty"`
		AgeRange            string `json:"age_range,omitempty"`
		MaritalStatus       string `json:"marital_status,omitempty"`
		EducationLevel      string `json:"education_level,omitempty"`
		TimeInPosition      string `json:"time_in_position,omitempty"`
		ShiftType           string `json:"shift_type,omitempty"`
		ShiftRotation       string `json:"shift_rotation,omitempty"`
		TotalWorkExperience string `json:"total_work_experience,omitempty"`
		Department          string `json:"department,omitempty"`
		Role                string `json:"role,omitempty"`
	} `json:"demographics,omitempty"`
}

// ListStaffRequest represents query parameters for listing staff
type ListStaffRequest struct {
	Limit  int `form:"limit,default=50"`
	Offset int `form:"offset,default=0"`
}

// CreateStaff creates a new staff member
// POST /api/v1/staff
func (h *StaffHandler) CreateStaff(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	var req CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map request to service request
	demographics := domain.DemographicsJSONB{
		Gender:              req.Demographics.Gender,
		AgeRange:            req.Demographics.AgeRange,
		MaritalStatus:       req.Demographics.MaritalStatus,
		EducationLevel:      req.Demographics.EducationLevel,
		TimeInPosition:      req.Demographics.TimeInPosition,
		ShiftType:           req.Demographics.ShiftType,
		ShiftRotation:       req.Demographics.ShiftRotation,
		TotalWorkExperience: req.Demographics.TotalWorkExperience,
		Department:          req.Demographics.Department,
		Role:                req.Demographics.Role,
	}

	serviceReq := services.CreateStaffRequest{
		CURP:         req.CURP,
		FullName:     req.FullName,
		Email:        req.Email,
		Demographics: demographics,
	}

	staff, err := h.staffService.CreateStaff(serviceReq, companyID)
	if err != nil {
		// Check for client errors (validation, duplicates) vs server errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "CURP must be exactly 18 characters") ||
			strings.Contains(errMsg, "duplicate key value violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, staff)
}

// ListStaff retrieves staff members with pagination
// GET /api/v1/staff
func (h *StaffHandler) ListStaff(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	var req ListStaffRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate pagination parameters
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	staff, total, err := h.staffService.ListStaff(companyID, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   staff,
		"total":  total,
		"limit":  req.Limit,
		"offset": req.Offset,
	})
}

// GetStaff retrieves a staff member by ID
// GET /api/v1/staff/:id
func (h *StaffHandler) GetStaff(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	staffIDStr := c.Param("id")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}

	staff, err := h.staffService.GetStaff(staffID, companyID)
	if err != nil {
		if err.Error() == "staff not found: staff not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// UpdateStaff updates an existing staff member
// PUT /api/v1/staff/:id
func (h *StaffHandler) UpdateStaff(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	staffIDStr := c.Param("id")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}

	var req UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map request to service request
	demographics := domain.DemographicsJSONB{
		Gender:              req.Demographics.Gender,
		AgeRange:            req.Demographics.AgeRange,
		MaritalStatus:       req.Demographics.MaritalStatus,
		EducationLevel:      req.Demographics.EducationLevel,
		TimeInPosition:      req.Demographics.TimeInPosition,
		ShiftType:           req.Demographics.ShiftType,
		ShiftRotation:       req.Demographics.ShiftRotation,
		TotalWorkExperience: req.Demographics.TotalWorkExperience,
		Department:          req.Demographics.Department,
		Role:                req.Demographics.Role,
	}

	serviceReq := services.UpdateStaffRequest{
		FullName:     req.FullName,
		Email:        req.Email,
		Demographics: demographics,
	}

	staff, err := h.staffService.UpdateStaff(staffID, companyID, serviceReq)
	if err != nil {
		if err.Error() == "staff not found: staff not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// DeleteStaff soft deletes a staff member
// DELETE /api/v1/staff/:id
func (h *StaffHandler) DeleteStaff(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	staffIDStr := c.Param("id")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff ID"})
		return
	}

	err = h.staffService.DeleteStaff(staffID, companyID)
	if err != nil {
		if err.Error() == "staff not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "staff deleted successfully"})
}

// ImportStaff imports staff members from a CSV file
// POST /api/v1/staff/import
func (h *StaffHandler) ImportStaff(c *gin.Context) {
	companyID, ok := middleware.RequireCompanyID(c)
	if !ok {
		return
	}

	// Get file from multipart form (FormFile automatically parses the multipart form)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided. Use 'file' as the form field name"})
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	// Import from CSV
	result, err := h.staffService.ImportFromCSV(file, companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AnalyzeCSV analyzes a CSV file and suggests column mappings
// POST /api/v1/staff/csv/analyze
func (h *StaffHandler) AnalyzeCSV(c *gin.Context) {
	companyID, exists := middleware.GetCompanyID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "company authentication required"})
		return
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload required"})
		return
	}

	// Validate file type
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only CSV files are supported"})
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	// Analyze CSV
	result, err := h.staffService.AnalyzeCSV(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PreviewCSVImport previews a CSV import with column mappings
// POST /api/v1/staff/csv/preview
func (h *StaffHandler) PreviewCSVImport(c *gin.Context) {
	companyID, exists := middleware.GetCompanyID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "company authentication required"})
		return
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload required"})
		return
	}

	// Get mappings from form data
	mappingsJSON := c.PostForm("mappings")
	if mappingsJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "column mappings required"})
		return
	}

	// Parse mappings (simplified - in practice you'd use JSON parsing)
	// For now, we'll assume mappings are provided as a simple format
	// This would need proper JSON parsing in a real implementation

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	// For now, create a basic preview without mappings
	// This is a placeholder - full implementation would parse the mappings JSON
	preview := &services.ImportPreview{
		TotalRows:   0,
		ValidRows:   0,
		InvalidRows: 0,
		SampleRecords: []services.ImportRecord{},
		Errors:      []string{"Preview with mappings not yet implemented"},
	}

	c.JSON(http.StatusOK, preview)
}
