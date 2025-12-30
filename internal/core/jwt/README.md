# JWT Service

High cohesion, low coupling JWT service implementation.

## Design Principles

- **High Cohesion**: Single responsibility - handles only JWT token operations
- **Low Coupling**: Uses interface (`Service`) to allow different implementations
- **Dependency Inversion**: Consumers depend on `jwt.Service` interface, not concrete implementation

## Usage

```go
import "github.com/entorno35/backend/internal/core/jwt"

// Create service
jwtService := jwt.NewService(config.JWT.Secret)

// Generate token
token, err := jwtService.GenerateToken(
    companyID,
    staffID,
    email,
    role,
    24 * time.Hour, // expiration
)

// Validate token
claims, err := jwtService.ValidateToken(tokenString)

// Refresh token
newToken, err := jwtService.RefreshToken(tokenString, 24 * time.Hour)
```

## Interface

The `Service` interface allows for:
- Easy testing (mock implementation)
- Future implementations (e.g., RS256, different token formats)
- Dependency injection

