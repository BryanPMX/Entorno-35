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

Billing and operational settings used by the current subscription implementation:

- Stripe checkout/webhooks:
  - `STRIPE_SECRET_KEY`
  - `STRIPE_WEBHOOK_SECRET`
  - `STRIPE_PRICE_MONTHLY`
  - `STRIPE_PRICE_YEARLY`
  - `STRIPE_SUCCESS_URL`
  - `STRIPE_CANCEL_URL`
  - `STRIPE_PORTAL_RETURN_URL`
- Pending registration cleanup worker:
  - `PENDING_REGISTRATION_CLEANUP_ENABLED`
  - `PENDING_REGISTRATION_CLEANUP_INTERVAL`
  - `PENDING_REGISTRATION_CLEANUP_RETENTION`
- Checkout rate limiting:
  - `CHECKOUT_RATE_LIMIT_ENABLED`
  - `CHECKOUT_RATE_LIMIT_LIMIT`
  - `CHECKOUT_RATE_LIMIT_WINDOW`

## Environment Setup

Production uses Portainer stack environment variables; no `.env` file on the server. For local development, create a `.env` file with the variables listed above (or see `docker-compose.prod.yml`). Never commit `.env` files to version control.

## Design Principles

- **Validation**: Database configuration is validated before use
- **Defaults**: Sensible defaults for non-critical settings
- **Security**: No hardcoded secrets or passwords
- **Type Safety**: Structured configuration types
- **Operational Hardening**: Billing, cleanup, and rate-limit settings are configurable without code changes
