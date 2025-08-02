package migration

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	pool          *pgxpool.Pool
	db            *sql.DB
	migrationsDir string
}

func New(pool *pgxpool.Pool, migrationsDir string) *Migrator {
	// Convert pgxpool to sql.DB for goose compatibility
	db := stdlib.OpenDBFromPool(pool)

	return &Migrator{
		pool:          pool,
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up(ctx context.Context, targetVersion int64) error {
	fmt.Printf("🚀 Applying migrations from %s\n", m.migrationsDir)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return fmt.Errorf("migration up failed: %w", err)
	}
	fmt.Println("✅ All migrations applied successfully")

	return nil
}
