# Ports Package

Defines interfaces (ports) for external dependencies following Hexagonal Architecture.

## Design Principles

- **High Cohesion**: Contains only interface definitions
- **Low Coupling**: Business logic depends on these interfaces, not concrete implementations
- **Dependency Inversion**: Enables dependency inversion principle

## Components

### AuthRepository

Interface for authentication repository operations:

```go
type AuthRepository interface {
    GetCompanyByRFC(rfc string) (*domain.Company, error)
    GetStaffByCURP(curp string, companyID string) (*domain.Staff, error)
}
```

**Implementation**: `internal/adapters/postgres/auth_repo.go`

## Error Re-exports

Repository errors are re-exported from adapters to allow handlers to depend on ports package only:

```go
import "github.com/entorno35/backend/internal/core/ports"

if err == ports.ErrCompanyNotFound {
    // Handle error
}
```

This maintains low coupling - handlers don't need to import adapter packages.

