# Complete API Reference

This document consolidates all REST API endpoints for the Entorno35 NOM-035 Compliance Platform: authentication, staff, assessments, scoring, and reports. Use it for integration, testing, or client implementation.

**Version**: 1.0  
**Last updated**: 2026-03-11

---

## Table of Contents

1. [Authentication](#authentication)
2. [Billing](#billing)
3. [Staff Management](#staff-management)
4. [Assessments](#assessments)
5. [Scoring](#scoring)
6. [Reports](#reports)

---

## Authentication

### POST /auth/login

Authenticates a company admin and returns a JWT token.

**Request**:
```json
{
  "identifier": "RFC123456789",
  "type": "COMPANY",
  "password": "string"
}
```

**Response (200 OK)**:
```json
{
  "token": "jwt_token_string"
}
```

**Error Responses**:
- `400 Bad Request`: Missing required fields
- `401 Unauthorized`: Invalid credentials
- `500 Internal Server Error`: Server error

### POST /auth/admin/login

Authenticates internal billing admins and returns a JWT token with `role=admin`.

**Request**:
```json
{
  "email": "admin@entorno35.com",
  "password": "string"
}
```

**Response (200 OK)**:
```json
{
  "token": "jwt_token_string"
}
```

**Error Responses**:
- `400 Bad Request`: Missing/invalid fields
- `401 Unauthorized`: Invalid credentials
- `503 Service Unavailable`: Admin auth not configured

---

## Billing

### POST /billing/checkout-session

Creates a Stripe checkout session for new paid registrations.

**Request**:
```json
{
  "rfc": "ABC123456789",
  "company_name": "Mi Empresa SA de CV",
  "address": "Calle 1, CDMX",
  "admin_email": "admin@empresa.com",
  "password": "string (min 8 chars)",
  "plan": "monthly or yearly",
  "employee_count": 50
}
```

**Response (200 OK)**:
```json
{
  "session_id": "cs_test_xxx",
  "checkout_url": "https://checkout.stripe.com/..."
}
```

### POST /billing/webhook

Stripe webhook endpoint used to activate/deactivate subscriptions.

### GET /billing/checkout-session/:id/verify

Verifies checkout completion and returns activation status + login identifier.

### POST /api/v1/billing/refund-request

Creates a refund request for the authenticated company.

**Request**:
```json
{
  "reason": "Texto explicando el motivo del reembolso (20-2000 caracteres)"
}
```

### GET /api/v1/billing/refund-requests

Lists refund requests for the authenticated company.

**Query Parameters**:
- `limit` (integer): Items per page (default: 50, max: 100)
- `offset` (integer): Pagination offset (default: 0)

### GET /api/v1/admin/billing/refund-requests

Lists refund requests for internal billing admins (requires admin JWT).

**Query Parameters**:
- `status` (string, optional): `requested|approved|rejected|refunded`
- `limit` (integer): Items per page (default: 50, max: 100)
- `offset` (integer): Pagination offset (default: 0)

### PATCH /api/v1/admin/billing/refund-requests/:id

Updates refund request decision (admin only).

**Request**:
```json
{
  "status": "approved|rejected|refunded",
  "resolution_note": "string (required for rejected/refunded)",
  "stripe_refund_id": "re_xxx (required when status=refunded)"
}
```

---

## Staff Management

### GET /api/v1/staff

Lists all staff members with pagination.

**Query Parameters**:
- `limit` (integer): Items per page (default: 50, max: 100)
- `offset` (integer): Pagination offset (default: 0)

**Response (200 OK)**:
```json
{
  "data": [
    {
      "id": "uuid",
      "full_name": "string",
      "email": "string",
      "curp": "string (optional)",
      "employee_id": "string (optional)",
      "demographics": {
        "department": "string",
        "shift_type": "string"
      }
    }
  ],
  "total": 100,
  "limit": 50,
  "offset": 0
}
```

### POST /api/v1/staff/import

Bulk import staff from CSV file.

**Request**: Multipart form-data with CSV file

**Response (200 OK)**:
```json
{
  "imported": 10,
  "errors": []
}
```

---

## Assessments

### POST /api/v1/assessments

Creates a new assessment. Guide type is automatically determined based on company employee count.

**Request**:
```json
{
  "staff_id": "uuid",
  "period": 2025
}
```

`period` is required and must be within the current year +/- 1.

**Response (201 Created)**:
```json
{
  "id": "uuid",
  "staff_id": "uuid",
  "company_id": "uuid",
  "period": 2025,
  "guide_type": "II or III (auto-determined)",
  "status": "pending"
}
```

**Guide Type Determination**:
- Guide II: 16-50 employees
- Guide III: >50 employees
- Default: Guide II for <16 employees

### GET /api/v1/assessments

Lists assessments with optional filters.

**Query Parameters**:
- `staff_id` (uuid): Filter by staff
- `period` (integer): Filter by period (must be within current year +/- 1)
- `status` (string): Filter by status
- `limit` (integer): Items per page (default: 50, max: 100)
- `offset` (integer): Pagination offset (default: 0)

### GET /api/v1/assessments/public/:token

Retrieves assessment for public access via secure token.

**Response (200 OK)**:
```json
{
  "assessment": { "id": "uuid", "guide_type": "II", ... },
  "questions": [ { "id": 1, "text": "...", ... } ]
}
```

### POST /api/v1/assessments/public/:token/submit

Submits assessment responses (public endpoint).

**Request**:
```json
{
  "responses": [
    { "question_id": 1, "value": 3 },
    { "question_id": 2, "value": 2 }
  ]
}
```

---

## Scoring

### POST /api/v1/assessments/:id/calculate

Calculates the score for a completed assessment.

**Response (200 OK)**:
```json
{
  "message": "Assessment scored successfully",
  "assessment_id": "uuid"
}
```

**Business Logic**:
1. Validates assessment belongs to company
2. Fetches all responses with questions
3. Selects scoring strategy based on guide type
4. Calculates scores with polarity inversion
5. Updates assessment with results

---

## Reports

### GET /api/v1/reports/individual/:assessment_id

Retrieves detailed individual assessment report.

**Response (200 OK)**:
```json
{
  "assessment_id": "uuid",
  "period": 2025,
  "guide_type": "II",
  "staff_name": "Juan Pérez",
  "total_score": 45.5,
  "risk_level": "alto",
  "category_scores": { "Category Name": 25.0 },
  "category_max_scores": { "Category Name": 40.0 },
  "category_risk_levels": { "Category Name": "bajo" },
  "domain_scores": { "Domain Name": 12.5 },
  "domain_max_scores": { "Domain Name": 24.0 },
  "domain_risk_levels": { "Domain Name": "medio" },
  "recommendations": [ "..." ]
}
```

**Important**: Scores are RECALCULATED from responses to ensure accuracy. Stored CalculatedScore values are not used to prevent issues with old data.

### GET /api/v1/reports/individual/:assessment_id/pdf

Downloads PDF report for an assessment.

**Response**: PDF file with complete assessment data including:
- Score/Max ratios for all categories and domains
- Color-coded risk levels
- Professional NOM-035 formatting

### GET /api/v1/reports/general

Retrieves company-wide general report with aggregations.

**Query Parameters**:
- `period` (integer, optional): Filter by period

**Response (200 OK)**:
```json
{
  "company_id": "uuid",
  "total_staff": 150,
  "completed_assessments": 120,
  "participation_rate": 80.0,
  "risk_distribution": [
    { "risk_level": "alto", "count": 15 }
  ],
  "department_heatmap": [
    { "department": "Producción", "risk_level": "alto", "count": 8 }
  ]
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid credentials or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |

---

## Authentication

All endpoints except `/auth/login`, `/auth/admin/login`, `/billing/*`, and `/assessments/public/*` require JWT authentication:

```
Authorization: Bearer <jwt_token>
```

## Rate Limiting

Checkout session creation endpoints use configurable fixed-window rate limiting:

- `POST /billing/checkout-session`
- `POST /api/v1/billing/checkout-session`

## Changelog

- 2026-01-12: Added note about score recalculation in reports
- 2026-01-12: Fixed staff login error handling (401 instead of 500)
- 2026-01-12: Consolidated from individual endpoint documentation files
