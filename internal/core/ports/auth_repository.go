package ports

import (
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// AuthRepository defines the interface for authentication repository operations
// (Low coupling - depends on domain types, not database implementation)
type AuthRepository interface {
	// GetCompanyByRFC retrieves a company by RFC
	// Returns error if company not found
	GetCompanyByRFC(rfc string) (*domain.Company, error)

	// CreateCompany creates a new company
	CreateCompany(company *domain.Company) error

	// UpdateCompany updates an existing company record
	UpdateCompany(company *domain.Company) error

	// GetCompanyByID retrieves a company by ID
	GetCompanyByID(id uuid.UUID) (*domain.Company, error)

	// GetCompanyByStripeSubscriptionID retrieves a company by Stripe subscription ID
	GetCompanyByStripeSubscriptionID(subscriptionID string) (*domain.Company, error)

	// GetStaffByCURP retrieves a staff member by CURP and company ID
	// Returns error if staff not found
	GetStaffByCURP(curp string, companyID string) (*domain.Staff, error)

	// GetStaffByIdentifier retrieves a staff member by either CURP or employee_id and company ID
	// Supports hybrid authentication for staff with or without CURPs
	// Returns error if staff not found
	GetStaffByIdentifier(identifier string, companyID string) (*domain.Staff, error)
}
