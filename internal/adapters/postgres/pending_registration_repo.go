package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrPendingRegistrationNotFound = errors.New("pending registration not found")

type pendingRegistrationRepository struct {
	db *gorm.DB
}

func NewPendingRegistrationRepository(db *gorm.DB) ports.PendingRegistrationRepository {
	return &pendingRegistrationRepository{db: db}
}

func (r *pendingRegistrationRepository) GetActiveByRFC(rfc string) (*domain.PendingCompanyRegistration, error) {
	var registration domain.PendingCompanyRegistration
	result := r.db.
		Where("rfc = ? AND completed_at IS NULL AND expires_at > ?", rfc, time.Now().UTC()).
		Order("created_at DESC").
		First(&registration)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPendingRegistrationNotFound
		}
		return nil, fmt.Errorf("failed to fetch pending registration by RFC: %w", result.Error)
	}
	return &registration, nil
}

func (r *pendingRegistrationRepository) GetByID(id uuid.UUID) (*domain.PendingCompanyRegistration, error) {
	var registration domain.PendingCompanyRegistration
	result := r.db.Where("id = ?", id).First(&registration)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPendingRegistrationNotFound
		}
		return nil, fmt.Errorf("failed to fetch pending registration by ID: %w", result.Error)
	}
	return &registration, nil
}

func (r *pendingRegistrationRepository) GetByStripeSubscriptionID(subscriptionID string) (*domain.PendingCompanyRegistration, error) {
	var registration domain.PendingCompanyRegistration
	result := r.db.
		Where("stripe_subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		First(&registration)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPendingRegistrationNotFound
		}
		return nil, fmt.Errorf("failed to fetch pending registration by subscription ID: %w", result.Error)
	}
	return &registration, nil
}

func (r *pendingRegistrationRepository) DeleteExpiredIncompleteBefore(cutoff time.Time) (int64, error) {
	result := r.db.Unscoped().
		Where("completed_at IS NULL AND expires_at < ?", cutoff.UTC()).
		Delete(&domain.PendingCompanyRegistration{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete expired incomplete pending registrations: %w", result.Error)
	}
	return result.RowsAffected, nil
}

func (r *pendingRegistrationRepository) UpsertPending(registration *domain.PendingCompanyRegistration) error {
	if registration == nil {
		return fmt.Errorf("pending registration is required")
	}

	existing, err := r.GetActiveByRFC(registration.RFC)
	if err != nil && !errors.Is(err, ErrPendingRegistrationNotFound) {
		return err
	}

	if err == nil {
		if saveErr := r.applyPendingDraft(existing, registration); saveErr != nil {
			return saveErr
		}
		*registration = *existing
		return nil
	}

	// Reuse the latest incomplete row even if expired to avoid unique-index conflicts on retries.
	var incomplete domain.PendingCompanyRegistration
	incompleteResult := r.db.
		Where("rfc = ? AND completed_at IS NULL", registration.RFC).
		Order("created_at DESC").
		First(&incomplete)
	if incompleteResult.Error == nil {
		if saveErr := r.applyPendingDraft(&incomplete, registration); saveErr != nil {
			return saveErr
		}
		*registration = incomplete
		return nil
	}
	if incompleteResult.Error != nil && !errors.Is(incompleteResult.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch incomplete pending registration: %w", incompleteResult.Error)
	}

	if createErr := r.db.Create(registration).Error; createErr != nil {
		return fmt.Errorf("failed to create pending registration: %w", createErr)
	}
	return nil
}

func (r *pendingRegistrationRepository) applyPendingDraft(existing, incoming *domain.PendingCompanyRegistration) error {
	if existing == nil || incoming == nil {
		return fmt.Errorf("pending registrations are required")
	}

	existing.CompanyName = incoming.CompanyName
	existing.Address = incoming.Address
	existing.AdminEmail = incoming.AdminEmail
	existing.AdminPasswordHash = incoming.AdminPasswordHash
	existing.EmployeeCount = incoming.EmployeeCount
	existing.Plan = incoming.Plan
	existing.ExpiresAt = incoming.ExpiresAt
	existing.StripeCheckoutSessionID = nil
	existing.StripeCustomerID = nil
	existing.StripeSubscriptionID = nil
	existing.CompanyID = nil
	existing.CompletedAt = nil

	if saveErr := r.db.Save(existing).Error; saveErr != nil {
		return fmt.Errorf("failed to update pending registration: %w", saveErr)
	}
	return nil
}

func (r *pendingRegistrationRepository) Update(registration *domain.PendingCompanyRegistration) error {
	if registration == nil {
		return fmt.Errorf("pending registration is required")
	}
	if err := r.db.Save(registration).Error; err != nil {
		return fmt.Errorf("failed to save pending registration: %w", err)
	}
	return nil
}
