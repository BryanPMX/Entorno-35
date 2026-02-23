package ports

import (
	"time"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// PendingRegistrationRepository persists pre-payment company registration attempts.
type PendingRegistrationRepository interface {
	GetActiveByRFC(rfc string) (*domain.PendingCompanyRegistration, error)
	GetByID(id uuid.UUID) (*domain.PendingCompanyRegistration, error)
	GetByStripeSubscriptionID(subscriptionID string) (*domain.PendingCompanyRegistration, error)
	DeleteExpiredIncompleteBefore(cutoff time.Time) (int64, error)
	UpsertPending(registration *domain.PendingCompanyRegistration) error
	Update(registration *domain.PendingCompanyRegistration) error
}
