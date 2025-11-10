package sinks

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"strings"

	"github.com/jackc/pgx/v5"
)

// BatchSize controls the max number of rows inserted at a time
const BatchSize = 1000

type PostgresSinkFactory struct {
	conn db.Pool
}

func NewPostgresSinkFactory(conn db.Pool) engine.SinkFactory {
	return &PostgresSinkFactory{conn}
}

func (f *PostgresSinkFactory) CreateSink(parserType engine.ParserType) (engine.DataSink, error) {
	switch parserType {
	case engine.Apconso,
		engine.Apdemande,
		engine.Cotisation,
		engine.Debit,
		engine.Delai,
		engine.Effectif,
		engine.EffectifEnt,
		engine.Sirene,
		engine.SireneUl:

		tableName := fmt.Sprintf("stg_%s", parserType)
		materializedTableUpdate := ""

		switch parserType {
		case engine.Apdemande:
			materializedTableUpdate = db.ViewStgApdemandePeriod
		case engine.SireneUl:
			materializedTableUpdate = db.ViewCleanFilter
		case engine.Effectif:
			materializedTableUpdate = db.ViewCleanFilter
		}

		return &PostgresSink{f.conn, tableName, materializedTableUpdate}, nil
	}

	return &engine.DiscardDataSink{}, nil
}

// PostgresSink writes the output to postgresql.
//
// The name of the table is defined by `parserType`, prefixed with "stg_".
// The table is expected to exist and be properly formatted. It is truncated
// before inserting new values.
//
// The columns of this table will correspond to the "Tuple.Headers()"
type PostgresSink struct {
	conn db.Pool

	// Name of the table to which to write
	table string

	// Name of a materialized view to refresh after write
	viewToRefresh string
}

func (s *PostgresSink) ProcessOutput(ctx context.Context, ch chan engine.Tuple) error {
	logger := slog.With("sink", "postgresql", "table", s.table)

	logger.Info("streaming output to PostgreSQL...")

	var currentBatch []engine.Tuple
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
			headers = ExtractTableColumns(tuple)
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

	logger.Info("output streaming to PostgreSQL ended successfully", "n_inserted", nInserted)

	if s.viewToRefresh != "" {
		_, err = s.conn.Exec(ctx, fmt.Sprintf("REFRESH MATERIALIZED VIEW %s", s.viewToRefresh))
		if err != nil {
			return fmt.Errorf("failed to refresh materialized view %s: %w", s.viewToRefresh, err)
		}

		logger.Debug("materialized view updated", "view", s.viewToRefresh)
	}

	return nil
}

func insertTuples(tuples []engine.Tuple, conn db.Pool, tableName string, columns []string) error {
	if len(tuples) == 0 {
		return nil
	}

	values := make([][]any, 0, len(tuples))

	// TODO rather than construct values
	// implement CopyFromSource interface
	for _, tuple := range tuples {
		row := ExtractTableRow(tuple)
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

// ExtractTableColumns extrait les noms des colonnes pour une table SQL via le tag "sql"
func ExtractTableColumns(tuple engine.Tuple) (header []string) {
	return parsing.ExtractFieldsByTags(tuple, "sql")
}

// ExtractTableRow extrait les valeurs des colonnes pour une table SQL via le tag "sql"
func ExtractTableRow(tuple engine.Tuple) (row []any) {
	rawValues := parsing.ExtractValuesByTags(tuple, "sql")
	for _, v := range rawValues {
		row = append(row, v.Interface())
	}
	return row
}
