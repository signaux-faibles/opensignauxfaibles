package engine

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"slices"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BatchSize controls the max number of rows inserted at a time
const BatchSize = 1000

// A PostgresSinkFactory create sinks that send data to the relevant
// postgresql table.
// At the end of all processes, it updates the materialized tables if needed.
// It implements the `SinkFactory` and `Finalizer` interfaces.
type PostgresSinkFactory struct {
	conn *pgxpool.Pool

	// Keys are the materialized viewsToRefresh to update
	// Values are specific parser types the viewsToRefresh depend on.
	// The materialized view will only be updated if any data has been imported
	// for this parser type.
	viewsToRefresh map[string][]base.ParserType

	// Parser Types for which a sink has been created
	// Use mutex for concurrent write access
	parserSinks []base.ParserType
	mu          sync.Mutex
}

func NewPostgresSinkFactory(conn *pgxpool.Pool, views map[string][]base.ParserType) SinkFactory {
	return &PostgresSinkFactory{conn: conn, viewsToRefresh: views}
}

func (f *PostgresSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	f.mu.Lock()
	f.parserSinks = append(f.parserSinks, parserType)
	f.mu.Unlock()

	switch parserType {
	case base.Apconso,
		base.Apdemande,
		base.Cotisation,
		base.Debit,
		base.Delai,
		base.Effectif,
		base.EffectifEnt,
		base.Sirene,
		base.SireneUl:

		tableName := fmt.Sprintf("stg_%s", parserType)
		return &PostgresSink{f.conn, tableName}, nil
	}

	return &DiscardDataSink{}, nil
}

func (f *PostgresSinkFactory) Finalize() error {
	for view, parserTypes := range f.viewsToRefresh {
		for _, parser := range parserTypes {
			if slices.Contains(f.parserSinks, parser) {
				// Data has been updated for a parser type on which the view depends
				// => udpdate the view
				_, err := f.conn.Exec(context.Background(), fmt.Sprintf("REFRESH MATERIALIZED VIEW %s", view))
				if err != nil {
					return fmt.Errorf("failed to refresh materialized view %s: %w", view, err)
				}

				slog.Info("Materialized View updated", "view", view)

				break
			}
		}
	}

	return nil
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

func (s *PostgresSink) ProcessOutput(ctx context.Context, ch chan marshal.Tuple) error {
	logger := slog.With("sink", "postgresql", "table", s.table)

	logger.Debug("stream output to PostgreSQL")

	var currentBatch []marshal.Tuple
	var headers []string

	logger.Debug("truncate table")

	_, err := s.conn.Exec(ctx, fmt.Sprintf("TRUNCATE %s", s.table))
	if err != nil {
		return fmt.Errorf("failed to truncate table: %w", err)
	}

	logger.Debug("data insertion")

	nInserted := 0

	for tuple := range ch {
		if headers == nil {
			headers = marshal.ExtractTableColumns(tuple)
		}

		currentBatch = append(currentBatch, tuple)

		if len(currentBatch) >= BatchSize {

			if err = insertTuples(currentBatch, s.conn, s.table, headers); err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}

			nInserted += len(currentBatch)

			currentBatch = currentBatch[:0] // Reset currentBatch slice
		}
	}

	// Insert remaining tuples after channel closes
	if len(currentBatch) > 0 {

		if err = insertTuples(currentBatch, s.conn, s.table, headers); err != nil {
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
