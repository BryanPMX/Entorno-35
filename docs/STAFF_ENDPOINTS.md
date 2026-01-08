# Staff Management Endpoints

## POST /api/v1/staff

Creates a new staff member for the authenticated company.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Staff must belong to authenticated company

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Body:**
```json
{
  "curp": "ABCD123456HIJKLM01",
  "full_name": "Juan Pérez García",
  "email": "juan.perez@example.com",
  "demographics": {
    "gender": "Masculino",
    "department": "Producción",
    "role": "Operador",
    "shift_type": "Diurno"
  }
}
```

**Fields:**
- `curp` (string, required) - CURP must be exactly 18 characters
- `full_name` (string, required) - Full name of staff member
- `email` (string, optional) - Email address
- `demographics` (object, optional) - Demographic information

### Response

**Success (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "770e8400-e29b-41d4-a716-446655440000",
  "curp": "ABCD123456HIJKLM01",
  "full_name": "Juan Pérez García",
  "email": "juan.perez@example.com",
  "demographics": {
    "gender": "Masculino",
    "department": "Producción",
    "role": "Operador",
    "shift_type": "Diurno"
  },
  "created_at": "2025-12-29T10:00:00Z"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "CURP must be exactly 18 characters, got 17"
}
```

---

## GET /api/v1/staff

Lists staff members for the authenticated company with pagination.

### Authentication

- **Required**: Yes (JWT Bearer token)

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters (all optional):**
- `limit` (integer, default: 50, max: 100) - Number of records per page
- `offset` (integer, default: 0) - Number of records to skip

### Response

**Success (200 OK):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "company_id": "770e8400-e29b-41d4-a716-446655440000",
      "curp": "ABCD123456HIJKLM01",
      "full_name": "Juan Pérez García",
      "email": "juan.perez@example.com",
      "demographics": {...},
      "created_at": "2025-12-29T10:00:00Z"
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

---

## GET /api/v1/staff/:id

Retrieves a staff member by ID.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Staff must belong to authenticated company

### Request

**URL Parameters:**
- `id` (UUID) - Staff member ID

**Headers:**
```
Authorization: Bearer <jwt_token>
```

### Response

**Success (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "770e8400-e29b-41d4-a716-446655440000",
  "curp": "ABCD123456HIJKLM01",
  "full_name": "Juan Pérez García",
  "email": "juan.perez@example.com",
  "demographics": {...},
  "created_at": "2025-12-29T10:00:00Z"
}
```

**Error (404 Not Found):**
```json
{
  "error": "staff not found"
}
```

---

## PUT /api/v1/staff/:id

Updates an existing staff member.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Staff must belong to authenticated company

### Request

**URL Parameters:**
- `id` (UUID) - Staff member ID

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Body:**
```json
{
  "full_name": "Juan Pérez García Updated",
  "email": "juan.perez.new@example.com",
  "demographics": {
    "department": "Calidad",
    "role": "Supervisor"
  }
}
```

**Fields (all optional):**
- `full_name` (string) - Updated full name
- `email` (string) - Updated email
- `demographics` (object) - Updated demographic information (merged with existing)

### Response

**Success (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "770e8400-e29b-41d4-a716-446655440000",
  "curp": "ABCD123456HIJKLM01",
  "full_name": "Juan Pérez García Updated",
  "email": "juan.perez.new@example.com",
  "demographics": {...},
  "updated_at": "2025-12-29T11:00:00Z"
}
```

---

## DELETE /api/v1/staff/:id

Soft deletes a staff member.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Staff must belong to authenticated company

### Request

**URL Parameters:**
- `id` (UUID) - Staff member ID

**Headers:**
```
Authorization: Bearer <jwt_token>
```

### Response

**Success (200 OK):**
```json
{
  "message": "staff deleted successfully"
}
```

**Error (404 Not Found):**
```json
{
  "error": "staff not found"
}
```

---

## POST /api/v1/staff/import

Bulk imports staff members from a CSV file.

### Authentication

- **Required**: Yes (JWT Bearer token)

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file` (file, required) - CSV file with staff data

**CSV Format:**
```csv
Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez García,ABCD123456HIJKLM01,juan.perez@example.com,Producción,Operador,Diurno,Masculino
María González López,EFGH567890MNOPQR02,maria.gonzalez@example.com,Recursos Humanos,Analista,Diurno,Femenino
```

**CSV Columns:**
- `Name` (required) - Full name
- `CURP` (required) - Must be exactly 18 characters
- `Email` (optional) - Email address
- `Area` (optional) - Maps to `demographics.department`
- `Job` (optional) - Maps to `demographics.role`
- `Shift` (optional) - Maps to `demographics.shift_type`
- `Gender` (optional) - Maps to `demographics.gender`

### Response

