# Domain Models

This package contains all GORM domain models for the Entorno35 platform.

## Models Overview

### Core Entities

1. **Company** - Multi-tenant companies
   - UUID primary key
   - RFC (unique, indexed)
   - Subscription management (status, dates, employee count)
   - Soft deletes enabled

2. **Staff** - Staff members (users)
   - UUID primary key
   - Company foreign key (CASCADE delete)
   - CURP (unique per company)
   - Demographics stored as JSONB (flexible schema)

3. **Question** - NOM-035 questionnaire items
   - Auto-increment ID
   - Guide type (I, II, III)
   - Type (Binary/Likert)
   - **Polarity (POSITIVE/NEGATIVE)** - Critical for scoring
   - Relationships to Category, Domain, Dimension

4. **Category/Domain/Dimension** - Normalized hierarchy
   - Normalized tables for efficient reporting
   - Category → Domain → Dimension → Question

5. **Assessment** - Assessment instances
   - UUID primary key
   - Staff and Company foreign keys
   - Period (e.g., 2025)
   - Pre-calculated TotalScore and RiskLevel
   - Medical attention flag (Guide I)

6. **Response** - Individual question responses
   - Assessment and Question foreign keys
   - SelectedValue (0-4)
   - CalculatedScore (after polarity inversion)
   - Unique constraint: one response per question per assessment

7. **AssessmentLink** - Secure token-based access
   - UUID token (unique index)
   - ExpiresAt timestamp
   - Links staff to assessments without authentication

## Key Features

- **UUID Support**: Primary keys use UUID for distributed systems
- **JSONB Demographics**: Flexible schema for staff demographic data
- **Soft Deletes**: Companies, Staff, and Assessments support soft deletes
- **CASCADE Deletes**: Proper data integrity (delete company → delete staff → delete assessments)
- **Enum Types**: Type-safe enums for status fields
- **BeforeCreate Hooks**: Automatic UUID generation
- **Indexes**: Performance-optimized indexes on foreign keys and lookup fields

## Usage

```go
import "github.com/entorno35/backend/internal/domain"

// Use models in GORM migrations
db.AutoMigrate(
    &domain.Company{},
    &domain.Staff{},
    &domain.Question{},
    &domain.Assessment{},
    // ...
)
```

## Important Notes

- **Polarity Field**: Critical for scoring engine - must be POSITIVE or NEGATIVE
- **Demographics JSONB**: Stores flexible fields without strict schema
- **Unique Constraints**: CURP per company, one response per question
- **Employee Count**: Determines Guide II vs Guide III selection

