# Report Endpoints

## GET /api/v1/reports/individual/:assessment_id/pdf

Generates and downloads a PDF report for an individual assessment with professional NOM-035 formatting.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Assessment must belong to authenticated company

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `assessment_id` (UUID, required) - Assessment ID

### Response

**Success (200 OK):**
- **Content-Type**: `application/pdf`
- **Content-Disposition**: `attachment; filename="NOM035_Report_{staff_name}_{date}.pdf"`
- **Body**: PDF file data

### Error Responses

**404 Not Found:**
```json
{
  "error": "assessment report not found"
}
```

**400 Bad Request:**
```json
{
  "error": "assessment has not been scored yet"
}
```

**401 Unauthorized:**
```json
{
  "error": "unauthorized"
}
```

### PDF Content Structure

The generated PDF includes:
- **Header**: NOM-035 Psychosocial Risk Assessment Report
- **Assessment Information**: Staff details, department, period, completion date
- **Risk Assessment**: Total score, risk level, medical attention warnings
- **Category Scores**: NOM-035 category breakdown (Guide II/III)
- **Domain Scores**: Detailed psychosocial risk domains
- **Recommendations**: Actionable NOM-035 Section 8 recommendations
- **Footer**: Generation timestamp and platform information

### Notes

- PDFs are generated server-side for consistent formatting and security
- Only available for completed, scored assessments
- Filename automatically generated with staff name and date
- Professional A4 layout suitable for HR documentation and compliance records

---

## GET /api/v1/reports/individual/:assessment_id

Retrieves an individual assessment report with detailed scoring breakdown and NOM-035 recommendations.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Assessment must belong to authenticated company

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `assessment_id` (UUID, required) - Assessment ID

### Response

**Success (200 OK):**
```json
{
  "assessment_id": "660e8400-e29b-41d4-a716-446655440000",
  "period": 2025,
  "guide_type": "II",
  "staff_name": "Juan Pérez García",
  "department": "Producción",
  "shift": "Diurno",
  "total_score": 45.5,
  "risk_level": "alto",
  "category_scores": {
    "Factores de riesgo psicosocial": 25.0,
    "Entorno organizacional": 20.5
  },
  "domain_scores": {
    "Carga de trabajo": 12.5,
    "Falta de control sobre el trabajo": 8.0,
    "Jornada de trabajo": 5.0,
    "Liderazgo": 10.0,
    "Relaciones en el trabajo": 10.0
  },
  "requires_medical_attention": false,
  "completed_at": "2025-12-30T10:00:00Z",
  "recommendations": [
    "Revisar y redistribuir la carga de trabajo. Implementar pausas activas y rotación de tareas.",
    "Capacitar a supervisores en liderazgo positivo. Establecer retroalimentación constructiva.",
    "NIVEL DE RIESGO ALTO: Se recomienda implementar medidas preventivas y correctivas. Realizar seguimiento periódico."
  ]
}
```

**Fields:**
- `assessment_id` (UUID) - Assessment identifier
- `period` (integer) - Assessment period (e.g., 2025)
- `guide_type` (string) - Guide type (I, II, or III)
- `staff_name` (string) - Full name of staff member
- `department` (string, optional) - Department from demographics
- `shift` (string, optional) - Shift type from demographics
- `total_score` (float) - Total assessment score
- `risk_level` (string) - Risk level: "nulo", "bajo", "medio", "alto", "muy_alto"
- `category_scores` (object) - Scores by category (Guide II/III only)
- `domain_scores` (object) - Scores by domain (Guide II/III only)
- `requires_medical_attention` (boolean) - Medical attention flag (Guide I only)
- `completed_at` (string, optional) - ISO 8601 timestamp when assessment was completed
- `recommendations` (array, optional) - NOM-035 Section 8 recommendations (only for "alto" or "muy_alto" risk levels)

### Error Responses

**404 Not Found:**
```json
{
  "error": "assessment report not found"
}
```

