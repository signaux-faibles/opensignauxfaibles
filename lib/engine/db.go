package engine

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/spf13/viper"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// MigrationFS are postgres migration files
var MigrationFS, _ = fs.Sub(migrationFS, "migrations")

// VersionTable is the name of the table to track migrations
const VersionTable = "sfdata_migrations"

// DB type centralisant les accès à une base de données
type DB struct {
	PostgresDB *pgxpool.Pool
}

// InitDB Initialisation de la la base de données PostgreSQL.
// Cette fonction réalise les migrations - le cas échéant - de la base
// PostgreSQL.
func InitDB() (DB, error) {
	conn, err := pgxpool.New(context.Background(), viper.GetString("POSTGRES_DB_URL"))

	if err == nil {
		// Test connectivity with postgreSQL database
		err = conn.Ping(context.Background())
	}

	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion à PostgreSQL : %w", err)
	}

	// Run database migrations
	if err := runMigrations(conn); err != nil {
		return DB{}, fmt.Errorf("erreur lors de l'exécution des migrations : %w", err)
	}

	return DB{
		PostgresDB: conn,
	}, nil
}

// runMigrations executes database migrations using Tern
func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("impossible d'acquérir une connexion du pool : %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "sfdata_migrations")
	if err != nil {
		return fmt.Errorf("impossible d'initialiser les migrations : %w", err)
	}

	if err := migrator.LoadMigrations(MigrationFS); err != nil {
		return fmt.Errorf("impossible de charger les migrations : %w", err)
	}

	currentVersion, errCV := migrator.GetCurrentVersion(ctx)

	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("échec de la migration : %w", err)
	}

	newVersion, errNV := migrator.GetCurrentVersion(ctx)

	if newVersion != currentVersion && errCV == nil && errNV == nil {
		slog.Info("Migrations exécutées", "version initiale", currentVersion, "nouvelle version", newVersion)
	}

	return nil
}
