package engine

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/marshal"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BatchSize controls the max number of rows inserted at a time
const BatchSize = 1000

// PostgresOutputStreamer writes the output to postgresql.
//
// The name of the table is defined by `parserType`, prefixed with "stg_".
// The table is created or replaced.
//
// The columns of this table will correspond to the "Tuple.Headers()"
type PostgresOutputStreamer struct {
	conn       *pgxpool.Pool
	parserType string
}

func NewPostgresOutputStreamer(conn *pgxpool.Pool, parserType string) *PostgresOutputStreamer {
	return &PostgresOutputStreamer{conn, parserType}
}

func (out *PostgresOutputStreamer) Stream(ch chan marshal.Tuple) error {
	// Temporary: only ap data
	if out.parserType != "apconso" && out.parserType != "apdemande" {
		for range ch {
			// discard data
		}
		return nil
	}
	// End temporary

	var currentBatch []marshal.Tuple
	var headers []string

	tableName := fmt.Sprintf("stg_%s", out.parserType)

	for tuple := range ch {
		if headers == nil {
			headers = tuple.Headers()
		}

		currentBatch = append(currentBatch, tuple)

		if len(currentBatch) >= BatchSize {
			if err := insertTuples(currentBatch, out.conn, tableName, headers); err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}
			currentBatch = currentBatch[:0] // Reset currentBatch slice
		}
	}

	// Insert remaining tuples after channel closes
	if len(currentBatch) > 0 {
		if err := insertTuples(currentBatch, out.conn, tableName, headers); err != nil {
			return fmt.Errorf("failed to execute final batch: %w", err)
		}
	}

	return nil
}

func insertTuples(tuples []marshal.Tuple, conn *pgxpool.Pool, tableName string, columns []string) error {
	if len(tuples) == 0 {
		return nil
	}

	// To store the required arguments for the SQL query
	valueArgs := make([]string, 0, len(tuples))
	for _, tuple := range tuples {
		valueArgs = append(valueArgs, fmt.Sprintf("(%s)", strings.Join(tuple.Values(), ", ")))
	}

	placeholders := make([]string, len(valueArgs))
	for i := range valueArgs {
		// Start at $3 as placeholders are 1-indexed
		// and first one is for column names
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf(
		fmt.Sprintf(
			`INSERT INTO %s (%s)VALUES %s`,
			tableName,
			strings.Join(columns, ", "),
			strings.Join(valueArgs, ", "),
		),
	)

	_, err := conn.Exec(
		context.Background(),
		query,
	)
	return err
}
