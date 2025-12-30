# Password Hasher

High cohesion, low coupling password hashing service.

## Design Principles

- **High Cohesion**: Single responsibility - handles only password hashing/verification
- **Low Coupling**: Uses interface (`Hasher`) to allow different implementations
- **Security**: Uses bcrypt with configurable cost

## Usage

```go
import "github.com/entorno35/backend/internal/core/password"

// Create hasher
hasher := password.NewDefaultHasher()

// Hash password
hashed, err := hasher.Hash("plaintext-password")

// Verify password
err := hasher.Verify(hashed, "plaintext-password")
if err != nil {
    // Password doesn't match
}
```

## Interface

The `Hasher` interface allows for:
- Easy testing (mock implementation)
- Future implementations (e.g., Argon2, scrypt)
- Dependency injection

