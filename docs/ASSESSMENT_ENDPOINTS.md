# Assessment Endpoints

## POST /api/v1/assessments

Creates a new assessment for a staff member. Guide type is automatically determined based on company employee count.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Staff must belong to authenticated company

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Body:**
```json
{
  "staff_id": "550e8400-e29b-41d4-a716-446655440000",
  "period": 2025
}
```

**Fields:**
- `staff_id` (UUID, required) - Staff member ID
- `period` (integer, required) - Assessment period (e.g., 2025)

### Response

**Success (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "staff_id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "770e8400-e29b-41d4-a716-446655440000",
  "period": 2025,
  "guide_type": "II",
  "status": "pending",
  "created_at": "2025-12-29T10:00:00Z"
}
```

### Guide Type Determination

- **Guide II**: 16-50 employees (or <16 employees default)
- **Guide III**: >50 employees

### Example

```bash
curl -X POST http://localhost:8080/api/v1/assessments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "staff_id": "550e8400-e29b-41d4-a716-446655440000",
    "period": 2025
  }'
```

---

## GET /api/v1/assessments

Lists assessments for the authenticated company with optional filters.

### Authentication

- **Required**: Yes (JWT Bearer token)

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters (all optional):**
- `staff_id` (UUID) - Filter by staff member
- `period` (integer) - Filter by period
- `status` (string) - Filter by status: `pending`, `completed`, `cancelled`

### Response

**Success (200 OK):**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "staff_id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "770e8400-e29b-41d4-a716-446655440000",
    "period": 2025,
    "guide_type": "II",
    "status": "pending",
    "created_at": "2025-12-29T10:00:00Z"
  }
]
```

### Example

```bash
# List all assessments
curl -X GET http://localhost:8080/api/v1/assessments \
  -H "Authorization: Bearer <token>"

# Filter by staff
curl -X GET "http://localhost:8080/api/v1/assessments?staff_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer <token>"

# Filter by period and status
curl -X GET "http://localhost:8080/api/v1/assessments?period=2025&status=completed" \
  -H "Authorization: Bearer <token>"
```

---

## GET /api/v1/assessments/:id

Retrieves a specific assessment by ID.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Assessment must belong to authenticated company

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
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "staff_id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "770e8400-e29b-41d4-a716-446655440000",
  "period": 2025,
  "guide_type": "II",
  "status": "pending",
  "total_score": null,
  "risk_level": null,
  "requires_medical_attention": false,
  "completed_at": null,
  "created_at": "2025-12-29T10:00:00Z",
  "updated_at": "2025-12-29T10:00:00Z"
}
```

**Error (404 Not Found):**
```json
{
  "error": "assessment not found"
}
```

### Example

```bash
curl -X GET http://localhost:8080/api/v1/assessments/660e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

---

## POST /api/v1/assessments/:id/links

Creates a secure link (token) for staff to access an assessment without authentication.

### Authentication

- **Required**: Yes (JWT Bearer token with staff ID)
- **Authorization**: Assessment must belong to authenticated staff member

### Request

**URL Parameters:**
- `id` (UUID) - Assessment ID

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Body:**
```json
{
  "expires_in_days": 7
}
```

**Fields:**
- `expires_in_days` (integer, required) - Link expiration in days (1-365)

### Response

**Success (201 Created):**
```json
{
  "token": "880e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2026-01-05T10:00:00Z",
  "link": ""
}
```

**Note**: The `link` field is empty - the frontend should construct the full URL using the token.

### Example

```bash
curl -X POST http://localhost:8080/api/v1/assessments/660e8400-e29b-41d4-a716-446655440000/links \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "expires_in_days": 7
  }'
```

### Business Logic

- Link token is a UUID
- Link expires after the specified number of days
- Link can be used by staff to access the assessment without authentication (public endpoint implemented)
- Link access is tracked via `accessed_at` timestamp

---

## GET /api/v1/assessments/public/:token

Retrieves assessment details and questions for public access via secure token.

### Authentication

- **Required**: No (public endpoint using token-based access)

### Request

**URL Parameters:**
- `token` (string, required) - Secure assessment link token

### Response

**Success (200 OK):**
```json
{
  "assessment": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "staff_id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "770e8400-e29b-41d4-a716-446655440000",
    "period": 2025,
    "guide_type": "II",
    "status": "pending",
    "requires_medical_attention": false,
    "created_at": "2025-12-29T10:00:00Z",
    "updated_at": "2025-12-29T10:00:00Z"
  },
  "questions": [
    {
      "id": 1,
      "question_number": 1,
      "guide_type": "II",
      "type": "likert",
      "text": "¿Con qué frecuencia te sientes agotado al final de la jornada laboral?",
      "polarity": "POSITIVE",
      "order_index": 0,
      "created_at": "2025-12-29T10:00:00Z",
      "updated_at": "2025-12-29T10:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Invalid or expired token, assessment not in pending status
- `500 Internal Server Error`: Server error

### Example

```bash
curl -X GET http://localhost:8080/api/v1/assessments/public/880e8400-e29b-41d4-a716-446655440000
```

### Business Logic

- Validates token exists and has not expired
- Ensures assessment is in "pending" status
- Returns assessment metadata and ordered question list
- Questions are ordered by `order_index` for proper flow

---

## POST /api/v1/assessments/public/:token/submit

Submits staff responses for a public assessment.

### Authentication

- **Required**: No (public endpoint using token-based access)

### Request

**URL Parameters:**
- `token` (string, required) - Secure assessment link token

**Body:**
```json
{
  "responses": [
    {
      "question_id": 1,
      "value": 3
    },
    {
      "question_id": 2,
      "value": 2
    }
  ]
}
```

**Fields:**
- `responses` (array, required) - Array of response objects
  - `question_id` (integer, required) - Question ID
  - `value` (integer, required) - Response value (0-4 for Likert scale)

### Response

**Success (200 OK):**
```json
{
  "message": "Assessment submitted successfully",
  "submitted_at": "2025-12-29T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid token, expired link, assessment not pending, etc.
- `500 Internal Server Error`: Server error

### Example

```bash
curl -X POST http://localhost:8080/api/v1/assessments/public/880e8400-e29b-41d4-a716-446655440000/submit \
  -H "Content-Type: application/json" \
  -d '{
    "responses": [
      {"question_id": 1, "value": 3},
      {"question_id": 2, "value": 2}
    ]
  }'
```

### Business Logic

- Validates token and assessment status
- Saves responses in transaction
- Marks assessment link as used (`accessed_at`)
- Triggers automatic scoring calculation
- Updates assessment status to "completed"

