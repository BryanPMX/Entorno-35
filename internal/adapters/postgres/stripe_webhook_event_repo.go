package postgres

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type stripeWebhookEventRepository struct {
	db *gorm.DB
}

func NewStripeWebhookEventRepository(db *gorm.DB) ports.StripeWebhookEventRepository {
	return &stripeWebhookEventRepository{db: db}
}

func (r *stripeWebhookEventRepository) TryBegin(eventID, eventType string) (bool, error) {
	if strings.TrimSpace(eventID) == "" {
		return false, fmt.Errorf("stripe event ID is required")
	}

	var shouldProcess bool
	err := r.db.Transaction(func(tx *gorm.DB) error {
		record := domain.StripeWebhookEvent{
			StripeEventID: eventID,
			EventType:     eventType,
			Status:        domain.StripeWebhookEventStatusProcessing,
			AttemptCount:  1,
		}
		createResult := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "stripe_event_id"}},
			DoNothing: true,
		}).Create(&record)
		if createResult.Error != nil {
			return fmt.Errorf("failed to create webhook event record: %w", createResult.Error)
		}
		if createResult.RowsAffected > 0 {
			shouldProcess = true
			return nil
		}

		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("stripe_event_id = ?", eventID).
			First(&record)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to fetch webhook event record: %w", result.Error)
			}
			return fmt.Errorf("webhook event record disappeared during processing")
		}

		if record.Status == domain.StripeWebhookEventStatusProcessed {
			shouldProcess = false
			return nil
		}

		record.EventType = eventType
		record.Status = domain.StripeWebhookEventStatusProcessing
		record.AttemptCount++
		record.LastError = nil
		record.ProcessedAt = nil
		if err := tx.Save(&record).Error; err != nil {
			return fmt.Errorf("failed to update webhook event record: %w", err)
		}

		shouldProcess = true
		return nil
	})
	if err != nil {
		return false, err
	}

	return shouldProcess, nil
}

func (r *stripeWebhookEventRepository) MarkProcessed(eventID string) error {
	now := time.Now().UTC()
	updates := map[string]interface{}{
		"status":       domain.StripeWebhookEventStatusProcessed,
		"processed_at": &now,
		"last_error":   nil,
	}
	if err := r.db.Model(&domain.StripeWebhookEvent{}).
		Where("stripe_event_id = ?", eventID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to mark webhook event as processed: %w", err)
	}
	return nil
}

func (r *stripeWebhookEventRepository) MarkFailed(eventID, failure string) error {
	if len(failure) > 2000 {
		failure = failure[:2000]
	}
	updates := map[string]interface{}{
		"status":       domain.StripeWebhookEventStatusFailed,
		"last_error":   failure,
		"processed_at": nil,
	}
	if err := r.db.Model(&domain.StripeWebhookEvent{}).
		Where("stripe_event_id = ?", eventID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to mark webhook event as failed: %w", err)
	}
	return nil
}
