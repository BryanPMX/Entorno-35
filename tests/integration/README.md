# Integration Tests

Integration tests that require a running PostgreSQL database.

## Test Organization

### Structure
- `main_test.go` - Shared test setup (TestMain, database connection, JWT service)
- `report_e2e_test.go` - E2E tests for the Reporting Engine ("The Compliance Audit")
- `staff_import_test.go` - Staff CSV import functionality tests

### Running Tests

```bash
# Run all integration tests
go test -tags=integration -v ./tests/integration/...

# Run specific test
go test -tags=integration -v ./tests/integration/... -run TestReportE2E

# Run with coverage
go test -tags=integration -v ./tests/integration/... -cover
```

## Prerequisites

1. **PostgreSQL Database**: Must be running and accessible
   ```bash
   docker-compose up -d
   # Or use existing PostgreSQL instance
   ```

2. **Environment Variables** (Required):
   ```bash
   # DB_URL is REQUIRED - integration tests will fail without it
   # Never commit passwords to the repository - use environment variables only
   export DB_URL=postgres://user:password@localhost:5432/entorno35_test?sslmode=disable
   export JWT_SECRET=your-jwt-secret-key  # Optional, defaults to test secret
   export API_BASE_URL=http://localhost:8080  # Optional, for HTTP endpoint tests
   ```
   
   **Security Note**: The `DB_URL` environment variable is required. Integration tests will fail immediately if it's not set. This prevents hardcoded credentials from being committed to the repository.

## Test Database

Tests use a separate test database (`entorno35_test`) to avoid conflicts with development data.

**Important**: Tests create test companies and clean them up after execution using `CleanupTestCompany()`. However, if tests are interrupted, you may need to manually clean up:

```sql
DELETE FROM responses WHERE assessment_id IN (
    SELECT id FROM assessments WHERE company_id IN (
        SELECT id FROM companies WHERE name LIKE 'Test%'
    )
);
DELETE FROM assessments WHERE company_id IN (
    SELECT id FROM companies WHERE name LIKE 'Test%'
);
DELETE FROM staff WHERE company_id IN (
    SELECT id FROM companies WHERE name LIKE 'Test%'
);
DELETE FROM companies WHERE name LIKE 'Test%';
```

## Test Cases

### Report E2E Tests (`TestReportE2E_TheComplianceAudit`)

**"The Compliance Audit"** - Comprehensive E2E test for the Reporting Engine:

1. **Setup (`setupAuditData`)**:
   - Creates a unique test company
   - Creates 2 staff members:
     - Staff A: IT Department, High Risk responses
     - Staff B: HR Department, Low Risk responses
   - Creates assessments for both staff
   - Submits responses (simulates completed tests)
   - Calculates scores using the scoring service

2. **Test A: Individual Report**:
   - Tests `GET /api/v1/reports/individual/:assessment_id`
   - Asserts HTTP 200
   - Asserts risk level matches expected (Alto/Muy Alto for high risk)
   - **Critical**: Asserts recommendations array is NOT empty (NOM-035 requirement)
   - Asserts domain scores are populated

3. **Test B: General Report (Aggregation)**:
   - Tests `GET /api/v1/reports/general`
   - Asserts HTTP 200
   - Asserts participation rate reflects 2 users (100%)
   - Asserts department heatmap contains "IT" and "HR"
   - Asserts risk distribution counts match seeded data

4. **Test C: HTTP Endpoints** (optional):
   - Tests actual HTTP endpoints if API server is running
   - Validates JSON response structure

5. **Teardown**:
   - Hard deletes test company and all related data (cascading)
   - Ensures cleanup even if test fails

### Staff Import Tests

See individual test functions in `staff_import_test.go`:
- `TestStaffImport_HappyPath` - Large batch import (2,500 rows)
- `TestStaffImport_SpanishCharacters` - UTF-8 encoding
- `TestStaffImport_DuplicateHandling` - Idempotency
- `TestStaffImport_MixedValidInvalid` - Partial success
- `TestStaffImport_TenantIsolation` - Multi-tenant security
- `TestStaffImport_UnauthorizedAccess` - Authentication

## Shared Test Infrastructure

### Global Variables (from `main_test.go`)
- `TestDB` - Global GORM database connection
- `JWTService` - JWT service for authentication
- `BaseURL` - API base URL (optional)

### Helper Functions
- `CreateAuthToken(companyID string)` - Generate JWT token
- `CleanupTestCompany(t *testing.T, companyID string)` - Hard delete test company

## Notes

- Tests are marked with `//go:build integration` to separate them from unit tests
- Tests require a real PostgreSQL database (not SQLite) to test JSONB queries and complex aggregations
- Tests clean up after themselves, but may leave data if interrupted
- HTTP endpoint tests are optional and will skip if API server is not running
