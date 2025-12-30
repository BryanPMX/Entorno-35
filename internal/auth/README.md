# Auth Package

Authentication context and types (high cohesion - auth domain).

## Design Principles

- **High Cohesion**: Contains only authentication-related domain types
- **Low Coupling**: Simple value types, no dependencies on other packages
- **Single Responsibility**: Defines the structure of auth context data

## Usage

```go
import "github.com/entorno35/backend/internal/auth"

// Auth context is set by middleware and retrieved via:
authCtx, exists := middleware.GetAuthContext(c)

if authCtx.IsStaff() {
    staffID := authCtx.GetStaffID()
}
```

