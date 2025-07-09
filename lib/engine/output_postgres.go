package engine

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/marshal"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BatchSize controls the max number of rows inserted at a time
const BatchSize = 1000

// PostgresOutputStreamer writes the output to postgresql.
//
// The name of the table is defined by `parserType`, prefixed with "stg_".
// The table is expected to exist and be properly formatted. It is truncated
// before inserting new values.
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
	logger := slog.With("parser", out.parserType)

	// Temporary: only ap data
	if out.parserType != "apconso" && out.parserType != "apdemande" {
		for range ch {
			// discard data
		}
		return nil
	}
	// End temporary

	logger.Debug("stream output to PostgreSQL")

	var currentBatch []marshal.Tuple
	var headers []string

	tableName := fmt.Sprintf("stg_%s", out.parserType)

	logger.Debug("truncate table", "table", tableName)

	_, err := out.conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s",
		tableName))

	if err != nil {
		return err
	}

	logger.Debug("data insertion", "table", tableName)

	nInserted := 0

	for tuple := range ch {
		if headers == nil {
			headers = tuple.Headers()
		}

		currentBatch = append(currentBatch, tuple)

		if len(currentBatch) >= BatchSize {

			if err := insertTuples(currentBatch, out.conn, tableName, headers); err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}

			nInserted += len(currentBatch)

			currentBatch = currentBatch[:0] // Reset currentBatch slice
		}
	}

	// Insert remaining tuples after channel closes
	if len(currentBatch) > 0 {

		if err := insertTuples(currentBatch, out.conn, tableName, headers); err != nil {
			return fmt.Errorf("failed to execute final batch: %w", err)
		}

		nInserted += len(currentBatch)
	}

	logger.Debug("output streaming to PostgreSQL ended successfully", "n_inserted", nInserted)

	return nil
}

func insertTuples(tuples []marshal.Tuple, conn *pgxpool.Pool, tableName string, columns []string) error {
	if len(tuples) == 0 {
		return nil
	}

	values := make([][]any, 0, len(tuples))

	// TODO rather than construct values
	// implement CopyFromSource interface
	for _, tuple := range tuples {
		newSlice := make([]any, 0, len(tuple.Values()))
		for _, value := range tuple.Values() {
			newSlice = append(newSlice, value)
		}
		values = append(values, newSlice)
	}
	lowerColumns := make([]string, len(columns))
	for i, c := range columns {
		lowerColumns[i] = strings.ToLower(c)
	}

	// Batch insertion
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{tableName},
		lowerColumns,
		pgx.CopyFromRows(values),
	)

	return err
}
