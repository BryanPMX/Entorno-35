# Testing Response Submission

## Prerequisites

1. Database must be migrated and seeded with questions
2. Create a test company and staff member
3. Create an assessment for the staff member
4. Generate an assessment link token

## Step-by-Step Test

### 1. Create Company and Staff (if not exists)

```bash
# Login as company (get JWT token)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "COMP001",
    "type": "COMPANY"
  }'

# Save the token from response
export TOKEN="<your-jwt-token>"
```

### 2. Create Assessment

```bash
# Replace STAFF_ID with actual staff UUID
curl -X POST http://localhost:8080/api/v1/assessments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "staff_id": "550e8400-e29b-41d4-a716-446655440000",
    "period": 2025
  }'

# Save the assessment ID from response
export ASSESSMENT_ID="<assessment-id-from-response>"
```

### 3. Generate Assessment Link

```bash
# Login as staff member first to get staff JWT token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "STAFF001",
    "type": "STAFF",
    "company_id": "770e8400-e29b-41d4-a716-446655440000"
  }'

# Save staff token
export STAFF_TOKEN="<staff-jwt-token>"

# Generate link
curl -X POST http://localhost:8080/api/v1/assessments/$ASSESSMENT_ID/links \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "expires_in_days": 7
  }'

# Save the token from response
export LINK_TOKEN="<token-from-link-response>"
```

### 4. Submit Responses

```bash
# Submit assessment responses (public endpoint, no auth required)
curl -X POST http://localhost:8080/api/v1/assessments/public/$LINK_TOKEN/submit \
  -H "Content-Type: application/json" \
  -d '{
    "responses": [
      {"question_id": 1, "value": 4},
      {"question_id": 2, "value": 3},
      {"question_id": 3, "value": 2},
      {"question_id": 4, "value": 1},
      {"question_id": 5, "value": 0}
    ]
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Assessment submitted successfully",
  "submitted_at": "2025-12-29T10:30:00Z"
}
```

### 5. Verify Assessment was Scored

```bash
# Check assessment status (should be "completed")
curl -X GET http://localhost:8080/api/v1/assessments/$ASSESSMENT_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response:**
```json
{
  "id": "...",
  "status": "completed",
  "total_score": 10.5,
  "risk_level": "bajo",
  "completed_at": "2025-12-29T10:30:00Z",
  ...
}
```

## Error Scenarios

### Invalid Token

```bash
curl -X POST http://localhost:8080/api/v1/assessments/public/invalid-token/submit \
  -H "Content-Type: application/json" \
  -d '{
    "responses": [{"question_id": 1, "value": 4}]
  }'
```

**Expected Response (400 Bad Request):**
```json
{
  "error": "invalid or expired assessment link: assessment link not found"
}
```

### Expired Token

If the link token has expired (past `expires_at`), you'll get:

**Expected Response (400 Bad Request):**
```json
{
  "error": "assessment link has expired"
}
```

### Already Used Token

If you try to submit with the same token twice:

**Expected Response (400 Bad Request):**
```json
{
  "error": "assessment link has already been used"
}
```

### Invalid Question Value

```bash
curl -X POST http://localhost:8080/api/v1/assessments/public/$LINK_TOKEN/submit \
  -H "Content-Type: application/json" \
  -d '{
    "responses": [
      {"question_id": 1, "value": 5}
    ]
  }'
```

**Expected Response (500 Internal Server Error):**
```json
{
  "error": "failed to save responses: selected value must be between 0 and 4, got 5 for question 1"
}
```

## Notes

- The endpoint is **public** (no authentication required)
- Token must be valid, not expired, and not previously used
- All responses are saved in a transaction (all or nothing)
- Scoring is triggered automatically after successful submission
- Assessment status is automatically updated to "completed"

