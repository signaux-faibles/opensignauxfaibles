package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/jackc/tern/v2/migrate"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// MigrationFS are postgres migration files
var MigrationFS, _ = fs.Sub(migrationFS, "migrations")

// VersionTable is the name of the table to track migrations
const VersionTable = "migrations"

// runMigrations executes database migrations using Tern
func runMigrations(ctx context.Context, pool Pool) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire connection from pool: %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), VersionTable)
	if err != nil {
		return fmt.Errorf("unable to initialize migrations: %w", err)
	}

	if err := migrator.LoadMigrations(MigrationFS); err != nil {
		return fmt.Errorf("unable to load migrations: %w", err)
	}

	currentVersion, errCV := migrator.GetCurrentVersion(ctx)

	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	newVersion, errNV := migrator.GetCurrentVersion(ctx)

	if newVersion != currentVersion && errCV == nil && errNV == nil {
		slog.Info("migrations executed", "initial version", currentVersion, "new version", newVersion)
	}

	return nil
}
