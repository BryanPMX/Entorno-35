# Entorno35 - Project Analysis

## Executive Summary

After analyzing the SRS document and the NOM-035 reference PDF, I've identified the domain logic, validated the proposed tech stack, and provide recommendations before implementation begins.

---

## 1. Domain Logic Analysis (Extracted from Reference PDF)

### 1.1 Assessment Types

#### Guide I - Trauma Assessment (ATS)
- **Type**: Binary (Yes/No) questions
- **Purpose**: Identify severe traumatic events requiring medical attention
- **Structure**:
  - Acontecimiento Traumático Severo (ATS) questions
  - Recuerdos persistentes (Persistent memories)
  - Esfuerzo por evitar (Effort to avoid)
  - Afectación (Affection/Symptoms)
- **Critical Logic**: If specific "Yes" answers detected → Flag as "Requires Medical Attention"

#### Guide II - Risk Assessment (16-50 employees)
- **Type**: Likert scale (5-point)
- **Scale**: Siempre, Casi Siempre, Algunas Veces, Casi Nunca, Nunca
- **Categories**:
  1. Ambiente de trabajo (1 domain, 3 dimensions, 3 questions)
  2. Factores propios de la actividad (2 domains, 9 dimensions, 20 questions)
  3. Organización del tiempo de trabajo (1 domain, 3 dimensions, 3 questions)
  4. Liderazgo y relaciones en el trabajo (3 domains, 5 dimensions, 19 questions)

#### Guide III - Risk Assessment (>50 employees)
- **Type**: Same as Guide II with additional detailed domain evaluation
- **Additional**: More granular domain analysis

### 1.2 Scoring Logic (Critical Business Rules)

#### Polarity Inversion Rules
- **Negative Questions** (e.g., "Me preocupa sufrir un accidente"):
  - Siempre = 4 points
  - Casi Siempre = 3 points
  - Algunas Veces = 2 points
  - Casi Nunca = 1 point
  - Nunca = 0 points

- **Positive Questions** (e.g., "Recibo capacitación útil"):
  - Siempre = 0 points
  - Casi Siempre = 1 point
  - Algunas Veces = 2 points
  - Casi Nunca = 3 points
  - Nunca = 4 points

#### Risk Level Thresholds (Total Score)
- **Nulo**: < 20 points
- **Bajo**: 20 - 44 points
- **Medio**: 45 - 69 points
- **Alto**: 70 - 89 points
- **Muy Alto**: ≥ 90 points

#### Category & Domain Thresholds
Each category and domain has specific thresholds (e.g., Ambiente de trabajo: Nulo <3, Bajo <5, Medio <7, Alto <9, Muy Alto ≥9)

#### Group Averaging
For General Reports: Average scores across all participants per department/domain

### 1.3 Report Structure

#### Individual Report Contains:
1. Introduction (Objective, Method)
2. Guide I Results (ATS, Persistent memories, Medical attention flag)
3. Guide II/III Results:
   - Total Score with risk meter
   - Category scores
   - Domain scores
   - Question-level breakdown (response frequency)
4. Conclusions (Action plan based on risk level)
5. Action Plan (Legal text by risk level)

#### General Report Contains:
1. Participation statistics (donut chart)
2. Demographic profiles (gender, age, marital status, education, etc.)
3. Global risk heatmap (% of staff at each risk level per domain)
4. Item analysis (stacked bar charts per question)

---

## 2. Tech Stack Validation

### 2.1 Approved Technologies ✅

| Component | Selected | Validation | Notes |
|-----------|----------|------------|-------|
| **Frontend** | Next.js 14 + React | ✅ Excellent | SSR/SSG for PDF previews, React Server Components |
| **Styling** | Tailwind CSS + Ant Design | ✅ Good | Ant Design has excellent form components for questionnaires |
| **Backend** | Go (Golang) + Gin | ✅ Excellent | High performance, great concurrency for async tasks |
| **Database** | PostgreSQL | ✅ Perfect | JSONB for demographics, strong relational integrity |
| **Infrastructure** | Docker Compose / K8s | ✅ Appropriate | Cloud-native ready |

