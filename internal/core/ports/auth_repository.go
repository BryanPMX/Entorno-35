package ports

import "github.com/entorno35/backend/internal/domain"

// AuthRepository defines the interface for authentication repository operations
// (Low coupling - depends on domain types, not database implementation)
type AuthRepository interface {
	// GetCompanyByRFC retrieves a company by RFC
	// Returns error if company not found or subscription is inactive
	GetCompanyByRFC(rfc string) (*domain.Company, error)

	// GetStaffByCURP retrieves a staff member by CURP and company ID
	// Returns error if staff not found
	GetStaffByCURP(curp string, companyID string) (*domain.Staff, error)

	// GetStaffByIdentifier retrieves a staff member by either CURP or employee_id and company ID
	// Supports hybrid authentication for staff with or without CURPs
	// Returns error if staff not found
	GetStaffByIdentifier(identifier string, companyID string) (*domain.Staff, error)
}

