# Postgres Adapters

PostgreSQL implementations of repository interfaces following Hexagonal Architecture.

## Design Principles

- **High Cohesion**: Each repository handles one domain entity
- **Low Coupling**: Implements ports interfaces, not concrete business logic
- **Dependency Inversion**: Business logic depends on ports, not these adapters

## Components

### AuthRepository

- `GetCompanyByRFC`: Company lookup with subscription validation
- `GetStaffByCURP`: Staff lookup with company validation

### AssessmentRepository

- `GetByID`: Fetch assessment with relationships loaded
- `GetByIDAndCompany`: Fetch with company authorization check
- `GetResponsesByAssessmentID`: Fetch responses with questions fully loaded (for scoring)
- `UpdateResult`: Update assessment with scoring results

## Testing

Unit tests use SQLite in-memory database for fast, isolated testing:

```bash
go test ./internal/adapters/postgres/... -v
```

Tests cover:
- Repository operations
- Authorization checks
- Relationship loading
- Error handling

