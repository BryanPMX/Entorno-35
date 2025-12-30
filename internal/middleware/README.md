# Middleware Package

HTTP middleware for authentication and multi-tenant isolation (high cohesion - HTTP concerns only).

## Design Principles

- **High Cohesion**: Single responsibility - HTTP middleware concerns only
- **Low Coupling**: Depends on interfaces (`jwt.Service`), not concrete implementations
- **Dependency Inversion**: Uses service interfaces to enable testing and extensibility

## Components

### AuthMiddleware

Validates JWT tokens and sets auth context in the request.

```go
jwtService := jwt.NewService(config.JWT.Secret)
router.Use(middleware.AuthMiddleware(jwtService))
```

### TenantMiddleware

Enforces multi-tenant isolation by validating company ID matches.

```go
router.Use(middleware.TenantMiddleware())
```

**Note**: Must be used after `AuthMiddleware` (depends on auth context).

### Helper Functions

- `GetAuthContext(c *gin.Context) (*auth.Context, bool)` - Extract auth context
- `RequireAuth(c *gin.Context) (*auth.Context, bool)` - Require and extract auth context
- `RequireCompanyID(c *gin.Context) (uuid.UUID, bool)` - Extract company ID
- `RequireStaffID(c *gin.Context) (uuid.UUID, bool)` - Extract staff ID

## Usage Example

```go
// Setup middleware
jwtService := jwt.NewService(config.JWT.Secret)
router.Use(middleware.AuthMiddleware(jwtService))
router.Use(middleware.TenantMiddleware())

// In handlers
func GetCompanyData(c *gin.Context) {
    companyID, ok := middleware.RequireCompanyID(c)
    if !ok {
        return // Already aborted with error
    }
    
    // Use companyID for data access
}
```

