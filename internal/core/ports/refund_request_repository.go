package ports

import (
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// RefundRequestRepository persists company refund request tickets.
type RefundRequestRepository interface {
	Create(request *domain.BillingRefundRequest) error
	GetLatestOpenByCompanyID(companyID uuid.UUID) (*domain.BillingRefundRequest, error)
	List(status string, limit, offset int) ([]domain.BillingRefundRequest, int64, error)
	ListByCompanyID(companyID uuid.UUID, limit, offset int) ([]domain.BillingRefundRequest, int64, error)
	GetByID(id uuid.UUID) (*domain.BillingRefundRequest, error)
	Update(request *domain.BillingRefundRequest) error
}
