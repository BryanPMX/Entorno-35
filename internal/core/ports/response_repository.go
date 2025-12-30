package ports

import (
	"github.com/entorno35/backend/internal/domain"
)

// ResponseRepository defines the interface for response repository operations
type ResponseRepository interface {
	// SaveResponses saves multiple responses in a transaction
	// If any response fails to save, the entire transaction is rolled back
	SaveResponses(responses []domain.Response) error
}

