# Pkg Package

Public reusable packages that can be imported by external applications.

## Structure

```
pkg/
└── errors/           # Error handling utilities
```

## Design Principles

- **Public API**: Code in `pkg/` is intended for external use
- **Reusability**: Generic utilities not tied to application-specific logic
- **Documentation**: Public packages should be well-documented

## Packages

### errors/

Error handling utilities for structured error responses. Provides helpers for creating API errors with appropriate HTTP status codes.

## Usage

External packages can import from `pkg/`:

```go
import "github.com/entorno35/backend/pkg/errors"
```

