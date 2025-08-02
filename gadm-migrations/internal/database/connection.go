package database

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfigFromEnv() *Config {
	port, _ := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))

	return &Config{
		Host:     getEnvOrDefault("DB_HOST", ""),
		Port:     port,
		User:     getEnvOrDefault("ADMIN_USER", ""),
		Password: getEnvOrDefault("ADMIN_PASSWORD", ""),
		DBName:   getEnvOrDefault("DB_NAME", ""),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", ""),
	}
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	var connectionString string

	if dbURL != "" {
		connectionString = dbURL
	} else {
		config := LoadConfigFromEnv()
		connectionString = config.ConnectionString()
	}

	fmt.Printf("Connecting to database...\n")

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection with retries
	if err := waitForDatabase(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database not available: %w", err)
	}

	fmt.Printf("✅ Database connection established\n")
	return pool, nil
}

func waitForDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		if err := pool.Ping(ctx); err == nil {
			return nil
		}

		if i == 0 {
			fmt.Printf("⏳ Waiting for database to be ready")
		} else {
			fmt.Printf(".")
		}

		time.Sleep(retryInterval)
	}

	fmt.Printf("\n")
	return fmt.Errorf("database not ready after %d attempts", maxRetries)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
