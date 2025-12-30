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
- Link can be used by staff to access the assessment without authentication (public endpoint to be implemented)
- Link access is tracked via `accessed_at` timestamp

