package ports

import (
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// CompanyRepository defines the interface for company repository operations
type CompanyRepository interface {
	// GetByID retrieves a company by ID
	GetByID(id uuid.UUID) (*domain.Company, error)
}

