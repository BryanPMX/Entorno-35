package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrBillingRefundRequestNotFound = errors.New("billing refund request not found")
var ErrBillingRefundRequestAlreadyOpen = errors.New("billing refund request already open")

type refundRequestRepository struct {
	db *gorm.DB
}

func NewRefundRequestRepository(db *gorm.DB) ports.RefundRequestRepository {
	return &refundRequestRepository{db: db}
}

func (r *refundRequestRepository) Create(request *domain.BillingRefundRequest) error {
	if request == nil {
		return fmt.Errorf("refund request is required")
	}

	request.Reason = strings.TrimSpace(request.Reason)
	if request.Reason == "" {
		return fmt.Errorf("refund reason is required")
	}

	if request.Status == "" {
		request.Status = domain.BillingRefundRequestStatusRequested
	}

	if err := r.db.Create(request).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "uq_billing_refund_requests_open_company") {
			return ErrBillingRefundRequestAlreadyOpen
		}
		return fmt.Errorf("failed to create refund request: %w", err)
	}

	return nil
}

func (r *refundRequestRepository) GetLatestOpenByCompanyID(companyID uuid.UUID) (*domain.BillingRefundRequest, error) {
	var request domain.BillingRefundRequest

	result := r.db.
		Where("company_id = ? AND status = ?", companyID, domain.BillingRefundRequestStatusRequested).
		Order("created_at DESC").
		First(&request)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrBillingRefundRequestNotFound
		}
		return nil, fmt.Errorf("failed to fetch refund request by company ID: %w", result.Error)
	}

	return &request, nil
}

func (r *refundRequestRepository) List(status string, limit, offset int) ([]domain.BillingRefundRequest, int64, error) {
	query := r.db.Model(&domain.BillingRefundRequest{})
	if strings.TrimSpace(status) != "" {
		query = query.Where("status = ?", strings.ToLower(strings.TrimSpace(status)))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count refund requests: %w", err)
	}

	var requests []domain.BillingRefundRequest
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&requests).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list refund requests: %w", err)
	}

	return requests, total, nil
}

func (r *refundRequestRepository) ListByCompanyID(companyID uuid.UUID, limit, offset int) ([]domain.BillingRefundRequest, int64, error) {
	query := r.db.Model(&domain.BillingRefundRequest{}).Where("company_id = ?", companyID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count company refund requests: %w", err)
	}

	var requests []domain.BillingRefundRequest
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&requests).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list company refund requests: %w", err)
	}

	return requests, total, nil
}

func (r *refundRequestRepository) GetByID(id uuid.UUID) (*domain.BillingRefundRequest, error) {
	var request domain.BillingRefundRequest
	result := r.db.Where("id = ?", id).First(&request)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrBillingRefundRequestNotFound
		}
		return nil, fmt.Errorf("failed to fetch refund request by ID: %w", result.Error)
	}
	return &request, nil
}

func (r *refundRequestRepository) Update(request *domain.BillingRefundRequest) error {
	if request == nil {
		return fmt.Errorf("refund request is required")
	}
	if err := r.db.Save(request).Error; err != nil {
		return fmt.Errorf("failed to update refund request: %w", err)
	}
	return nil
}
