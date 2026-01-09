# Manual E2E Testing Guide - Golden Path

This guide walks through manual testing of the complete Entorno35 workflow.

## Prerequisites

1. **Start Infrastructure** (Root directory):
   ```bash
   make docker-up
   make migrate-up
   ```

2. **Start Services** (Two terminals):
   ```bash
   # Terminal 1: Backend
   export DB_URL="postgres://entorno35:entorno35@localhost:5432/entorno35?sslmode=disable"
   export JWT_SECRET="test-secret"
   export CORS_ORIGIN="http://localhost:3000"
   make run

   # Terminal 2: Frontend
   cd web/frontend && npm run dev
   ```

## Golden Path Test Steps

### 1. Admin Registration/Login

**Via Frontend:**
- Open http://localhost:3000
- Click "Company Login"
- Enter test credentials (any RFC will work for demo)

**Via API:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TEST123456789",
    "type": "COMPANY"
  }'
```

### 2. Staff Import

**Via Frontend:**
- Navigate to Staff Management
- Upload CSV file with staff data

**Sample CSV:**
```csv
Nombre Completo,CURP,Correo Electrónico,Departamento,Género,Turno
Juan Pérez García,PEGJ900101HDFRRR01,juan.perez@test.com,Sistemas,Masculino,Diurno
María González López,GOLM901202MDFNNN02,maria.gonzalez@test.com,Recursos Humanos,Femenino,Diurno
```

**Via API:**
```bash
# Get admin token first, then:
curl -X POST http://localhost:8080/api/v1/staff/import \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@staff.csv"
```

### 3. Create Assessment

**Via Frontend:**
- Go to Assessments page
- Click "Create Assessment"
- Select staff member and period (2025)

**Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/assessments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "staff_id": "STAFF_ID_HERE",
    "period": 2025
  }'
```

### 4. Generate Assessment Link

**Via Frontend:**
- Click assessment in list
- Click "Generate Link"

**Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/assessments/ASSESSMENT_ID/links \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"expires_in_days": 7}'
```

### 5. Staff Takes Assessment

**Via Browser (Mobile Simulation):**
- Open assessment URL: `http://localhost:8080/assessment/TOKEN_HERE`
- Answer questions using Likert scale (0-4)
- Submit assessment

**Via API:**
```bash
# Get questions
curl http://localhost:8080/api/v1/assessments/public/TOKEN

# Submit responses
curl -X POST http://localhost:8080/api/v1/assessments/public/TOKEN/submit \
  -H "Content-Type: application/json" \
  -d '{
    "responses": [
      {"question_id": 1, "value": 0},
      {"question_id": 2, "value": 1},
      {"question_id": 3, "value": 2}
    ]
  }'
```

### 6. Admin Views Reports

**Via Frontend:**
- Go to Reports section
- View individual assessment reports
- Check company-wide analytics

**Via API:**
```bash
# Individual report
curl http://localhost:8080/api/v1/reports/individual/ASSESSMENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"

# Company report
curl http://localhost:8080/api/v1/reports/general \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Expected Results

- ✅ Assessment creation successful
- ✅ Secure link generation
- ✅ Staff can access assessment
- ✅ Questions load properly
- ✅ Responses save and score automatically
- ✅ Risk level calculation (Nulo, Bajo, Medio, Alto, Muy Alto)
- ✅ Individual reports with recommendations
- ✅ Company analytics with participation rates
- ✅ Department risk heatmaps

## Cleanup

```bash
# Stop services
make docker-down

# Kill background processes if needed
pkill -f "go run cmd/api/main.go"
pkill -f "npm run dev"
```

## Troubleshooting

- **Backend won't start**: Check DB_URL environment variable
- **Frontend errors**: Ensure CORS_ORIGIN is set correctly
- **Assessment not found**: Verify token and expiration
- **Reports empty**: Wait for scoring to complete (2-3 seconds)

## Alternative: Automated E2E

For automated testing, run:
```bash
./e2e_walkthrough.sh
```