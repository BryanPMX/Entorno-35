# Migrations

Database schema migrations using SQL files.

## Overview

Migrations are SQL scripts that version control database schema changes. Use `golang-migrate` tool to apply migrations.

## Migration Files

- `001_initial_schema.up.sql` - Creates initial database schema
- `001_initial_schema.down.sql` - Reverts initial schema

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

