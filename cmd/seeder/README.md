# Database Seeder

This tool populates the database with NOM-035 questions from the JSON file.

## Usage

### Prerequisites

1. Database must be migrated (run migrations first)
2. Environment variables must be set (DB_USER, DB_PASSWORD, DB_NAME, etc.)
3. JSON file `nom035_questions.json` must be in the project root

### Running the Seeder

```bash
# Using go run (defaults to nom035_questions.json)
go run cmd/seeder/main.go

# Or specify a custom JSON file path
go run cmd/seeder/main.go path/to/questions.json
```

### Using Makefile (recommended)

Add to your Makefile:

```makefile
seed: ## Seed database with questions
	@echo "Seeding database..."
	@go run cmd/seeder/main.go
```

Then run:
```bash
make seed
```

## What It Does

1. **Connects to Database**: Uses configuration from environment variables
2. **Parses JSON**: Reads `nom035_questions.json` structure
3. **Hierarchy Creation**: 
   - For Guide I: Creates categories/domains from section/subsection (for organizational purposes)
   - For Guide II/III: Creates normalized Category → Domain → Dimension hierarchy
4. **Question Creation**: Creates or updates questions with proper relationships
5. **Idempotent**: Safe to run multiple times (uses Find-or-Create pattern)

## Guide I Mapping

- `section` → Category table (created for consistency) + Question.Section field
- `subsection` → Domain table (created for consistency) + Question.Subsection field
- **Note**: Guide I questions store section/subsection in their own fields
- No CategoryID/DomainID FK relationships (these are NULL in Question table)

## Guide II/III Mapping

- `category` → Category table (find or create by name)
- `domain` → Domain table (find or create by name + category_id)
- `dimension` → Dimension table (find or create by name + domain_id)
- `polarity` → Question.Polarity field (POSITIVE/NEGATIVE enum)

## Output

The seeder will log:
- Database connection status
- Each question created/updated
- Completion status for each guide
- Total completion message

## Notes

- **Risk Thresholds**: Risk threshold data is NOT seeded by this tool. The thresholds (stored in `RiskStrategy.json` or similar) should be loaded as configuration in Phase 2's scoring engine (e.g., `internal/core/services/scoring_rules.go`). These are static business rules and don't need database storage.
- **Idempotency**: Running the seeder multiple times will update existing questions, not create duplicates (uses Find-or-Create pattern).
- **Transaction Safety**: Consider wrapping the entire seed operation in a transaction for production use if you need atomic updates.