### 2.2 Recommendations & Additions

#### Critical Additions:

1. **PDF Generation Library**
   - **Recommendation**: Use `github.com/jung-kurt/gofpdf` or `github.com/signintech/gopdf` for Go
   - **Alternative**: Consider `puppeteer` via microservice if complex charts needed
   - **Reason**: Legal PDF reports must match NOM-035 format exactly

2. **Background Job Processing**
   - **Recommendation**: Add **Redis** + **Worker Queue** (e.g., `github.com/hibiken/asynq`)
   - **Reason**: Observer pattern implementation needs async report recalculation
   - **Current**: SRS mentions Observer pattern but no queue system specified

3. **Authentication & Authorization**
   - **Recommendation**: JWT tokens (Go: `github.com/golang-jwt/jwt/v5`)
   - **Multi-tenant**: Row-level security or tenant ID isolation
   - **Reason**: Multi-tenant SaaS requires robust auth

4. **Data Encryption**
   - **Recommendation**: Use PostgreSQL `pgcrypto` extension + application-level encryption
   - **For PII**: Consider `github.com/gtank/cryptopasta` or similar
   - **Reason**: NFR-SEC-01 requires AES-256 encryption at rest

5. **Caching Layer**
   - **Recommendation**: Redis for General Report caching
   - **Reason**: Performance requirement (NFR-PERF-01: 5,000 employees)

6. **CSV Import**
   - **Recommendation**: Use `encoding/csv` (Go standard library)
   - **Validation**: Strict CURP format validation required

7. **Chart Generation**
   - **Recommendation**: 
     - Frontend: Recharts or Chart.js (for web dashboards)
     - PDF: Custom SVG rendering or chart images embedded in PDF
   - **Reason**: Reports need donut charts, heatmaps, stacked bars

#### Optional but Recommended:

8. **API Documentation**
   - **Recommendation**: OpenAPI/Swagger (`github.com/swaggo/swag`)
   - **Reason**: Better developer experience, API contract clarity

9. **Logging & Monitoring**
   - **Recommendation**: Structured logging (`github.com/sirupsen/logrus`) + Prometheus metrics
   - **Reason**: Production-grade observability

10. **Testing Framework**
    - **Recommendation**: `github.com/stretchr/testify` for Go
    - **Reason**: Scoring logic needs extensive unit tests

---

## 3. Design Patterns Analysis

### 3.1 Strategy Pattern ✅ APPROVED

**Application**: Scoring Engine
- `TraumaStrategy`: Binary Yes/No scoring
- `RiskStrategy`: Likert scale with polarity inversion

**Recommendation**: 
```go
type ScoringStrategy interface {
    CalculateScore(response Response) (int, error)
    ValidateResponse(response Response) error
}
```

**Additional Consideration**: Consider adding a `StrategyFactory` to select strategy based on guide type.

### 3.2 Composite Pattern ✅ APPROVED

**Application**: Result Aggregation (Questions → Dimensions → Domains → Categories)

**Recommendation**:
```go
type Scorable interface {
    CalculateScore() float64
    GetChildren() []Scorable
}
```

**Implementation Structure**:
- `Question` (leaf)
- `Dimension` (composite)
- `Domain` (composite)
- `Category` (composite)
- `TotalScore` (root composite)

**Critical**: Ensure recursive traversal handles polarity inversion at question level before aggregation.

### 3.3 Observer Pattern ⚠️ NEEDS INFRASTRUCTURE

**Application**: Async Stats Recalculation

**Current Gap**: SRS mentions Observer but no event/message queue specified.

**Recommendation**: 
- Use Redis Pub/Sub or message queue (asynq)
- Event-driven architecture:
  ```
  AssessmentSubmitted → Event → Worker → RecalculateGeneralReport → Cache
  ```

