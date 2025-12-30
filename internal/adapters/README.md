# Adapters Package

HTTP and database adapters following Hexagonal Architecture (Ports & Adapters pattern).

## Design Principles

- **High Cohesion**: Each adapter package handles one type of external dependency
- **Low Coupling**: Adapters implement interfaces defined in `internal/core/ports`
- **Dependency Inversion**: Business logic depends on ports (interfaces), not adapters (implementations)

## Package Structure

```
internal/adapters/
├── http/           # HTTP adapters (handlers, request/response DTOs)
└── postgres/       # PostgreSQL adapters (repository implementations)
```

## Components

### HTTP Adapters (`internal/adapters/http`)

- **AuthHandler**: Handles authentication HTTP requests
  - Login endpoint
  - Request/Response DTOs

### Postgres Adapters (`internal/adapters/postgres`)

- **AuthRepository**: PostgreSQL implementation of `ports.AuthRepository`
  - Company lookup by RFC
  - Staff lookup by CURP
  - Subscription status validation

## Usage

Adapters are initialized in `cmd/api/main.go` and wired to business logic services.

```go
// Initialize repository (adapter)
authRepo := postgres.NewAuthRepository(db)

// Initialize handler (adapter)
authHandler := http.NewAuthHandler(jwtService, authRepo, tokenExpiry)

// Register routes
router.POST("/auth/login", authHandler.Login)
```

## Error Handling

Repository errors are defined in the adapter package and re-exported through the ports package to maintain low coupling.

