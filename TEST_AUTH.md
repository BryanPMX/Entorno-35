# Testing Authentication Endpoints

## Quick Test Guide

### 1. Start Services

```bash
# Start Docker services
make docker-up

# Start API server (in another terminal)
go run cmd/api/main.go
```

### 2. Setup Test Data

**Option A: Via SQL (Quick Test)**

```sql
-- Connect to database
psql -h localhost -U entorno35 -d entorno35

-- Create test company
INSERT INTO companies (rfc, name, subscription_status, employee_count)
VALUES ('TEST123456789', 'Test Company', 'active', 10);

-- Get company ID (copy this UUID)
SELECT id FROM companies WHERE rfc = 'TEST123456789';

-- Create test staff (replace COMPANY_ID with UUID from above)
INSERT INTO staff (company_id, curp, full_name, email)
VALUES ('COMPANY_ID_HERE', 'TESTCURP12345678901', 'Test Staff', 'test@example.com');
```

**Option B: Via API** (requires company/staff creation endpoints - future work)

### 3. Test Company Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TEST123456789",
    "type": "COMPANY"
  }'
```

**Expected Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. Test Staff Login

```bash
# Replace COMPANY_ID with the UUID from step 2
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TESTCURP12345678901",
    "type": "STAFF",
    "company_id": "COMPANY_ID_HERE"
  }'
```

**Expected Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 5. Test Error Cases

**Invalid Company:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "INVALID123",
    "type": "COMPANY"
  }'
```

**Inactive Company:**
```sql
-- Set company to inactive
UPDATE companies SET subscription_status = 'inactive' WHERE rfc = 'TEST123456789';
```

Then try login again - should return 401.

**Missing company_id for STAFF:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TESTCURP12345678901",
    "type": "STAFF"
  }'
```

### 6. Use Token in Requests

```bash
# Save token from login response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Use in subsequent requests (when protected endpoints are added)
curl -X GET http://localhost:8080/api/protected \
  -H "Authorization: Bearer $TOKEN"
```

## Environment Variables Required

```bash
# Database
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=entorno35
export DB_PASSWORD=entorno35  # or your password
export DB_NAME=entorno35

# JWT
export JWT_SECRET=your-secret-key-here  # Use a strong secret in production
export JWT_EXPIRY=24h

# Server
export PORT=8080
```

## Complete Test Script

```bash
#!/bin/bash

# Test Company Login
echo "Testing Company Login..."
COMPANY_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TEST123456789",
    "type": "COMPANY"
  }')

echo "Response: $COMPANY_RESPONSE"
TOKEN=$(echo $COMPANY_RESPONSE | jq -r '.token')
echo "Token: ${TOKEN:0:50}..."

# Test Staff Login (requires COMPANY_ID from database)
echo ""
echo "Testing Staff Login..."
STAFF_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"identifier\": \"TESTCURP12345678901\",
    \"type\": \"STAFF\",
    \"company_id\": \"COMPANY_ID_HERE\"
  }")

echo "Response: $STAFF_RESPONSE"
```

