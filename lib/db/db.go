// Package db defines an interface for Database operations.
package db

import (
	"context"
	"fmt"
	"log"

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
		log.Println("Mode NO_DB: pas de lecture et d'écriture de la base de données")
		mockPool, _ := pgxmock.NewPool()

		DB = mockPool

		return nil
	}

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, viper.GetString("POSTGRES_DB_URL"))

	if err != nil {
		return fmt.Errorf("erreur de connexion à PostgreSQL : %w", err)
	}

	// Test connectivity with postgreSQL database
	err = conn.Ping(ctx)

	if err != nil {
		return fmt.Errorf("erreur de connexion à PostgreSQL : %w", err)
	}

	// Run database migrations
	if err := runMigrations(ctx, conn); err != nil {
		return fmt.Errorf("erreur lors de l'exécution des migrations : %w", err)
	}

	DB = conn
	return nil
}
