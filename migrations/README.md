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

