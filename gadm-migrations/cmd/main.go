package main

import (
	"context"
	"fmt"
	"log"

	"gadm-migrations/internal/database"
	"gadm-migrations/internal/migration"

	"github.com/spf13/cobra"
)

var (
	dbURL         string
	migrationsDir string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gadm-migrate",
		Short: "GADM API Database Migration Tool",
	}

	rootCmd.PersistentFlags().StringVar(&dbURL, "db-url", "", "Database connection URL")
	rootCmd.PersistentFlags().StringVar(&migrationsDir, "migrations-dir", "./migrations", "Migrations directory")

	rootCmd.AddCommand(upCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func upCmd() *cobra.Command {
	var targetVersion int64
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			pool, err := database.Connect(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}
			defer pool.Close()

			migrator := migration.New(pool, migrationsDir)
			return migrator.Up(ctx, targetVersion)
		},
	}
	cmd.Flags().Int64Var(&targetVersion, "target", 0, "Target version (0 for latest)")
	return cmd
}
