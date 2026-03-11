# Adapters Package

HTTP and database adapters following Hexagonal Architecture (Ports & Adapters pattern).

## Design Principles

- **High Cohesion**: Each adapter package handles one type of external dependency
- **Low Coupling**: Adapters implement interfaces defined in `internal/core/ports`
- **Dependency Inversion**: Business logic depends on ports (interfaces), not adapters (implementations)

## Package Structure

```
internal/adapters/
├── http/           # HTTP adapters (handlers, request/response DTOs, billing endpoints)
└── postgres/       # PostgreSQL adapters (repository implementations, billing state/idempotency)
```

## Components

### HTTP Adapters (`internal/adapters/http`)

- **AuthHandler**: Handles authentication HTTP requests
  - Login endpoint
  - Request/Response DTOs
- **BillingHandler**: Handles Stripe billing flows
  - Public registration checkout/webhook/verification
  - Authenticated existing-company checkout and customer portal

### Postgres Adapters (`internal/adapters/postgres`)

- **AuthRepository**: PostgreSQL implementation of `ports.AuthRepository`
  - Company lookup by RFC
  - Staff lookup by CURP
  - Subscription status validation
- **PendingRegistrationRepository**: Pre-payment registration drafts and cleanup operations
- **StripeWebhookEventRepository**: Webhook idempotency persistence

## Usage

Adapters are initialized in `cmd/api/main.go` and wired to business logic services.

```go
// Initialize repository (adapter)
authRepo := postgres.NewAuthRepository(db)

// Initialize handlers (adapters)
authHandler := http.NewAuthHandler(jwtService, authRepo, tokenExpiry, cfg.Admin.BillingEmail, cfg.Admin.BillingPasswordHash)
billingHandler := http.NewBillingHandler(authRepo, pendingRegistrationRepo, stripeWebhookEventRepo, refundRequestRepo, stripeService, successURL, cancelURL, portalReturnURL)

// Register routes
router.POST("/auth/login", authHandler.Login)
router.POST("/auth/admin/login", authHandler.AdminLogin)
router.POST("/billing/checkout-session", billingHandler.CreateCheckoutSession)
```

## Error Handling

Repository errors are defined in adapter packages (for example `internal/adapters/postgres`) and are handled explicitly by callers where semantic branching is required.
