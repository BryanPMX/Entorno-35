# Integration Tests

Integration tests that require a running PostgreSQL database and API server.

## Prerequisites

1. **PostgreSQL Database**: Must be running and accessible
   ```bash
   docker-compose up -d
   # Or use existing PostgreSQL instance
   ```

2. **API Server**: Must be running (optional, but recommended for full integration)
   ```bash
   go run cmd/api/main.go
   ```

## Configuration

Set environment variables:

```bash
# Database connection (defaults to localhost test database)
export DB_URL=postgres://entorno35:entorno35@localhost:5432/entorno35_test?sslmode=disable

# JWT secret (defaults to test secret)
export JWT_SECRET=your-jwt-secret-key

# API base URL (defaults to http://localhost:8080)
export API_BASE_URL=http://localhost:8080
```

## Running Tests

### Run All Integration Tests

```bash
go test -tags=integration ./tests/integration/... -v
```

### Run Specific Test

```bash
go test -tags=integration ./tests/integration/... -run TestStaffImport_HappyPath -v
```

### Run with Coverage

```bash
go test -tags=integration ./tests/integration/... -cover
```

## Test Database Setup

Tests use a separate test database (`entorno35_test`) to avoid conflicts with development data.

**Important**: Tests create test companies and clean them up after execution. However, if tests are interrupted, you may need to manually clean up:

```sql
DELETE FROM staff WHERE company_id IN (
    SELECT id FROM companies WHERE name LIKE 'Test Company%'
);
DELETE FROM companies WHERE name LIKE 'Test Company%';
```

## Test Cases

### TestStaffImport_HappyPath
- Tests successful import of 2,500 valid rows (forces batch processing with batch size 1000)
- Verifies all records are inserted correctly across multiple batches
- Validates data integrity and batch processing behavior
- Queries database directly to confirm persistence

### TestStaffImport_SpanishCharacters
- Tests UTF-8 encoding with Spanish characters (á, é, í, ó, ú, ñ, Ñ)
- Verifies special characters are stored correctly in database
- Tests apostrophes in names (O'Brien)

### TestStaffImport_DuplicateHandling
- Tests idempotency (same CSV uploaded twice)
- Verifies ON CONFLICT behavior (duplicates are skipped)
- Ensures database count doesn't double

### TestStaffImport_MixedValidInvalid
- Tests partial success with error reporting (10 valid rows + 1 invalid row)
- Validates that invalid rows are skipped while valid rows are processed
- Verifies error messages are returned for invalid rows
- Ensures only valid rows are inserted into the database

### TestStaffImport_TenantIsolation
- Tests multi-tenant isolation
- Verifies Company A cannot see Company B's staff
- Validates company_id filtering works correctly

### TestStaffImport_UnauthorizedAccess
- Tests that requests without authentication are rejected
- Verifies 401 Unauthorized response

## Notes

- Tests are marked with `//go:build integration` to separate them from unit tests
- Tests require a real PostgreSQL database (not SQLite) to test ON CONFLICT and transactions
- Tests clean up after themselves, but may leave data if interrupted
- Tests can run against a local API server or use direct database access

