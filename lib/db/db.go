// Package db defines all necessary interfaces for writing to and
// reading from the PostgreSQL database.
//
// It also defines all tables and views in the form of
// migrations.
//
// The database has a two-layer architecture :
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
	"strings"

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
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close()
}

// Init set Db global variable to a connection pool
func Init(schema string, shouldMigrate bool) error {

	connStr := viper.GetString("POSTGRES_DB_URL")

	conf, err := pgx.ParseConfig(connStr)
	if err != nil {
		// Do not log connexion string as it can contain user / password
		// information
		slog.Error("could not properly parse connexion string provided by POSTGRES_DB_URL environment variable")
		return err
	}

	logger := slog.With("host", conf.Host, "port", conf.Port, "database", conf.Database, "schema", schema)
	logger.Info("connecting to database...")

	ctx := context.Background()

	// Step 1: Create schema using a temporary connection (without search_path)
	tmpConn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	_, err = tmpConn.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", pgx.Identifier{schema}.Sanitize()))
	tmpConn.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to create schema %q: %w", schema, err)
	}

	// Step 2: Create the pool with search_path set as a connection parameter.
	conn, err := pgxpool.New(ctx, appendSearchPath(connStr, schema))
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
	if shouldMigrate {
		logger.Info("running database migrations...")

		if err := runMigrations(ctx, conn); err != nil {
			return fmt.Errorf("failed to execute database migrations: %w", err)
		}
		logger.Info("database migrated with success")
	}

	DB = conn
	return nil
}

// appendSearchPath adds search_path to a PostgreSQL connection string.
func appendSearchPath(connStr string, schema string) string {
	if strings.Contains(connStr, "?") {
		return connStr + "&search_path=" + schema
	}
	return connStr + "?search_path=" + schema
}

func InitMock() {

	slog.Info("NO DB mode : no reading from and writing to the database")
	mockPool, _ := pgxmock.NewPool()

	DB = mockPool
}
