package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// companyRepository implements the CompanyRepository interface using GORM
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository creates a new Postgres implementation of CompanyRepository
func NewCompanyRepository(db *gorm.DB) ports.CompanyRepository {
	return &companyRepository{db: db}
}

// GetByID retrieves a company by ID
func (r *companyRepository) GetByID(id uuid.UUID) (*domain.Company, error) {
	var company domain.Company

	result := r.db.Where("id = ?", id).First(&company)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("failed to fetch company: %w", result.Error)
	}

	return &company, nil
}

