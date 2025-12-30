# Entorno35 - Implementation Schedule & Branching Strategy

## Branching Strategy Overview

We will use a **Git Flow**-inspired strategy adapted for a phased SaaS project with both frontend and backend components.

### Branch Types

| Branch Type | Purpose | Lifecycle |
|------------|---------|-----------|
| `main` | Production-ready code | Permanent |
| `develop` | Integration branch for completed phases | Permanent |
| `phase-{n}-{feature}` | Feature branches for each phase | Temporary |
| `hotfix/*` | Critical production fixes | Temporary |
| `release/*` | Release preparation branches | Temporary |

---

## Total Branch Count

### Permanent Branches (2)
1. `main` - Production branch
2. `develop` - Development integration branch

### Phase Feature Branches (15-18)
**Phase 1**: 4-5 branches
**Phase 2**: 2-3 branches
**Phase 3**: 4-5 branches (frontend + backend integration)
**Phase 4**: 3-4 branches
**Phase 5**: 2-3 branches

**Total Estimated Branches**: ~17-23 feature branches across all phases

---

## Detailed Phase Breakdown

### Phase 1: Core Infrastructure

**Duration Estimate**: 2-3 weeks  
**Dependencies**: None  
**Parallelization**: Limited (sequential setup)

#### Branches:

1. **`phase-1/project-setup`**
   - Go module initialization
   - Next.js project setup
   - Docker Compose configuration
   - Basic folder structure (`cmd/`, `internal/`, `pkg/`, `web/`)
   - **Merge to**: `develop`

2. **`phase-1/database-schema`**
   - PostgreSQL schema design
   - Migration tool setup (golang-migrate)
   - Initial migrations (companies, staff, subscriptions)
   - **Merge to**: `develop` (after phase-1/project-setup)

3. **`phase-1/question-seeder`**
   - Question catalog schema
   - Seeder script for 72 NOM-035 questions
   - Categories, domains, dimensions setup
   - **Merge to**: `develop` (after phase-1/database-schema)

4. **`phase-1/auth-foundation`**
   - JWT authentication middleware
   - Multi-tenant context middleware
   - Basic user/company models
   - **Merge to**: `develop` (can be parallel with question-seeder)

5. **`phase-1/api-foundation`** (Optional, if time permits)
   - Gin router setup
   - Health check endpoints
   - Basic error handling
   - **Merge to**: `develop`

**Phase 1 Completion**: Merge all branches to `develop`, then `develop` → `main` (tagged as `v0.1.0`)

---

### Phase 2: Scoring Logic & Engine

**Duration Estimate**: 3-4 weeks  
**Dependencies**: Phase 1 complete  
**Parallelization**: Moderate (Strategy pattern can be developed separately)

#### Branches:

1. **`phase-2/scoring-strategy-pattern`**
   - Strategy interface definition
   - TraumaStrategy implementation (Binary Yes/No)
   - RiskStrategy interface
   - Strategy factory
   - **Merge to**: `develop`

2. **`phase-2/polarity-inversion`**
   - Polarity inversion logic
   - Question scoring service
   - Unit tests with table-driven tests
   - **Merge to**: `develop` (can be parallel with strategy-pattern)

3. **`phase-2/composite-aggregation`**
   - Composite pattern implementation
   - Hierarchical score calculation (Question → Dimension → Domain → Category)
   - Risk level determination logic
   - Threshold configuration
   - **Merge to**: `develop` (after phase-2/polarity-inversion)

4. **`phase-2/scoring-integration`** (Optional integration branch)
   - Integration tests
   - End-to-end scoring flow
   - API endpoints for scoring
   - **Merge to**: `develop`

**Phase 2 Completion**: Merge all branches to `develop`, then `develop` → `main` (tagged as `v0.2.0`)

---

### Phase 3: Frontend & Assessment Engine

**Duration Estimate**: 4-5 weeks  
**Dependencies**: Phase 2 complete  
**Parallelization**: High (Frontend and backend APIs can be developed in parallel)

#### Backend Branches:

1. **`phase-3/api-assessments`**
   - Assessment creation endpoints
   - Assessment link generation (secure tokens)
   - Assessment status management
   - Session tracking for partial completions
   - **Merge to**: `develop`

2. **`phase-3/api-responses`**
   - Response submission endpoints
   - Response validation
   - Assessment completion logic
   - Medical attention flagging (Guide I)
   - **Merge to**: `develop` (can be parallel with api-assessments)

3. **`phase-3/api-admin-staff`**
   - Staff CRUD endpoints
   - CSV import endpoint
   - Bulk operations
   - CURP validation
   - **Merge to**: `develop`

