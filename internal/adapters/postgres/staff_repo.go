package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// Create creates a new staff member
func (r *staffRepository) Create(staff *domain.Staff) error {
	result := r.db.Create(staff)
	if result.Error != nil {
		return fmt.Errorf("failed to create staff: %w", result.Error)
	}
	return nil
}

// Update updates an existing staff member
func (r *staffRepository) Update(staff *domain.Staff) error {
	// Ensure both ID and company_id match (authorization check)
	result := r.db.Model(&domain.Staff{}).
		Where("id = ? AND company_id = ?", staff.ID, staff.CompanyID).
		Updates(map[string]interface{}{
			"full_name":    staff.FullName,
			"email":        staff.Email,
			"demographics": staff.Demographics,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update staff: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrStaffNotFound
	}
	return nil
}

// Delete soft deletes a staff member
func (r *staffRepository) Delete(id uuid.UUID, companyID uuid.UUID) error {
	result := r.db.Where("id = ? AND company_id = ?", id, companyID).Delete(&domain.Staff{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete staff: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrStaffNotFound
	}
	return nil
}

// ListByCompany retrieves all staff members for a company with optional filters
func (r *staffRepository) ListByCompany(companyID uuid.UUID, limit, offset int) ([]domain.Staff, int64, error) {
	var staff []domain.Staff
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Staff{}).Where("company_id = ?", companyID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count staff: %w", err)
	}

	// Fetch paginated records
	result := r.db.Where("company_id = ?", companyID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&staff)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list staff: %w", result.Error)
	}

	return staff, total, nil
}

// BulkCreate creates multiple staff members in a transaction
// Uses ON CONFLICT DO NOTHING for CURP collisions to avoid crashing the whole batch
// Returns the number of records successfully inserted
func (r *staffRepository) BulkCreate(staff []*domain.Staff) (int, error) {
	if len(staff) == 0 {
		return 0, nil
	}

	var insertedCount int

	// Use transaction for atomicity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Use Clauses with OnConflict to handle CURP collisions gracefully
		// The unique constraint is on (company_id, curp), so we use DoNothing
		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "company_id"}, {Name: "curp"}},
			DoNothing: true,
		}).Create(staff)

		if result.Error != nil {
			return fmt.Errorf("failed to bulk create staff: %w", result.Error)
		}

		insertedCount = int(result.RowsAffected)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return insertedCount, nil
}

