package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"

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
// Uses ON CONFLICT DO NOTHING for employee_id collisions to avoid crashing the whole batch
// Returns the number of records successfully inserted
func (r *staffRepository) BulkCreate(staff []*domain.Staff) (int, error) {
	if len(staff) == 0 {
		return 0, nil
	}

	var insertedCount int

	// Use transaction for atomicity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Use Clauses with OnConflict to handle employee_id collisions gracefully
		// The unique constraint is on employee_id, so we use DoNothing
		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "employee_id"}},
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

// HasCompletedAssessments checks if a staff member has any completed assessments
func (r *staffRepository) HasCompletedAssessments(staffID uuid.UUID) (bool, error) {
	var count int64

	result := r.db.Model(&domain.Assessment{}).
		Where("staff_id = ? AND status = ?", staffID, domain.AssessmentStatusCompleted).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check for completed assessments: %w", result.Error)
	}

	return count > 0, nil
}

// GenerateAndReserveEmployeeID atomically generates and reserves a unique employee ID
// Uses PostgreSQL advisory locks to prevent race conditions across concurrent requests
func (r *staffRepository) GenerateAndReserveEmployeeID(companyID uuid.UUID, prefix string) (string, error) {
	var employeeID string
	prefixPattern := prefix + "%"

	// Generate a unique lock key based on company ID
	// PostgreSQL advisory locks use bigint, so we hash the UUID string to int64
	// This ensures the same company always uses the same lock key
	h := fnv.New64a()
	h.Write([]byte(companyID.String()))
	lockKey := int64(h.Sum64())
	// Ensure positive value for advisory lock (PostgreSQL requires positive)
	if lockKey < 0 {
		lockKey = -lockKey
	}

	// Use transaction with advisory lock to ensure atomicity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Acquire advisory lock (blocks until available, prevents concurrent ID generation)
		// pg_advisory_xact_lock is automatically released when transaction ends
		lockQuery := `SELECT pg_advisory_xact_lock(?)`
		if err := tx.Exec(lockQuery, lockKey).Error; err != nil {
			return fmt.Errorf("failed to acquire advisory lock: %w", err)
		}

		var maxSeq sql.NullInt64

		// Now safely query for max sequence (we have exclusive lock)
		query := `
			SELECT MAX(
				CAST(
					SUBSTRING(employee_id FROM POSITION('-' IN employee_id) + 1) AS INTEGER
				)
			) as max_seq
			FROM staff
			WHERE company_id = ?
			  AND employee_id IS NOT NULL
			  AND employee_id LIKE ?
			  AND deleted_at IS NULL
		`

		if err := tx.Raw(query, companyID, prefixPattern).Scan(&maxSeq).Error; err != nil {
			return fmt.Errorf("failed to get max employee ID sequence: %w", err)
		}

		// Calculate next sequence number
		nextSequence := 1
		if maxSeq.Valid {
			nextSequence = int(maxSeq.Int64) + 1
		}

		// Format as 4-digit zero-padded number
		employeeID = fmt.Sprintf("%s%04d", prefix, nextSequence)

		// Double-check the ID doesn't exist (shouldn't happen with lock, but safety check)
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM staff WHERE employee_id = ? AND deleted_at IS NULL)`
		if err := tx.Raw(checkQuery, employeeID).Scan(&exists).Error; err != nil {
			return fmt.Errorf("failed to verify employee ID availability: %w", err)
		}

		if exists {
			// If it exists, increment (shouldn't happen with proper locking, but handle it)
			nextSequence++
			employeeID = fmt.Sprintf("%s%04d", prefix, nextSequence)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return employeeID, nil
}

// CreateWithEmployeeID generates an employee ID and creates the staff record atomically
// This prevents race conditions where ID generation and insertion happen separately
func (r *staffRepository) CreateWithEmployeeID(staff *domain.Staff, companyID uuid.UUID, prefix string) error {
	prefixPattern := prefix + "%"

	// Generate a unique lock key based on company ID
	h := fnv.New64a()
	h.Write([]byte(companyID.String()))
	lockKey := int64(h.Sum64())
	if lockKey < 0 {
		lockKey = -lockKey
	}

	// Use transaction with advisory lock to ensure atomicity
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Acquire advisory lock (blocks until available, prevents concurrent ID generation)
		lockQuery := `SELECT pg_advisory_xact_lock(?)`
		if err := tx.Exec(lockQuery, lockKey).Error; err != nil {
			return fmt.Errorf("failed to acquire advisory lock: %w", err)
		}

		var maxSeq sql.NullInt64

		// Query for max sequence (we have exclusive lock)
		query := `
			SELECT MAX(
				CAST(
					SUBSTRING(employee_id FROM POSITION('-' IN employee_id) + 1) AS INTEGER
				)
			) as max_seq
			FROM staff
			WHERE company_id = ?
			  AND employee_id IS NOT NULL
			  AND employee_id LIKE ?
			  AND deleted_at IS NULL
		`

		if err := tx.Raw(query, companyID, prefixPattern).Scan(&maxSeq).Error; err != nil {
			return fmt.Errorf("failed to get max employee ID sequence: %w", err)
		}

		// Calculate next sequence number
		nextSequence := 1
		if maxSeq.Valid {
			nextSequence = int(maxSeq.Int64) + 1
		}

		// Format as 4-digit zero-padded number
		employeeID := fmt.Sprintf("%s%04d", prefix, nextSequence)

		// Double-check the ID doesn't exist (shouldn't happen with lock, but safety check)
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM staff WHERE employee_id = ? AND deleted_at IS NULL)`
		if err := tx.Raw(checkQuery, employeeID).Scan(&exists).Error; err != nil {
			return fmt.Errorf("failed to verify employee ID availability: %w", err)
		}

		if exists {
			// If it exists, increment (shouldn't happen with proper locking, but handle it)
			nextSequence++
			employeeID = fmt.Sprintf("%s%04d", prefix, nextSequence)
		}

		// Set the employee ID on the staff record
		staff.EmployeeID = sql.NullString{String: employeeID, Valid: true}

		// Create the staff record in the same transaction
		if err := tx.Create(staff).Error; err != nil {
			return fmt.Errorf("failed to create staff: %w", err)
		}

		return nil
	})
}

