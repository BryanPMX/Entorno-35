# Services Package

Business logic services that orchestrate repositories and core logic.

## Overview

Services implement use cases and coordinate between repositories (data access) and core business logic. They follow the Service Layer pattern from Domain-Driven Design.

## Design Principles

- **High Cohesion**: Each service handles one business domain
- **Low Coupling**: Depends on repository interfaces (ports), not concrete implementations
- **Single Responsibility**: Each service has a focused purpose
- **Transaction Boundaries**: Services typically represent transaction boundaries

## Services

### ScoringService

Handles assessment scoring operations:
- Calculates scores using appropriate strategy based on guide type
- Coordinates between AssessmentRepository and scoring strategies
- Updates assessment records with calculated results

## Usage

Services are typically instantiated in `cmd/api/main.go` and passed to HTTP handlers:

```go
assessmentRepo := postgres.NewAssessmentRepository(db)
scoringService := services.NewScoringService(assessmentRepo)
```

## Future Services

Additional services will be added as the application grows:
- AssessmentService: Assessment creation and management
- ResponseService: Response submission and validation
- StaffService: Staff CRUD operations
- ReportService: Report generation and aggregation

