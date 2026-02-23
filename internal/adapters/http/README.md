# HTTP Adapters

HTTP handlers and request/response DTOs following Hexagonal Architecture.

## Overview

HTTP adapters translate HTTP requests into domain operations and format domain results as HTTP responses. They depend on ports and reusable application services (for example JWT, Stripe integration), while keeping transport concerns inside the adapter layer.

## Structure

- **AuthHandler**: Company authentication (`/auth/login`, RFC + password)
- **BillingHandler**: Stripe billing flows (public registration checkout, webhooks, verification, authenticated billing management)
- **ScoringHandler**: Scoring endpoints (calculate assessment scores)
- **AssessmentHandler**: Assessment management endpoints (create, list, retrieve)
- **StaffHandler**: Staff management endpoints (CRUD, CSV import)
- **ReportHandler**: Report generation endpoints (individual and general reports)

## Design Principles

- **High Cohesion**: HTTP concerns only (request/response handling)
- **Low Coupling**: Depends on ports and reusable services instead of HTTP-internal concrete concerns
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

## Billing Handler Responsibilities

`BillingHandler` now orchestrates two distinct flows while keeping HTTP concerns in the adapter layer:

- **Public paid registration flow**
  - `POST /billing/checkout-session`
  - `POST /billing/webhook`
  - `GET /billing/checkout-session/:id/verify`
  - Creates pending registrations only (not `companies`) before payment
  - Activates company accounts only after Stripe confirms checkout/subscription state

- **Authenticated existing-company billing flow**
  - `POST /api/v1/billing/checkout-session` (reactivation/new plan checkout for existing company)
  - `POST /api/v1/billing/customer-portal` (Stripe Billing Portal)
  - Derives company identity from JWT auth context only (never from request payload)

## Operational Hardening Implemented in HTTP Layer

- Stripe webhook signature verification (delegated to Stripe service)
- Webhook event idempotency integration (duplicate deliveries return `200`)
- Async activation fallback support via subscription webhooks
- Public checkout rate limiting applied by middleware (route-level)