**Success (200 OK):**
```json
{
  "total_processed": 10,
  "success_count": 8,
  "skipped_count": 2,
  "errors": [
    "Row 3: CURP must be exactly 18 characters, got 'ABCD123456HIJKLM0' (17)",
    "Row 7: Name is required"
  ]
}
```

**Error (400 Bad Request):**
```json
{
  "error": "missing required header: curp"
}
```

### Business Logic

1. **CSV Parsing**: Reads CSV file with case-insensitive header matching
2. **Validation**: 
   - Validates CURP format (18 characters)
   - Validates required fields (Name, CURP)
3. **Demographics Mapping**:
   - `Area` → `demographics.department`
   - `Job` → `demographics.role`
   - `Shift` → `demographics.shift_type`
   - `Gender` → `demographics.gender`
4. **Bulk Insert**: Uses transaction with `ON CONFLICT DO NOTHING` for CURP collisions
5. **Result**: Returns count of successfully inserted records and any errors

### Notes

- Duplicate CURPs (same company) are skipped without error (counted in `skipped_count`)
- Invalid rows are reported in `errors` array but don't stop the import
- All valid rows are processed in a single transaction
- See `docs/staff_import_template.csv` for a sample template

---

## Example: Create Staff

```bash
curl -X POST http://localhost:8080/api/v1/staff \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "curp": "ABCD123456HIJKLM01",
    "full_name": "Juan Pérez García",
    "email": "juan.perez@example.com",
    "demographics": {
      "gender": "Masculino",
      "department": "Producción",
      "role": "Operador",
      "shift_type": "Diurno"
    }
  }'
```

## Example: Import from CSV

```bash
curl -X POST http://localhost:8080/api/v1/staff/import \
  -H "Authorization: Bearer <token>" \
  -F "file=@staff_import.csv"
```


---

## POST /api/v1/staff/csv/analyze

Analyzes a CSV file and provides intelligent column mapping suggestions using fuzzy matching and data type inference.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Company-level access

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file` (file, required) - CSV file to analyze

### Response

**Success (200 OK):**
```json
{
  "mappings": [
    {
      "csv_header": "Full Name",
      "expected_field": "name",
      "confidence": 0.95,
      "match_type": "synonym",
      "inferred_type": "text",
      "sample_values": ["Juan Pérez", "María García", "Carlos López"]
    },
    {
      "csv_header": "Correo",
      "expected_field": "email",
      "confidence": 0.88,
      "match_type": "fuzzy",
      "inferred_type": "email",
      "sample_values": ["juan@example.com", "maria@example.com"]
    }
  ],
  "unmapped_headers": ["Address", "Phone"],
  "overall_confidence": 0.85,
  "warnings": ["Low confidence mapping for 'Correo' -> 'email'"],
  "requires_review": false
}
```

---

## POST /api/v1/staff/csv/preview

Previews a CSV import with the specified column mappings and validates data.

### Authentication

- **Required**: Yes (JWT Bearer token)
- **Authorization**: Company-level access

### Request

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file` (file, required) - CSV file to preview
- `mappings` (string, required) - JSON array of column mappings

### Response

**Success (200 OK):**
```json
{
  "total_rows": 5,
  "valid_rows": 4,
  "invalid_rows": 1,
  "sample_records": [
    {
      "row_number": 1,
      "data": {
        "name": "Juan Pérez",
        "email": "juan@example.com",
        "area": "IT"
      },
      "errors": []
    },
    {
      "row_number": 2,
      "data": {
        "name": "",
        "email": "invalid-email",
        "area": "HR"
      },
      "errors": ["name is required", "invalid email format"]
    }
  ],
  "errors": []
}
```

---

## CSV Column Mapping Features

### Intelligent Detection
- **Exact Match**: Perfect header matches
- **Synonym Recognition**: Common variations (Name ↔ Full Name, Email ↔ Correo)
- **Fuzzy Matching**: Levenshtein distance similarity scoring
- **Data Type Inference**: Automatic detection of emails, numbers, enums

### Supported Field Mappings

| Expected Field | Type | Synonyms |
|----------------|------|----------|
| name | text | Full Name, Employee Name, Nombre, Nombre Completo |
| curp | text | RFC, ID, Employee ID, Identificación |
| email | email | Correo, Email Address, Mail |
| area | text | Department, Área, Departamento, Sector |
| job | text | Position, Role, Cargo, Puesto, Ocupación |
| shift | enum | Turno, Jornada, Schedule, Horario |
| gender | enum | Sexo, Género, Sex |

### Confidence Scoring
- **1.0**: Exact match or perfect synonym
- **0.8-0.9**: High fuzzy match confidence
- **0.6-0.8**: Medium confidence, may need review
- **<0.6**: Low confidence, requires manual mapping

### Data Type Inference
- **Text**: Names, departments, positions
- **Email**: RFC 5322 compliant email patterns
- **Number**: CURP (18 chars), employee IDs
- **Enum**: Limited values (gender, shift types)