**Implementation**:
```go
type EventBus interface {
    Publish(event Event) error
    Subscribe(eventType string, handler EventHandler)
}
```

### 3.4 Factory Pattern ✅ APPROVED

**Application**: Report Generation (Guide II vs Guide III, subscription year)

**Recommendation**:
```go
type ReportFactory interface {
    CreateIndividualReport(assessment Assessment) (Report, error)
    CreateGeneralReport(company Company, year int) (Report, error)
}
```

**Considerations**:
- PDF template selection based on guide type
- Dynamic section inclusion/exclusion
- Legal text injection based on risk level

---

## 4. Database Schema Recommendations

### 4.1 Enhancements to Proposed Schema

#### Missing Critical Tables:

1. **`subscriptions`** (or extend `company`)
   - `subscription_start_date`
   - `subscription_end_date`
   - `status` (active, expired, cancelled)
   - `stripe_subscription_id` (for Phase 5)

2. **`assessment_links`**
   - `token` (UUID for secure link access)
   - `expires_at`
   - `accessed_at`
   - `assessment_id`
   - **Reason**: FR-ASM-01/02 mentions "link-based delivery"

3. **`categories`** and **`domains`**
   - Currently only referenced in `question_catalog`
   - Need explicit tables for hierarchy navigation
   - Store thresholds per category/domain

4. **`general_report_cache`**
   - `company_id`, `year`
   - `cache_data` (JSONB)
   - `last_updated_at`
   - **Reason**: Observer pattern needs cache storage

5. **`assessment_sessions`**
   - Track progress for partial completions
   - Store current page/section
   - **Reason**: Long questionnaires may be completed over time

#### Schema Enhancements:

6. **`question_catalog`**: Add `dimension_id`, `order_index`
7. **`responses`**: Add `answered_at` timestamp
8. **`staff`**: Consider separate `demographics` table or JSONB with schema validation
9. **`companies`**: Add `employee_count_range` (16-50, >50) for Guide II/III selection

### 4.2 Indexing Strategy

**Critical Indexes**:
- `assessments(company_id, status, created_at)`
- `responses(assessment_id, question_id)`
- `staff(company_id, curp)` (with unique constraint)
- `assessments(staff_id, subscription_year)`

---

## 5. Critical Implementation Considerations

### 5.1 Scoring Engine Complexity

**Risk**: Polarity inversion must happen BEFORE aggregation, not after.

**Recommendation**: 
```go
// Wrong: Aggregate first, invert later
// Correct: Invert at question level, then aggregate
question.Score = polarityInverter.Apply(response, question.Polarity)
domain.Score = sum(questions.Scores)
```

**Testing**: Create comprehensive test cases with known question sets and expected outcomes.

### 5.2 PDF Generation Challenges

**Risk**: Legal compliance requires exact format matching.

**Recommendations**:
1. Create reusable PDF template components
2. Use exact fonts and spacing from reference PDF
3. Generate charts as images before embedding
4. Consider PDF/A format for long-term archival
5. **Test**: Compare generated PDFs pixel-by-pixel with reference

### 5.3 Multi-Tenancy Security

**Critical**: Row-level isolation is mandatory.

**Recommendations**:
1. Always filter by `company_id` in queries
2. Use middleware to inject tenant context
3. Consider PostgreSQL Row Level Security (RLS)
4. Audit all queries for tenant isolation

### 5.4 Data Anonymization (NFR-LEG-01)

**Requirement**: Suppress charts if sample size < 5

**Implementation**:
- Check count before rendering demographic charts
- Return `null` or "Insufficient data" message
- Log suppression events for audit

### 5.5 Annual Subscription Cycle

**Critical Logic**:
- Assessments must be tied to subscription year
- Cannot start new assessments if subscription expired
- Historical reports must remain accessible

**Recommendation**:
- Store `subscription_year` with assessments
- Implement subscription validation middleware
- Separate "read" vs "write" permissions for expired subscriptions

---

## 6. Additional Recommendations

