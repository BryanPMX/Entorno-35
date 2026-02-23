# Core Package

Core business logic and interfaces following Hexagonal Architecture principles.

## Structure

```
core/
├── jwt/              # JWT token service
├── password/         # Password hashing service
├── ports/            # Port interfaces (hexagonal architecture)
├── scoring/          # NOM-035 scoring logic
└── services/         # Business logic services (assessment, reporting, Stripe, cleanup)
```

## Design Principles

- **High Cohesion**: Each package handles one domain of business logic
- **Low Coupling**: Dependencies on interfaces (ports), not concrete implementations
- **Dependency Inversion**: Business logic defines interfaces, adapters implement them
- **Pure Functions**: Core logic is testable without external dependencies

## Package Overview

### jwt/
JWT token generation, validation, and refresh operations. Provides interface-based service for authentication.

### password/
Password hashing and verification using bcrypt. Interface-based design allows for different hashing algorithms.

### ports/
Interface definitions for external dependencies (repositories, external services). Following Hexagonal Architecture, these are the "ports" that adapters implement.

### scoring/
NOM-035 scoring logic including:
- Strategy pattern implementations for different guide types
- Polarity inversion logic
- Risk level calculation
- Scoring rules and thresholds

### services/
Business logic services that orchestrate repositories and core logic. These services implement use cases and coordinate between different layers.

Current notable services include:

- `assessment_service.go` / `staff_service.go` / `report_service.go`
- `stripe_service.go` (Stripe Checkout, subscription fetch, webhook signature verification, Billing Portal sessions)
- `pending_registration_cleanup_service.go` (background cleanup for expired unpaid registration drafts)
