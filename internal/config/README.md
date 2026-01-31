# Config Package

Configuration management for the application.

## Overview

Loads configuration from environment variables and provides structured access to application settings.

## Usage

```go
import "github.com/entorno35/backend/internal/config"

cfg := config.Load()
dsn, err := cfg.DatabaseURL()
```

## Configuration Sources

Configuration is loaded from environment variables. Required variables:

- Database: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- JWT: `JWT_SECRET`, `JWT_EXPIRY`
- Redis: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- Server: `PORT`, `ENV`
- CORS: `CORS_ORIGIN` (defaults to `http://localhost:3000`)

## Environment Setup

Production uses Portainer stack environment variables; no `.env` file on the server. For local development, create a `.env` file with the variables listed above (or see `docker-compose.prod.yml`). Never commit `.env` files to version control.

## Design Principles

- **Validation**: Database configuration is validated before use
- **Defaults**: Sensible defaults for non-critical settings
- **Security**: No hardcoded secrets or passwords
- **Type Safety**: Structured configuration types

