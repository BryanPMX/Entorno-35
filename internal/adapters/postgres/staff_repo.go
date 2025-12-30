package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// staffRepository implements the StaffRepository interface using GORM
type staffRepository struct {
	db *gorm.DB
}

// NewStaffRepository creates a new Postgres implementation of StaffRepository
func NewStaffRepository(db *gorm.DB) ports.StaffRepository {
	return &staffRepository{db: db}
}

// GetByIDAndCompany retrieves a staff member by ID and company ID (for authorization)
func (r *staffRepository) GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Staff, error) {
	var staff domain.Staff

	result := r.db.Where("id = ? AND company_id = ?", id, companyID).First(&staff)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound
		}
		return nil, fmt.Errorf("failed to fetch staff: %w", result.Error)
	}

	return &staff, nil
}

