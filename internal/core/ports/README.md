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
// excerpt
type AuthRepository interface {
    GetCompanyByRFC(rfc string) (*domain.Company, error)
    GetStaffByCURP(curp string, companyID string) (*domain.Staff, error)
}
```

**Implementation**: `internal/adapters/postgres/auth_repo.go`

Current `AuthRepository` also includes company lifecycle methods used by billing:

- company create/update
- lookup by company ID
- lookup by Stripe subscription ID

### ReportRepository

Interface for report repository operations:

```go
type ReportRepository interface {
    GetIndividualReport(assessmentID uuid.UUID, companyID uuid.UUID) (*domain.IndividualReportDTO, error)
    GetGeneralReport(companyID uuid.UUID, period *int) (*domain.GeneralReportDTO, error)
}
```

**Implementation**: `internal/adapters/postgres/report_repo.go`

### PendingRegistrationRepository

Interface for pre-payment company registration drafts used by Stripe onboarding:

- create/update pending registration drafts
- lookup by RFC / ID / Stripe subscription ID
- cleanup expired incomplete drafts

**Implementation**: `internal/adapters/postgres/pending_registration_repo.go`

### StripeWebhookEventRepository

Interface for webhook idempotency state:

- begin processing by Stripe `event.id`
- mark processed / failed

**Implementation**: `internal/adapters/postgres/stripe_webhook_event_repo.go`

## Error Handling Note

This package defines interfaces only. Concrete repository error types are currently declared in adapter packages (for example `internal/adapters/postgres`) and handled explicitly where needed.
