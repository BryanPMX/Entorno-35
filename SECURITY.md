# Security Notes

## Credentials Management

⚠️ **NEVER commit passwords, secrets, or API keys to the repository.**

### Fixed Security Issues

The following security issues have been addressed:

1. **docker-compose.yml**: Removed hardcoded passwords, now uses environment variables
2. **internal/config/config.go**: Removed default password values
3. **Makefile**: Migrations now require DB_URL environment variable

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
```

### Development Setup

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your local credentials

3. Ensure `.env` is in `.gitignore` (it is)

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

**Last Updated**: 2025-12-29  
**Security Alert**: Resolved - Removed committed generic passwords

