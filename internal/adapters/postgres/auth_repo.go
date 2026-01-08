package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"gorm.io/gorm"
)

// Package-level errors (following Go error handling best practices)
var (
	ErrCompanyNotFound = errors.New("company not found")
	ErrCompanyInactive = errors.New("company subscription is inactive")
	ErrStaffNotFound   = errors.New("staff not found")
)

// authRepository implements the AuthRepository interface using GORM (High cohesion - persistence concerns only)
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new Postgres implementation of AuthRepository
func NewAuthRepository(db *gorm.DB) ports.AuthRepository {
	return &authRepository{db: db}
}

// GetCompanyByRFC retrieves a company by RFC and validates subscription status
func (r *authRepository) GetCompanyByRFC(rfc string) (*domain.Company, error) {
	var company domain.Company

	result := r.db.Where("rfc = ?", rfc).First(&company)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("failed to fetch company: %w", result.Error)
	}

	// Critical: Check subscription status - don't allow login if inactive
	if company.SubscriptionStatus != domain.SubscriptionStatusActive {
		return nil, ErrCompanyInactive
	}

	return &company, nil
}

// GetStaffByCURP retrieves a staff member by CURP and company ID
func (r *authRepository) GetStaffByCURP(curp string, companyID string) (*domain.Staff, error) {
	var staff domain.Staff

	// First, verify company exists and is active
	var company domain.Company
	result := r.db.Where("id = ?", companyID).First(&company)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("failed to fetch company: %w", result.Error)
	}

	if company.SubscriptionStatus != domain.SubscriptionStatusActive {
		return nil, ErrCompanyInactive
	}

	// Then fetch staff member
	result = r.db.Where("curp = ? AND company_id = ?", curp, companyID).First(&staff)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", result.Error)
	}

	// Set company for reference
	staff.Company = company

	return &staff, nil
}

// GetStaffByIdentifier retrieves a staff member by either CURP or employee_id and company ID
// Supports hybrid authentication for staff with or without CURPs
func (r *authRepository) GetStaffByIdentifier(identifier string, companyID string) (*domain.Staff, error) {
	var staff domain.Staff

	// First, verify company exists and is active
	var company domain.Company
	result := r.db.Where("id = ?", companyID).First(&company)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("failed to fetch company: %w", result.Error)
	}

	if company.SubscriptionStatus != domain.SubscriptionStatusActive {
		return nil, ErrCompanyInactive
	}

	// Try to find staff by CURP first (if identifier looks like a CURP - 18 characters)
	if len(identifier) == 18 {
		result = r.db.Where("curp = ? AND company_id = ?", identifier, companyID).First(&staff)
		if result.Error == nil {
			// Found by CURP
			staff.Company = company
			return &staff, nil
		}
		// If CURP lookup failed but error is not "not found", return the error
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to fetch staff by CURP: %w", result.Error)
		}
	}

	// Try to find staff by employee_id
	result = r.db.Where("employee_id = ? AND company_id = ?", identifier, companyID).First(&staff)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound
		}
		return nil, fmt.Errorf("failed to fetch staff by employee ID: %w", result.Error)
	}

	// Set company for reference
	staff.Company = company

	return &staff, nil
}
