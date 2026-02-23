package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Stripe   StripeConfig
	Cleanup  CleanupConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Env        string
	Port       string
	APIVersion string
	LogLevel   string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Expiry string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Origin string
}

// StripeConfig holds Stripe billing configuration.
type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	MonthlyPriceID string
	YearlyPriceID  string
	SuccessURL     string
	CancelURL      string
}

// CleanupConfig holds background data-retention job settings.
type CleanupConfig struct {
	PendingRegistrationEnabled   bool
	PendingRegistrationInterval  string
	PendingRegistrationRetention string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Env:        getEnv("ENV", "development"),
			Port:       getEnv("PORT", "8080"),
			APIVersion: getEnv("API_VERSION", "v1"),
			LogLevel:   getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-this-in-production"),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
		CORS: CORSConfig{
			Origin: getEnv("CORS_ORIGIN", "http://localhost:3000"),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
			MonthlyPriceID: getEnv("STRIPE_PRICE_MONTHLY", ""),
			YearlyPriceID:  getEnv("STRIPE_PRICE_YEARLY", ""),
			SuccessURL:     getEnv("STRIPE_SUCCESS_URL", ""),
			CancelURL:      getEnv("STRIPE_CANCEL_URL", ""),
		},
		Cleanup: CleanupConfig{
			PendingRegistrationEnabled:   getEnvBool("PENDING_REGISTRATION_CLEANUP_ENABLED", true),
			PendingRegistrationInterval:  getEnv("PENDING_REGISTRATION_CLEANUP_INTERVAL", "1h"),
			PendingRegistrationRetention: getEnv("PENDING_REGISTRATION_CLEANUP_RETENTION", "168h"),
		},
	}
}

// DatabaseURL returns the PostgreSQL connection URL
// Returns an error if required database configuration fields are missing
func (c *Config) DatabaseURL() (string, error) {
	if c.Database.User == "" {
		return "", fmt.Errorf("database user (DB_USER) is required")
	}
	if c.Database.Password == "" {
		return "", fmt.Errorf("database password (DB_PASSWORD) is required")
	}
	if c.Database.Name == "" {
		return "", fmt.Errorf("database name (DB_NAME) is required")
	}
	if c.Database.Host == "" {
		return "", fmt.Errorf("database host (DB_HOST) is required")
	}
	if c.Database.Port == "" {
		return "", fmt.Errorf("database port (DB_PORT) is required")
	}

	// URL-encode user, password, and database name to handle special characters
	// This prevents DSN malformation and security vulnerabilities
	encodedUser := url.QueryEscape(c.Database.User)
	encodedPassword := url.QueryEscape(c.Database.Password)
	encodedDBName := url.QueryEscape(c.Database.Name)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		encodedUser,
		encodedPassword,
		c.Database.Host,
		c.Database.Port,
		encodedDBName,
		url.QueryEscape(c.Database.SSLMode),
	)
	return dsn, nil
}

// RedisURL returns the Redis connection URL
func (c *Config) RedisURL() string {
	if c.Redis.Password != "" {
		// URL-encode password to handle special characters
		encodedPassword := url.QueryEscape(c.Redis.Password)
		return fmt.Sprintf("redis://:%s@%s:%s", encodedPassword, c.Redis.Host, c.Redis.Port)
	}
	return fmt.Sprintf("redis://%s:%s", c.Redis.Host, c.Redis.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
