package ports

import (
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// ReportRepository defines the interface for report repository operations
type ReportRepository interface {
	// GetIndividualReport retrieves data for an individual assessment report
	// Returns assessment with staff demographics and calculated scores
	GetIndividualReport(assessmentID uuid.UUID, companyID uuid.UUID) (*domain.IndividualReportDTO, error)

	// GetGeneralReport retrieves aggregated data for a company-wide general report
	// Returns participation metrics, risk distribution, and department heatmap
	GetGeneralReport(companyID uuid.UUID, period *int) (*domain.GeneralReportDTO, error)
}

