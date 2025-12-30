# Internal Package

Private application code following Go best practices and Hexagonal Architecture.

## Structure

```
internal/
├── adapters/          # External adapters (HTTP, database)
├── auth/              # Authentication context and types
├── config/            # Configuration management
├── core/              # Core business logic and ports
├── database/          # Database connection utilities
├── domain/            # Domain models and entities
├── middleware/        # HTTP middleware
└── services/          # Business logic services
```

## Design Principles

- **Private Package**: Code in `internal/` is not importable by external packages
- **High Cohesion**: Each package has a single, well-defined responsibility
- **Low Coupling**: Packages depend on interfaces, not concrete implementations
- **Hexagonal Architecture**: Clear separation between domain, ports, and adapters

## Package Overview

### adapters/
External adapters implementing ports defined in `core/ports/`. Includes HTTP handlers and database repositories.

### auth/
Authentication context structures and helper functions for managing authenticated user sessions.

### config/
Configuration management for loading environment variables and providing structured access to application settings.

### core/
Core business logic including:
- **ports/**: Interface definitions (hexagonal architecture ports)
- **scoring/**: NOM-035 scoring logic and strategies
- **jwt/**: JWT token service
- **password/**: Password hashing service

### database/
Database connection utilities and migration helpers.

### domain/
GORM domain models representing the business entities (Company, Staff, Assessment, Question, etc.).

### middleware/
HTTP middleware for authentication, authorization, and multi-tenant isolation.

### services/
Business logic services that orchestrate repositories and core logic.

