# Security Notes

This document describes credential management, pre-commit checks, required environment variables, and security practices for the Entorno35 platform. Production uses Portainer stack environment variables only; no env files are committed.

## Table of contents

1. [Pre-commit security check](#pre-commit-security-check)
2. [Credentials management](#credentials-management)
3. [Securing existing env files](#securing-existing-env-files)
4. [Fixed security issues](#fixed-security-issues)
5. [Required environment variables](#required-environment-variables)
6. [Development setup](#development-setup)
7. [Docker Compose](#docker-compose)
8. [Running migrations](#running-migrations)

---

## Pre-commit security check

Before every commit:

1. **Env files**: Run `git status`. Ensure `.env`, `.env.stripe`, and any `.env.*` do **not** appear. If they are staged, run `git reset HEAD .env .env.stripe` and confirm they are in `.gitignore`.
2. **Secrets in code**: No hardcoded passwords, API keys, or tokens in tracked files. `docker-compose.prod.yml` and `docker-compose.yml` use only `${VAR}` or `${VAR:-default}`; real values are set in Portainer or local env.
3. **GitHub Actions**: Workflow uses `${{ secrets.* }}` only; no literal credentials in `.github/workflows/`.
4. **Docs**: SECURITY.md, README, and CI_CD refer to placeholders (e.g. `your_password`, `change-this-in-production`), not real values.

---

## Credentials Management

**WARNING**: NEVER commit passwords, secrets, or API keys to the repository.

### Securing existing env files

No env files are committed. All env files must stay local and must never be committed:

| File | Purpose | Committed? |
|------|---------|------------|
| `.env` | Local backend config (DB, JWT, CORS, etc.) | No (gitignored) |
| `.env.stripe` | Local Stripe keys (test/live); used only if Stripe is wired in | No (gitignored) |
| `web/frontend/.env.local` | Local frontend config (e.g. `NEXT_PUBLIC_API_URL`) | No (gitignored) |

Required backend env vars are documented in `docker-compose.prod.yml` (variable names) and in the "Required Environment Variables" section below.

**Production**: No `.env` or `.env.stripe` files on the server. All secrets (DB, JWT, CORS, Stripe, SMTP) are set in **Portainer stack environment variables** (or server env). The repo never contains real credentials.

**Before every commit**:
1. Run `git status` and ensure `.env`, `.env.stripe`, and any `.env.*` do not appear.
2. If any of them are staged, run `git reset HEAD .env .env.stripe` (and the file name) then confirm they are listed in `.gitignore`.

### Fixed Security Issues

The following security issues have been addressed:

1. **docker-compose.yml**: Removed hardcoded passwords, now uses environment variables
2. **internal/config/config.go**: Removed default password values
3. **Makefile**: Migrations now require DB_URL environment variable
4. **tests/integration/main_test.go**: Removed hardcoded password from default DB_URL fallback - DB_URL is now required
5. **vercel.json CSP**: Removed 'unsafe-eval' from Content-Security-Policy (kept 'unsafe-inline' for Next.js compatibility)
6. **.gitignore**: `.env`, `.env.*`, and `.env.stripe` are ignored. No env files are committed.

### Required Environment Variables

All sensitive configuration must be provided via environment variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password  # NEVER commit this
DB_NAME=your_database

# JWT
JWT_SECRET=your-secret-key  # NEVER commit this

# Redis (if password protected)
REDIS_PASSWORD=your_redis_password  # NEVER commit this

# SMTP Email Configuration
SMTP_HOST=smtp.gmail.com  # SMTP server hostname
SMTP_PORT=587  # SMTP port (587 for TLS, 465 for SSL)
SMTP_USER=your-email@example.com  # SMTP username/email
SMTP_PASSWORD=your-app-password  # NEVER commit this - Use App Passwords for Gmail
SMTP_FROM_ADDRESS=noreply@example.com  # Sender email address
SMTP_FROM_NAME=Entorno35 - NOM-035  # Sender display name
APP_BASE_URL=http://localhost:3000  # Base URL of your application
```

**For detailed SMTP setup instructions, see [docs/SMTP_CONFIGURATION.md](../docs/SMTP_CONFIGURATION.md)**

### Development Setup

1. For local development only: create a `.env` file with the variables listed in "Required Environment Variables" below (or in `docker-compose.prod.yml`). Never commit `.env` (it is in `.gitignore`).

2. For production: set all variables in Portainer stack environment variables; no `.env` file on the server.

### Docker Compose

For local development, docker-compose.yml uses environment variable fallbacks:
```yaml
POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-entorno35}
```

**For production**: Always set explicit environment variables, never rely on defaults.

### Running Migrations

Migrations require the `DB_URL` environment variable:

```bash
export DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
make migrate-up
```

---

**Last updated**: 2026-01-31  
**Security status**: No hardcoded secrets; credentials via Portainer env vars; pre-commit checks documented; PDF generation with secure file handling.

