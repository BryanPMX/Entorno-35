package ports

import (
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// StaffRepository defines the interface for staff repository operations
type StaffRepository interface {
	// GetByIDAndCompany retrieves a staff member by ID and company ID (for authorization)
	GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Staff, error)

	// Create creates a new staff member
	Create(staff *domain.Staff) error

	// Update updates an existing staff member
	Update(staff *domain.Staff) error

	// Delete soft deletes a staff member
	Delete(id uuid.UUID, companyID uuid.UUID) error

	// ListByCompany retrieves all staff members for a company with optional filters
	ListByCompany(companyID uuid.UUID, limit, offset int) ([]domain.Staff, int64, error)

	// BulkCreate creates multiple staff members in a transaction
	// Uses ON CONFLICT DO NOTHING for CURP collisions to avoid crashing the whole batch
	// Returns the number of records successfully inserted
	BulkCreate(staff []*domain.Staff) (int, error)

	// HasCompletedAssessments checks if a staff member has any completed assessments
	HasCompletedAssessments(staffID uuid.UUID) (bool, error)
}

