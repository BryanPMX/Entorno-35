# Scoring Package

NOM-035 scoring logic following Test-Driven Development (TDD) principles.

## Design Principles

- **TDD Approach**: Tests written first, implementation follows
- **High Cohesion**: Scoring logic only
- **Low Coupling**: Depends on domain types, not database or HTTP
- **Pure Functions**: No side effects, easy to test
- **Strategy Pattern**: Different scoring strategies for different guide types

## Components

### Interfaces

- `ScoringStrategy`: Interface for scoring implementations
- `AssessmentResult`: Result structure with scores and risk levels

### Strategies

- **TraumaStrategy**: Guide I (Trauma Assessment) - Binary Yes/No logic
- **RiskStrategy**: Guide II/III (Psychosocial Risk) - Likert scale with polarity inversion

### Polarity Application

Applies polarity inversion to Likert scale responses (Guide II/III):
- **Positive**: No inversion (0→0, 1→1, 2→2, 3→3, 4→4)
- **Negative**: Inverted (0→4, 1→3, 2→2, 3→1, 4→0)

### Risk Level Calculation

Calculates risk levels (Nulo, Bajo, Medio, Alto, Muy Alto) based on score thresholds:
- Uses official NOM-035 thresholds from official document
- Boundary-inclusive logic (e.g., score = 50 → Bajo, not Nulo)

### Scoring Rules

Hardcoded configuration from official NOM-035 document:
- Guide II thresholds (16-50 employees)
- Guide III thresholds (>50 employees)
- Category and Domain thresholds

## Testing

All scoring logic is thoroughly tested:

```bash
go test ./internal/core/scoring/... -v
```

## Usage

```go
import "github.com/entorno35/backend/internal/core/scoring"

// Create strategy
strategy := scoring.NewTraumaStrategy() // or NewRiskStrategy()

// Calculate scores
result, err := strategy.Calculate(responses)

// Apply polarity
score, err := scoring.ApplyPolarity(selectedValue, polarity)

// Load rules
rules := scoring.LoadScoringRules()
```

