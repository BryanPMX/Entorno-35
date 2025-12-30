package ports

import "github.com/entorno35/backend/internal/adapters/postgres"

// Re-export postgres adapter errors for use by handlers
// This allows handlers to depend on ports package, not adapters (low coupling)
var (
	ErrCompanyNotFound = postgres.ErrCompanyNotFound
	ErrCompanyInactive = postgres.ErrCompanyInactive
	ErrStaffNotFound   = postgres.ErrStaffNotFound
)

