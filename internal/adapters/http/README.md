# HTTP Adapters

HTTP handlers and request/response DTOs following Hexagonal Architecture.

## Overview

HTTP adapters translate HTTP requests into domain operations and format domain results as HTTP responses. They depend on service interfaces, not concrete implementations.

## Structure

- **AuthHandler**: Authentication endpoints (login)
- **ScoringHandler**: Scoring endpoints (calculate assessment scores)
- **AssessmentHandler**: Assessment management endpoints (create, list, retrieve)
- **StaffHandler**: Staff management endpoints (CRUD, CSV import)
- **ReportHandler**: Report generation endpoints (individual and general reports)

## Design Principles

- **High Cohesion**: HTTP concerns only (request/response handling)
- **Low Coupling**: Depends on service interfaces, not concrete implementations
- **Thin Handlers**: Business logic belongs in services, not handlers
- **Error Handling**: Appropriate HTTP status codes and error messages

## Handler Pattern

Handlers follow this pattern:

```go
type Handler struct {
    service ServiceInterface
}

func (h *Handler) HandleRequest(c *gin.Context) {
    // 1. Parse and validate request
    // 2. Call service method
    // 3. Format and return response
}
```

## Request/Response DTOs

Handlers define DTOs (Data Transfer Objects) for request and response payloads. These are separate from domain models to allow for API versioning and evolution.

## Authentication

Most handlers require authentication via middleware. Handlers extract company context from the request context set by middleware.