#### Frontend Branches:

4. **`phase-3/ui-authentication`**
   - Login/registration pages
   - Protected route wrapper
   - JWT token management
   - **Merge to**: `develop`

5. **`phase-3/ui-admin-dashboard`**
   - Admin layout components
   - Staff roster management UI
   - CSV upload component
   - Company settings
   - **Merge to**: `develop` (can be parallel with ui-authentication)

6. **`phase-3/ui-assessment-taking`**
   - Questionnaire rendering (Guide I, II, III)
   - Likert scale components
   - Binary Yes/No components
   - Progress tracking
   - Form validation
   - **Merge to**: `develop` (after api-assessments and api-responses)

7. **`phase-3/ui-assessment-links`**
   - Secure link page (public, no auth)
   - Assessment initialization
   - **Merge to**: `develop` (with ui-assessment-taking)

**Phase 3 Completion**: Merge all branches to `develop`, integration testing, then `develop` → `main` (tagged as `v0.3.0`)

---

### Phase 4: Reporting System

**Duration Estimate**: 4-5 weeks  
**Dependencies**: Phase 3 complete  
**Parallelization**: Moderate (Individual and General reports can be developed separately)

#### Branches:

1. **`phase-4/api-individual-report`**
   - Individual report data aggregation
   - Score calculation for reports
   - Report data model
   - **Merge to**: `develop`

2. **`phase-4/pdf-generation`**
   - PDF generation service (gofpdf)
   - PDF template system
   - Individual report PDF structure
   - Risk meter visualization
   - Action plan text injection
   - Legal format compliance
   - **Merge to**: `develop` (after phase-4/api-individual-report)

3. **`phase-4/api-general-report`**
   - General report aggregation SQL
   - Demographic statistics
   - Risk heatmap data
   - Item analysis data
   - Participation statistics
   - **Merge to**: `develop` (can be parallel with pdf-generation)

4. **`phase-4/observer-pattern-async`**
   - Redis setup
   - Message queue implementation (asynq)
   - Event bus for assessment completion
   - Background worker for report recalculation
   - Cache invalidation logic
   - **Merge to**: `develop` (critical for performance)

5. **`phase-4/ui-reports`**
   - Report viewing pages
   - Chart components (donut, heatmap, stacked bars)
   - PDF preview/download
   - General report dashboard
   - **Merge to**: `develop` (after api-general-report)

6. **`phase-4/data-anonymization`**
   - Sample size validation
   - Chart suppression logic (< 5 samples)
   - **Merge to**: `develop`

**Phase 4 Completion**: Merge all branches to `develop`, extensive testing, then `develop` → `main` (tagged as `v0.4.0`)

---

### Phase 5: Monetization & Production

**Duration Estimate**: 3-4 weeks  
**Dependencies**: Phase 4 complete  
**Parallelization**: Low (Stripe integration requires sequential steps)

#### Branches:

1. **`phase-5/subscription-management`**
   - Subscription lifecycle management
   - Annual cycle enforcement
   - Subscription validation middleware
   - Historical data access logic
   - **Merge to**: `develop`

2. **`phase-5/stripe-integration`**
   - Stripe API client setup
   - Payment processing
   - Webhook handlers
   - Subscription creation/renewal
   - **Merge to**: `develop` (after phase-5/subscription-management)

3. **`phase-5/ui-billing`**
   - Billing dashboard
   - Subscription management UI
   - Payment forms
   - Invoice display
   - **Merge to**: `develop` (after phase-5/stripe-integration)

4. **`phase-5/production-hardening`** (Optional but recommended)
   - Error monitoring (Sentry or similar)
   - Logging improvements
   - Performance optimization
   - Security audit
   - **Merge to**: `develop`

**Phase 5 Completion**: Merge all branches to `develop`, production testing, then `develop` → `main` (tagged as `v1.0.0`)

---

## Branch Merge Strategy

### Standard Merge Flow

```
feature branch → develop → main
```

1. **Feature branches** are created from `develop`
2. Work is completed and tested locally
3. Pull request is created to `develop`
4. Code review and CI/CD checks
5. Merge to `develop`
6. Integration testing on `develop`
7. After phase completion: `develop` → `main` (with tag)

### Hotfix Flow

```
main → hotfix/branch-name → main + develop
```

- Hotfixes branch from `main`
- Fix is applied and tested
- Merged to both `main` and `develop`
- Tagged as patch version (e.g., `v1.0.1`)

### Release Flow (Optional)

```
develop → release/v1.0.0 → main + develop
```

