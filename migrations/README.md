# Migrations

Database schema migrations using SQL files.

## Overview

Migrations are SQL scripts that version control database schema changes. Use `golang-migrate` tool to apply migrations.

## Migration Files

- `001_initial_schema.up.sql` - Creates initial database schema
- `001_initial_schema.down.sql` - Reverts initial schema
- `002_add_staff_password.up.sql` - Adds password authentication to staff table
- `002_add_staff_password.down.sql` - Removes password column from staff table
- `003_add_employee_id.up.sql` - Adds employee_id field for staff without CURP
- `003_add_employee_id.down.sql` - Removes employee_id column from staff table
- `004_add_company_billing_and_admin_fields.up.sql` - Adds company admin credentials and Stripe billing fields
- `004_add_company_billing_and_admin_fields.down.sql` - Removes company billing/admin fields
- `005_add_pending_company_registrations.up.sql` - Adds pending pre-payment company registration table and indexes
- `005_add_pending_company_registrations.down.sql` - Removes pending registration table and indexes
- `006_add_stripe_webhook_events.up.sql` - Adds Stripe webhook idempotency table and indexes
- `006_add_stripe_webhook_events.down.sql` - Removes webhook idempotency table and indexes

## Usage

### Apply Migrations

```bash
export DB_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
make migrate-up
```

### Rollback Migrations

```bash
export DB_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
make migrate-down
```

## Naming Convention

Migrations follow the pattern:
- `{version}_{description}.up.sql` - Forward migration
- `{version}_{description}.down.sql` - Rollback migration

## Best Practices

- Always create both up and down migrations
- Test migrations on a copy of production data
- Never modify existing migrations (create new ones instead)
- Ensure migrations are idempotent when possible

## Important Production Note (Billing Migrations)

The API also runs GORM `AutoMigrate` at startup, which can create tables/basic indexes. However, SQL migrations remain the source of truth for SQL-only constraints and triggers, including:

- partial unique index `uq_pending_company_registrations_open_rfc`
- `update_*_updated_at` triggers on billing tables

For production, apply SQL migrations (`make migrate-up`) even if the API starts successfully after `AutoMigrate`.
