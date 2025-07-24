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

type PostgresSinkFactory struct {
	conn *pgxpool.Pool
}

func NewPostgresSinkFactory(conn *pgxpool.Pool) SinkFactory {
	return &PostgresSinkFactory{conn}
}

func (f *PostgresSinkFactory) CreateSink(parserType string) (DataSink, error) {
	// Temporary: only ap data
	if parserType == "apconso" || parserType == "apdemande" {
		tableName := fmt.Sprintf("stg_%s", parserType)
		return &PostgresSink{f.conn, tableName}, nil
	}
	return &DiscardDataSink{}, nil

}

// PostgresSink writes the output to postgresql.
//
// The name of the table is defined by `parserType`, prefixed with "stg_".
// The table is expected to exist and be properly formatted. It is truncated
// before inserting new values.
//
// The columns of this table will correspond to the "Tuple.Headers()"
type PostgresSink struct {
	conn *pgxpool.Pool

	// Name of the table to which to write
	table string
}

func NewPostgresSink(conn *pgxpool.Pool) *PostgresSink {
	return &PostgresSink{conn, ""}
}

func (s *PostgresSink) ProcessOutput(ch chan marshal.Tuple) error {
	logger := slog.With("sink", "postgresql", "table", s.table)

	logger.Debug("stream output to PostgreSQL")

	var currentBatch []marshal.Tuple
	var headers []string

	logger.Debug("truncate table")

	_, err := s.conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s", s.table))

	if err != nil {
		return err
	}

	logger.Debug("data insertion")

	nInserted := 0

	for tuple := range ch {
		if headers == nil {
			headers = marshal.ExtractTableColumns(tuple)
		}

		currentBatch = append(currentBatch, tuple)

		if len(currentBatch) >= BatchSize {

			if err := insertTuples(currentBatch, s.conn, s.table, headers); err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}

			nInserted += len(currentBatch)

			currentBatch = currentBatch[:0] // Reset currentBatch slice
		}
	}

	// Insert remaining tuples after channel closes
	if len(currentBatch) > 0 {

		if err := insertTuples(currentBatch, s.conn, s.table, headers); err != nil {
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
		row := marshal.ExtractTableRow(tuple)
		values = append(values, row)
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
