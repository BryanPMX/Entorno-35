# PostgreSQL Repository Adapters

This directory contains PostgreSQL implementation of repository interfaces using GORM ORM.

## Architecture

Implements the Repository Pattern to abstract data persistence concerns from business logic. Each repository provides CRUD operations and domain-specific queries for its entity.

### Repository Implementations

- **auth_repo.go**: Authentication and company verification
- **company_repo.go**: Company CRUD operations  
- **staff_repo.go**: Staff management and bulk operations
- **assessment_repo.go**: Assessment lifecycle and response management
- **response_repo.go**: Assessment response persistence
- **report_repo.go**: Report generation with dynamic scoring calculations

## Testing Strategy

### Unit Tests

Repository tests use an **in-memory SQLite database** for fast, isolated testing without requiring external dependencies.

#### Test Files

- `assessment_repo_test.go`: Assessment CRUD and query operations
- `report_repo_test.go`: Report generation, scoring calculations, and risk level determination

#### Running Tests

```bash
# Run all repository tests
go test ./internal/adapters/postgres -v

# Run specific repository tests
go test ./internal/adapters/postgres -run TestReportRepository -v
go test ./internal/adapters/postgres -run TestAssessmentRepository -v

# Run with coverage
go test ./internal/adapters/postgres -cover
```

### Test Database Setup

Tests use SQLite with manually created schemas that mirror the PostgreSQL production schema. This approach provides:

- Fast test execution (in-memory database)
- No external dependencies
- Deterministic test data
- Transaction isolation between tests

**Note**: One test (`TestReportRepository_GetIndividualReport_DomainRiskLevelCalculation`) is skipped due to SQLite relationship loading limitations. This test verifies complex GORM relationship preloading which works correctly in production PostgreSQL but has limitations in SQLite's foreign key handling.

### Test Coverage

Current test coverage focuses on:

- **Assessment Repository**: CRUD operations, question loading, response management
- **Report Repository**: 
  - Individual report generation
  - Domain and category score calculations
  - Dynamic max score computation
  - Risk level determination
  - General report aggregations
  - Error handling (not found, not scored)

### Integration Tests

For end-to-end testing with actual PostgreSQL database, see:
- `tests/integration/report_e2e_test.go`
- `tests/integration/staff_import_test.go`

These tests require a running PostgreSQL instance and test the complete flow from API endpoints through services to repositories.

## Error Handling

Package-level error variables provide semantic error types:

```go
var (
    ErrCompanyNotFound = errors.New("company not found")
    ErrStaffNotFound   = errors.New("staff not found")
    ErrReportNotFound  = errors.New("report not found")
    // ... additional errors
)
```

## Database Migrations

Schema migrations are managed separately in `/migrations` directory. Repositories assume the schema is already migrated.

## GORM Configuration

- **Soft Deletes**: Enabled via `DeletedAt` field using `gorm.DeletedAt`
- **Preloading**: Explicit relationship preloading for performance
- **Transactions**: Used for bulk operations and data consistency
- **JSON Fields**: PostgreSQL JSONB for flexible demographic data

## Best Practices

1. **Repository Methods**: Follow SOLID principles
   - Single responsibility per method
   - Return domain entities, not ORM models
   - Handle errors explicitly

2. **Query Optimization**:
   - Use preloading to avoid N+1 queries
   - Apply indexes via migrations
   - Use pagination for list operations

3. **Testing**:
   - Isolate each test with fresh database
   - Use table-driven tests for multiple scenarios
   - Mock external dependencies, not database
   - Test error paths, not just happy path

4. **Error Handling**:
   - Return semantic errors, not GORM errors
   - Wrap errors with context using fmt.Errorf
   - Check for specific errors (e.g., gorm.ErrRecordNotFound)

## Future Improvements

- Add benchmark tests for query performance
- Implement query result caching for read-heavy operations
- Add database connection pooling configuration
- Implement audit logging for critical operations
