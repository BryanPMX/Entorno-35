# Scoring Endpoints

## POST /api/v1/assessments/:id/calculate

Calculates the score for a completed assessment and updates the assessment record with results.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Assessment must belong to the authenticated company

### Request

**URL Parameters:**
- `id` (UUID) - Assessment ID

**Headers:**
```
Authorization: Bearer <jwt_token>
```

### Response

**Success (200 OK):**
```json
{
  "message": "Assessment scored successfully",
  "assessment_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "invalid assessment ID"
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "authorization header required"
}
```

**Error (403 Forbidden):**
```json
{
  "error": "access denied: company ID mismatch"
}
```

**Error (404 Not Found):**
```json
{
  "error": "assessment not found"
}
```

**Error (500 Internal Server Error):**
```json
{
  "error": "failed to calculate assessment score: <details>"
}
```

### Example

```bash
# Calculate score for an assessment
curl -X POST http://localhost:8080/api/v1/assessments/550e8400-e29b-41d4-a716-446655440000/calculate \
  -H "Authorization: Bearer <token>"
```

### Business Logic

1. Validates assessment belongs to authenticated company
2. Fetches all responses with questions loaded
3. Selects scoring strategy based on GuideType (I, II, or III)
4. Calculates scores using the appropriate strategy:
   - **Guide I**: Binary Yes/No logic with medical attention flag
   - **Guide II/III**: Likert scale with polarity inversion and hierarchical aggregation
5. Updates assessment record with:
   - Total score
   - Risk level
   - Category scores (Guide II/III)
   - Domain scores (Guide II/III)
   - Medical attention flag (Guide I)
   - Status set to "completed"
   - Completed timestamp

### Notes

- Assessment must exist and belong to the authenticated company
- Responses must already be saved in the database
- The calculation is idempotent (can be run multiple times safely)
- Results overwrite previous calculations