- For major releases, create release branch
- Final testing and bug fixes
- Merge to `main` and `develop`
- Tag as release version

---

## Parallelization Opportunities

### Can Work in Parallel:

1. **Phase 2**: Strategy pattern + Polarity inversion (separate concerns)
2. **Phase 3**: All frontend branches + Backend API branches (separate codebases)
3. **Phase 4**: General report API + PDF generation (after individual report API)

### Must Be Sequential:

1. **Phase 1**: Database schema → Question seeder → Auth foundation
2. **Phase 2**: Strategy pattern → Composite aggregation (dependencies)
3. **Phase 4**: Individual report API → PDF generation
4. **Phase 5**: Subscription management → Stripe integration → UI billing

---

## Timeline Estimate

| Phase | Duration | Branches | Key Deliverables |
|-------|----------|----------|------------------|
| **Phase 1** | 2-3 weeks | 4-5 | Database, Auth, Project structure |
| **Phase 2** | 3-4 weeks | 2-3 | Scoring engine, Strategy/Composite patterns |
| **Phase 3** | 4-5 weeks | 6-7 | Frontend UI, Assessment APIs |
| **Phase 4** | 4-5 weeks | 5-6 | Reporting system, PDF generation |
| **Phase 5** | 3-4 weeks | 3-4 | Stripe integration, Production ready |
| **Total** | **16-21 weeks** | **~20-25 branches** | **MVP to Production** |

**Note**: Timeline assumes 1-2 developers. With more developers, phases 3 and 4 can be accelerated through parallelization.

---

## Branch Naming Convention

```
{type}/{phase}-{feature-description}

Examples:
- phase-1/project-setup
- phase-2/scoring-strategy-pattern
- phase-3/api-assessments
- phase-3/ui-assessment-taking
- phase-4/pdf-generation
- phase-5/stripe-integration
- hotfix/critical-scoring-bug
- release/v1.0.0
```

### Branch Name Guidelines:

- Use kebab-case (lowercase with hyphens)
- Include phase number for organization
- Be descriptive but concise
- Prefix with `phase-{n}/` for feature branches
- Use `hotfix/` for urgent production fixes
- Use `release/` for release preparation

---

## CI/CD Considerations

### Required Checks Before Merge:

1. **Linting**: Code style and quality
2. **Unit Tests**: All tests must pass
3. **Integration Tests**: For affected modules
4. **Build Check**: Project compiles successfully
5. **Security Scan**: Dependency vulnerabilities

### Branch Protection Rules:

- **`main`**: Require PR review, all checks must pass, no direct commits
- **`develop`**: Require PR review, all checks must pass
- **Feature branches**: No protection (developer discretion)

---

## Git Commands Reference

### Creating a Phase Branch:

```bash
# Start from develop
git checkout develop
git pull origin develop
git checkout -b phase-1/project-setup

# Work and commit
git add .
git commit -m "feat: initialize Go module and project structure"

# Push and create PR
git push origin phase-1/project-setup
```

### Merging Phase to Develop:

```bash
# After PR approval, merge via GitHub/GitLab UI
# Or via command line:
git checkout develop
git merge phase-1/project-setup
git push origin develop
```

### Tagging Phase Completion:

```bash
git checkout main
git merge develop
git tag -a v0.1.0 -m "Phase 1: Core Infrastructure Complete"
git push origin main --tags
```

---

## Risk Mitigation

### Branch Management Risks:

| Risk | Mitigation |
|------|------------|
| **Merge conflicts** | Regular rebasing from develop, small focused branches |
| **Long-lived branches** | Phase branches merged quickly, feature branches < 1 week |
| **Integration issues** | Continuous integration testing, merge to develop frequently |
| **Broken develop** | Feature flags, rollback strategy, staging environment |

### Recommendations:

1. **Keep branches small**: Each branch should focus on one feature/module
2. **Merge frequently**: Don't let branches diverge too far from develop
3. **Test before merge**: Run tests locally and in CI before creating PR
4. **Code review**: At least one reviewer for each PR to develop/main
5. **Documentation**: Update README/docs when merging significant changes

---

## Summary

- **Total Branches**: ~20-25 feature branches across 5 phases
- **Permanent Branches**: 2 (main, develop)
- **Branch Strategy**: Git Flow adapted for phased development
- **Timeline**: 16-21 weeks for complete implementation
- **Parallelization**: High in Phase 3 (frontend/backend), Moderate in Phase 4

This branching strategy provides clear organization, enables parallel development where possible, and maintains code quality through structured merging and testing processes.