### 6.1 API Design

**Recommendations**:
1. RESTful endpoints with versioning (`/api/v1/`)
2. Separate admin vs. participant endpoints
3. Use HTTP status codes correctly (201 for creation, 202 for async)
4. Implement pagination for large lists
5. Rate limiting for public assessment links

### 6.2 Error Handling

**Recommendations**:
1. Structured error responses with error codes
2. Validation errors return 422 (Unprocessable Entity)
3. Log errors with context (company_id, user_id, request_id)
4. Don't expose internal errors to clients

### 6.3 Internationalization (i18n)

**Observation**: All content is in Spanish (Mexican regulation).

**Recommendation**: 
- Use translation keys even if only Spanish initially
- Consider `golang.org/x/text/message` or similar
- Store question text in database (already planned)

### 6.4 Migration Strategy

**Recommendations**:
1. Use database migrations (e.g., `github.com/golang-migrate/migrate`)
2. Seed script for 72 NOM-035 questions (Phase 1)
3. Version control for all schema changes
4. Rollback strategy for each migration

---

## 7. Phased Implementation Recommendations

### Phase 1 Enhancements:
- ✅ Add migration tool setup
- ✅ Include seed script structure for questions
- ✅ Set up basic project structure (cmd/, internal/, pkg/)
- ✅ Initialize Go modules with dependencies

### Phase 2 Enhancements:
- ✅ Unit tests for scoring strategies
- ✅ Integration tests with test data
- ✅ Consider table-driven tests for polarity inversion

### Phase 3 Enhancements:
- ✅ Component library structure for reusable UI
- ✅ Form validation for questionnaires
- ✅ Progress tracking for long forms

### Phase 4 Enhancements:
- ✅ PDF generation service (separate microservice candidate)
- ✅ Chart generation utilities
- ✅ Template system for reports

### Phase 5 Enhancements:
- ✅ Webhook handling for Stripe events
- ✅ Subscription lifecycle management
- ✅ Email notifications for renewal reminders

---

## 8. Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| Scoring logic errors | HIGH | Extensive unit tests, reference data validation |
| PDF format mismatch | HIGH | Automated comparison tests, manual review process |
| Multi-tenant data leak | CRITICAL | Security audit, RLS, tenant isolation testing |
| Performance with 5K employees | MEDIUM | Load testing, caching strategy, query optimization |
| Legal compliance | CRITICAL | Legal review of reports, compliance testing |
| Data encryption failures | HIGH | Encryption testing, key management best practices |

---

## 9. Final Recommendations Summary

### ✅ APPROVED AS-IS:
- Tech Stack (Next.js, Go, PostgreSQL, Docker)
- Design Patterns (Strategy, Composite, Observer, Factory)
- Database approach (PostgreSQL with JSONB)

### 🔧 REQUIRED ADDITIONS:
1. **Redis** for caching and message queue
2. **PDF Generation Library** (gofpdf or similar)
3. **JWT Authentication** library
4. **Migration Tool** (golang-migrate)
5. **Testing Framework** (testify)

### 💡 RECOMMENDED ENHANCEMENTS:
1. API documentation (Swagger)
2. Structured logging
3. Monitoring/metrics
4. Rate limiting
5. Webhook infrastructure for Stripe

### ⚠️ CRITICAL CONSIDERATIONS:
1. Scoring polarity must be applied at question level
2. Multi-tenant security must be enforced at every layer
3. PDF generation must match legal format exactly
4. Annual subscription logic must handle edge cases
5. Data anonymization (NFR-LEG-01) must be implemented correctly

---

## 10. Next Steps

1. **Project Initialization**: Set up Go module, Next.js project, Docker Compose
2. **Database Design**: Finalize schema with recommended enhancements
3. **Architecture Documentation**: Document service boundaries and interfaces
4. **Development Environment**: Docker Compose with PostgreSQL, Redis
5. **CI/CD Setup**: Automated testing and deployment pipeline

---

