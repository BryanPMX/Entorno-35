# Middleware Package

HTTP middleware for authentication, authorization, and multi-tenant isolation.

## Overview

Middleware functions that handle cross-cutting concerns for HTTP requests. Middleware is applied to route groups in the HTTP router.

## Components

### AuthMiddleware

Validates JWT tokens from incoming requests and sets authentication context in the Gin context.

**Usage:**
```go
router.Use(middleware.AuthMiddleware(jwtService))
```

**Requirements:**
- JWT token in `Authorization: Bearer <token>` header
- Valid, non-expired token

### TenantMiddleware

Enforces multi-tenant isolation by validating that requests are scoped to the authenticated company. Must be used after `AuthMiddleware`.

**Usage:**
```go
router.Use(middleware.AuthMiddleware(jwtService))
router.Use(middleware.TenantMiddleware())
```

**Features:**
- Validates company ID matches authenticated user
- Sets company ID and staff ID in context for easy access

### Helper Functions

- `GetAuthContext(c *gin.Context) (*auth.Context, bool)` - Extract auth context
- `RequireAuth(c *gin.Context) (*auth.Context, bool)` - Require and extract auth context (aborts if missing)
- `RequireCompanyID(c *gin.Context) (uuid.UUID, bool)` - Extract company ID (aborts if missing)
- `RequireStaffID(c *gin.Context) (uuid.UUID, bool)` - Extract staff ID (aborts if missing)

## Design Principles

- **High Cohesion**: HTTP middleware concerns only
- **Low Coupling**: Depends on interfaces (jwt.Service), not concrete implementations
- **Chain of Responsibility**: Middleware chain processes requests sequentially
- **Context Propagation**: Auth data passed via Gin context
