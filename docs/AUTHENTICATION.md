# Authentication Architecture

## Overview

The authentication system follows **high cohesion** and **low coupling** principles, using interfaces and dependency inversion to enable testability and extensibility.

## Architecture Principles

### High Cohesion

Each package has a single, well-defined responsibility:

- **`internal/core/jwt`**: JWT token operations only
- **`internal/core/password`**: Password hashing/verification only
- **`internal/auth`**: Auth domain types only
- **`internal/middleware`**: HTTP middleware concerns only

### Low Coupling

- Components depend on **interfaces**, not concrete implementations
- Services can be mocked for testing
- Easy to swap implementations (e.g., different JWT libraries, hashing algorithms)

### Dependency Inversion

- High-level modules (middleware) depend on abstractions (interfaces)
- Low-level modules (implementations) implement those interfaces
- Both depend on abstractions

## Package Structure

```
internal/
├── core/
│   ├── jwt/
│   │   ├── service.go      # JWT Service interface + implementation
│   │   └── README.md
│   └── password/
│       ├── hasher.go       # Password Hasher interface + implementation
│       └── README.md
├── auth/
│   ├── context.go          # Auth context domain type
│   └── README.md
└── middleware/
    ├── auth.go             # JWT authentication middleware
    ├── tenant.go           # Multi-tenant isolation middleware
    └── README.md
```

## Components

### 1. JWT Service (`internal/core/jwt`)

**Interface**: `jwt.Service`

```go
type Service interface {
    GenerateToken(companyID, staffID, email, role string, expiresIn time.Duration) (string, error)
    ValidateToken(tokenString string) (*Claims, error)
    RefreshToken(tokenString string, expiresIn time.Duration) (string, error)
}
```

**Implementation**: `jwtService` (HS256, JWT v5)

**Benefits**:
- Can be mocked for testing
- Easy to swap for RS256 or other algorithms
- Single responsibility (JWT operations only)

### 2. Password Hasher (`internal/core/password`)

**Interface**: `password.Hasher`

```go
type Hasher interface {
    Hash(password string) (string, error)
    Verify(hashedPassword, password string) error
}
```

**Implementation**: `bcryptHasher` (bcrypt with configurable cost)

**Benefits**:
- Can be mocked for testing
- Easy to swap for Argon2, scrypt, etc.
- Single responsibility (password operations only)

### 3. Auth Context (`internal/auth`)

Simple value type containing authenticated user's context:

```go
type Context struct {
    CompanyID uuid.UUID
    StaffID   *uuid.UUID  // Optional
    Email     string
    Role      string
}
```

**Benefits**:
- No dependencies on other packages
- Pure data structure
- Easy to serialize/pass around

### 4. Middleware (`internal/middleware`)

#### AuthMiddleware

- Validates JWT tokens from `Authorization: Bearer <token>` header
- Sets auth context in Gin context
- Aborts request if token is invalid/expired

#### TenantMiddleware

- Enforces multi-tenant isolation
- Validates company ID matches authenticated user's company
- Sets company/staff IDs in context for easy access

**Benefits**:
- Depends on `jwt.Service` interface, not concrete type
- Can be tested with mock JWT service
- Single responsibility (HTTP middleware only)

## Usage Flow

### 1. Token Generation (in auth service/endpoint)

```go
jwtService := jwt.NewService(config.JWT.Secret)
token, err := jwtService.GenerateToken(
    companyID.String(),
    staffID.String(),
    email,
    role,
    24 * time.Hour,
)
```

### 2. Request Flow

```
Request → AuthMiddleware → TenantMiddleware → Handler
          (validate JWT)    (validate tenant)  (use context)
```

### 3. Handler Usage

```go
func GetStaff(c *gin.Context) {
    // Extract auth context
    authCtx, ok := middleware.RequireAuth(c)
    if !ok {
        return // Already aborted
    }
    
    // Or extract IDs directly
    companyID, ok := middleware.RequireCompanyID(c)
    if !ok {
        return
    }
    
    // Use for data access
    // ...
}
```

## Testing Strategy

### Unit Tests

- Mock `jwt.Service` interface for middleware tests
- Mock `password.Hasher` interface for auth service tests
- Test each component in isolation

### Integration Tests

- Use real implementations
- Test token generation → validation flow
- Test middleware chain

## Security Considerations

1. **JWT Secret**: Must be strong, stored in environment variables
2. **Token Expiration**: Configurable per use case
3. **Password Hashing**: Uses bcrypt with appropriate cost
4. **Tenant Isolation**: Enforced at middleware level
5. **HTTPS**: Required in production (tokens in headers)

## Future Extensibility

### Possible Enhancements

1. **Token Blacklisting**: Add Redis-based token blacklist
2. **Refresh Tokens**: Separate refresh token mechanism
3. **OAuth2**: Add OAuth2 provider support
4. **MFA**: Multi-factor authentication
5. **Session Management**: Alternative to JWT (if needed)

All can be added without breaking existing code due to interface-based design.

---

**Last Updated**: 2025-12-29  
**Design Pattern**: Strategy Pattern, Dependency Inversion Principle