**400 Bad Request:**
```json
{
  "error": "assessment has not been scored yet"
}
```

**401 Unauthorized:**
```json
{
  "error": "unauthorized"
}
```

### Notes

- Reports are only available for assessments that have been scored
- Recommendations are automatically generated based on:
  - Risk level (only for "alto" and "muy_alto")
  - Top 3 highest-scoring domains
- Category and domain scores are recalculated from responses for accurate reporting
- Guide I (Trauma) assessments will have empty category_scores and domain_scores

---

## GET /api/v1/reports/general

Retrieves a company-wide general report with aggregated metrics, risk distribution, and department heatmap.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Must be authenticated company

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `period` (integer, optional) - Filter by assessment period (e.g., 2025)

### Response

**Success (200 OK):**
```json
{
  "company_id": "770e8400-e29b-41d4-a716-446655440000",
  "company_name": "Acme Corporation",
  "period": 2025,
  "total_staff": 150,
  "completed_assessments": 120,
  "participation_rate": 80.0,
  "risk_distribution": [
    {
      "risk_level": "nulo",
      "count": 20
    },
    {
      "risk_level": "bajo",
      "count": 45
    },
    {
      "risk_level": "medio",
      "count": 35
    },
    {
      "risk_level": "alto",
      "count": 15
    },
    {
      "risk_level": "muy_alto",
      "count": 5
    }
  ],
  "department_heatmap": [
    {
      "department": "Producción",
      "risk_level": "alto",
      "count": 8
    },
    {
      "department": "Producción",
      "risk_level": "medio",
      "count": 12
    },
    {
      "department": "Recursos Humanos",
      "risk_level": "bajo",
      "count": 5
    },
    {
      "department": "Recursos Humanos",
      "risk_level": "medio",
      "count": 3
    }
  ]
}
```

**Fields:**
- `company_id` (UUID) - Company identifier
- `company_name` (string) - Company name
- `period` (integer, optional) - Assessment period if filtered
- `total_staff` (integer) - Total number of staff members in company
- `completed_assessments` (integer) - Number of completed assessments
- `participation_rate` (float) - Participation percentage (completed / total * 100)
- `risk_distribution` (array) - Count of assessments by risk level
  - `risk_level` (string) - Risk level: "nulo", "bajo", "medio", "alto", "muy_alto"
  - `count` (integer) - Number of assessments with this risk level
- `department_heatmap` (array) - Risk distribution by department
  - `department` (string) - Department name from staff demographics
  - `risk_level` (string) - Risk level
  - `count` (integer) - Number of assessments in this department/risk combination

### Error Responses

**404 Not Found:**
```json
{
  "error": "company not found"
}
```

**401 Unauthorized:**
```json
{
  "error": "unauthorized"
}
```

### Notes

- Only completed assessments are included in aggregations
- Participation rate is calculated as: `(completed_assessments / total_staff) * 100`
- Department heatmap only includes staff with department information in demographics
- Risk distribution includes all risk levels, even if count is 0
- If `period` is not provided, all completed assessments are included regardless of period
- Department heatmap groups by both department and risk level for cross-tabulation analysis

### Use Cases

- **Compliance Reporting**: Generate official NOM-035 reports for STPS submission
- **Risk Analysis**: Identify departments with high risk concentrations
- **Participation Tracking**: Monitor assessment completion rates
- **Trend Analysis**: Compare metrics across different periods

---

## Report Data Structure

### Individual Report Use Cases

- **PDF Generation**: Frontend can use this data to generate official NOM-035 individual reports
- **Staff Review**: Staff members can view their own assessment results
- **HR Analysis**: HR can review individual assessments for intervention planning

### General Report Use Cases

- **Executive Dashboard**: High-level metrics for management
- **Compliance Documentation**: Official NOM-035 general report data
- **Department Analysis**: Identify departments requiring intervention
- **Risk Monitoring**: Track risk distribution across the organization

---

**Last Updated**: 2026-01-09
**API Version**: v1

