# Response Submission Endpoints

## POST /api/v1/assessments/public/:token/submit

Public endpoint for staff to submit their assessment responses. This endpoint does not require authentication - it uses a secure token from the assessment link.

### Authentication

- **Required**: No (public endpoint)
- **Authorization**: Token-based (assessment link token)

### Request

**URL Parameters:**
- `token` (UUID string) - Assessment link token (from `POST /api/v1/assessments/:id/links`)

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "responses": [
    {
      "question_id": 1,
      "value": 4
    },
    {
      "question_id": 2,
      "value": 3
    }
  ]
}
```

**Fields:**
- `responses` (array, required) - Array of responses
  - `question_id` (integer, required) - Question ID
  - `value` (integer, required) - Selected value (0-4 for Likert, 0-1 for Binary)

### Response

**Success (200 OK):**
```json
{
  "message": "Assessment submitted successfully",
  "submitted_at": "2025-12-29T10:30:00Z"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "assessment link has expired"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "assessment link has already been used"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "assessment is not in pending status (current status: completed)"
}
```

**Error (500 Internal Server Error):**
```json
{
  "error": "failed to save responses: <details>"
}
```

### Business Logic

1. **Link Validation**: Fetches assessment link by token
   - Verifies link exists
   - Checks link has not expired
   - Verifies link has not been used (accessed_at is null)

2. **Assessment Validation**: 
   - Verifies assessment exists
   - Checks assessment is in "pending" status

3. **Save Responses**: Saves all responses in a database transaction
   - If any response fails, entire transaction is rolled back
   - Ensures AssessmentID is set for all responses
   - Validates SelectedValue is in range (0-4)

4. **Mark Link as Used**: Updates link's `accessed_at` timestamp

5. **Trigger Scoring**: Automatically calculates assessment scores
   - Uses ScoringService to calculate scores
   - Updates assessment with total score, risk level, category/domain scores
   - Sets assessment status to "completed"
   - Sets `completed_at` timestamp

### Example

```bash
# Submit responses using assessment link token
curl -X POST http://localhost:8080/api/v1/assessments/public/550e8400-e29b-41d4-a716-446655440000/submit \
  -H "Content-Type: application/json" \
  -d '{
    "responses": [
      {"question_id": 1, "value": 4},
      {"question_id": 2, "value": 3},
      {"question_id": 3, "value": 2}
    ]
  }'
```

### Notes

- This endpoint is **public** (no authentication required)
- Token must be valid and not expired
- Token can only be used once (marked as used after first submission)
- Responses are saved atomically (all or nothing)
- Scoring is triggered automatically after successful submission
- Assessment status is automatically updated to "completed"

