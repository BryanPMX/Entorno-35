# Database Package

Database connection utilities and migration helpers.

## Overview

Provides functions to establish database connections and manage migrations using GORM.

## Usage

```go
import "github.com/entorno35/backend/internal/database"

dsn := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
db, err := database.Connect(dsn)
defer database.Close()
```

## Connection Management

The package maintains a global database connection that should be closed on application shutdown.

## Connection Pool

Connection pool settings are configured automatically:
- Max idle connections: 10
- Max open connections: 100

## Migration Support

The `Migrate` function supports GORM AutoMigrate for schema updates. For production, use SQL migration files in the `migrations/` directory.

## Design Principles

- **Single Connection**: Global connection instance for the application
- **Resource Management**: Proper connection cleanup
- **Configuration**: Connection pool tuning for production use

