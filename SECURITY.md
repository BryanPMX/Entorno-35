# Security Notes

## Credentials Management

**WARNING**: NEVER commit passwords, secrets, or API keys to the repository.

### Fixed Security Issues

The following security issues have been addressed:

1. **docker-compose.yml**: Removed hardcoded passwords, now uses environment variables
2. **internal/config/config.go**: Removed default password values
3. **Makefile**: Migrations now require DB_URL environment variable
4. **tests/integration/main_test.go**: Removed hardcoded password from default DB_URL fallback - DB_URL is now required
5. **vercel.json CSP**: Removed 'unsafe-eval' from Content-Security-Policy (kept 'unsafe-inline' for Next.js compatibility)
6. **.gitignore**: Added `.env.*` pattern to ensure all environment files (including `.env.stripe`) are properly ignored

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

**Last Updated**: 2026-01-09
**Security Status**: All security issues resolved - No hardcoded secrets, proper environment variable usage, PDF generation with secure file handling

