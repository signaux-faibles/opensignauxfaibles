// Package db defines all necessary interfaces for writing to and
// reading from the PostgreSQL database.
//
// It also defines all tables and views in the form of
// migrations.
//
// The database has a two-layer architecture :
// - Tables prefixed with `stg_` represent imported data, relatively raw
// (although a number of quality operations are already
// performed at import time).
// - Tables and views prefixed with `clean_` are enriched and
// cleaned tables.
package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/spf13/viper"
)

// DB is a connection pool, set with `Init`
var DB Pool

// Pool is the subset of the pgxpool.Pool interface we actually use
type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close()
}

// Init set Db global variable to a connection pool, or a mock pool if noDB is
// set to true
func Init(noDB bool) error {

	if noDB {
		slog.Info("NO_DB mode : no reading from and writing to the database")
		mockPool, _ := pgxmock.NewPool()

		DB = mockPool

		return nil
	}

	connStr := viper.GetString("POSTGRES_DB_URL")

	conf, err := pgx.ParseConfig(connStr)
	if err != nil {
		// Do not log connexion string as it can contain user / password
		// information
		slog.Error("could not properly parse connexion string provided by POSTGRES_DB_URL environment variable")
		return err
	}

	logger := slog.With("host", conf.Host, "port", conf.Port, "database", conf.Database)
	logger.Info("connecting to database...")

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, viper.GetString("POSTGRES_DB_URL"))

	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test connectivity with postgreSQL database
	err = conn.Ping(ctx)

	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	logger.Info("database connexion established")

	// Run database migrations
	logger.Info("running database migrations...")

	if err := runMigrations(ctx, conn); err != nil {
		return fmt.Errorf("failed to execute database migrations: %w", err)
	}
	logger.Info("database migrated with success")

	DB = conn
	return nil
}
