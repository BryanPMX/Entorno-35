# Authentication Endpoints

## POST /auth/login

Authenticates a user (company or staff) and returns a JWT token.

### Request Body

```json
{
  "identifier": "string",  // RFC for COMPANY, CURP for STAFF
  "type": "COMPANY" | "STAFF",
  "company_id": "uuid",    // Required for STAFF type
  "password": "string"     // Required for STAFF type
}
```

### Request Examples

**Company Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "ABC123456789",
    "type": "COMPANY"
  }'
```

**Staff Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "CURP12345678901234",
    "type": "STAFF",
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "password": "staffpassword123"
  }'
```

### Response

**Success (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error (400 Bad Request):**
```json
{
  "error": "validation error message"
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "company not found"
}
```
or
```json
{
  "error": "company subscription is inactive"
}
```

### Validation Rules

1. **identifier**: Required, must be valid RFC (for COMPANY) or CURP (for STAFF)
2. **type**: Required, must be "COMPANY" or "STAFF"
3. **company_id**: Required when type is "STAFF"
4. **password**: Required when type is "STAFF"

### Business Rules

1. **Company Authentication**:
   - Company must exist (by RFC)
   - Company subscription must be active
   - Returns token with company_id (no staff_id)

2. **Staff Authentication**:
   - Company must exist and be active
   - Staff must exist (by CURP) within the specified company
   - Password must be provided and match the stored hash
   - Returns token with both company_id and staff_id

### Security Notes

Staff authentication now uses bcrypt password hashing for secure credential verification.

### Testing with Seeded Data

After seeding the database, you can test login with:

1. **Create a test company** (via SQL or seeder extension):
```sql
INSERT INTO companies (rfc, name, subscription_status, employee_count)
VALUES ('TEST123456789', 'Test Company', 'active', 10);
```

2. **Create test staff**:
```sql
INSERT INTO staff (company_id, curp, full_name, email)
SELECT id, 'TESTCURP12345678901', 'Test Staff', 'test@example.com'
FROM companies WHERE rfc = 'TEST123456789';
```

3. **Test Company Login**:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TEST123456789",
    "type": "COMPANY"
  }'
```

4. **Test Staff Login** (use company ID from previous response or database):
```bash
# First get company ID from database or previous login token
COMPANY_ID=$(psql -d entorno35 -t -c "SELECT id FROM companies WHERE rfc='TEST123456789';" | xargs)

curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"identifier\": \"TESTCURP12345678901\",
    \"type\": \"STAFF\",
    \"company_id\": \"$COMPANY_ID\"
  }"
```

### Using the Token

After successful login, include the token in subsequent requests:

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/api/protected-endpoint \
  -H "Authorization: Bearer $TOKEN"
```

