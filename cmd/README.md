# Cmd Package

Application entry points (CLI tools and servers).

## Structure

```
cmd/
├── api/              # Main API server
└── seeder/           # Database seeding CLI tool
```

## Design Principles

- **Single Responsibility**: Each command has one purpose
- **Thin Entry Points**: Commands orchestrate services, do not contain business logic
- **Dependency Injection**: Services and repositories are initialized here

## Commands

### api/

Main HTTP API server entry point. Initializes:
- Database connection
- HTTP router (Gin)
- Services and repositories
- Middleware
- Route handlers

### seeder/

CLI tool for populating the database with initial data:
- NOM-035 questions from JSON file
- Categories, domains, dimensions
- Idempotent operations (safe to run multiple times)

Usage:
```bash
go run cmd/seeder/main.go
```

