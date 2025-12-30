package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(dsn string) (*gorm.DB, error) {
	var err error

	// Configure GORM logger based on environment
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return DB, nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Migrate runs database migrations for all models
func Migrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	return DB.AutoMigrate(models...)
}

