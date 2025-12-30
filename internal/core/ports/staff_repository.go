package ports

import (
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// StaffRepository defines the interface for staff repository operations
type StaffRepository interface {
	// GetByIDAndCompany retrieves a staff member by ID and company ID (for authorization)
	GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Staff, error)
}

