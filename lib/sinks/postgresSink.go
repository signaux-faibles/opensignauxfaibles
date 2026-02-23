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
const BatchSize = 100000

// MaterializedViewsWorkMem is the value of Postgresql's WORK_MEM option to set locally for materialized views updates.
const MaterializedViewsWorkMem = "512MB"

// MaintenanceWorkMem is the value of Postgresql's MAINTENANCE_WORK_MEM option to set locally for index recreation.
const MaintenanceWorkMem = "512MB"

// tmpTableSavedIndexes is the name of the temporary table to store indexes
const tmpTableSavedIndexes = "tmp_saved_indexes"

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
		engine.SireneUl,
		engine.Procol,
		engine.SireneHisto:

		tableName := fmt.Sprintf("stg_%s", parserType)
		var viewsToRefresh []string

		switch parserType {
		case engine.Apdemande:
			viewsToRefresh = []string{db.ViewStgApdemandePeriod, db.ViewCleanAp}
		case engine.Apconso:
			viewsToRefresh = []string{db.ViewCleanAp}
		case engine.SireneUl:
			viewsToRefresh = []string{db.ViewSirenBlacklist}
		case engine.Effectif:
			viewsToRefresh = []string{db.ViewSirenBlacklist}
		case engine.Debit:
			viewsToRefresh = []string{db.IntermediateViewDebits, db.ViewDebit}
		}

		return &PostgresSink{f.conn, tableName, viewsToRefresh}, nil
	}

	slog.Warn("type de parser non supporté pour envoi des données à PostgreSQL", "parser", parserType)

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

	// Names of materialized views to refresh after write, in order
	viewsToRefresh []string
}

type indexInfo struct {
	IndexName string `db:"indexname"`
	IndexDef  string `db:"indexdef"`
}

// recreateSavedIndexes reads index definitions from the saved_indexes table
// and recreates them. If all indexes are recreated successfully, their entries
// are removed from the table. Otherwise, entries are kept for the next run.
func (s *PostgresSink) recreateSavedIndexes(ctx context.Context, logger *slog.Logger) {
	savedRows, err := s.conn.Query(ctx, fmt.Sprintf(
		"SELECT index_name AS indexname, index_def AS indexdef FROM %s WHERE table_name = $1",
		tmpTableSavedIndexes,
	), s.table)
	if err != nil {
		logger.Error("failed to read saved indexes", "error", err)
		return
	}
	savedIndexes, err := pgx.CollectRows(savedRows, pgx.RowToStructByName[indexInfo])
	if err != nil {
		logger.Error("failed to collect saved indexes", "error", err)
		return
	}

	if len(savedIndexes) == 0 {
		return
	}

	logger.Info("recreating indexes", "count", len(savedIndexes))

	allRecreated := true
	for _, idx := range savedIndexes {
		tx, err := s.conn.Begin(ctx)
		if err != nil {
			logger.Error("failed to begin transaction for index recreation", "error", err)
			allRecreated = false
			continue
		}

		_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL maintenance_work_mem = '%s'", MaintenanceWorkMem))
		if err != nil {
			logger.Error("failed to set maintenance_work_mem", "error", err)
			tx.Rollback(ctx) //nolint:errcheck
			allRecreated = false
			continue
		}

		idempotentDef := strings.Replace(idx.IndexDef, "CREATE INDEX ", "CREATE INDEX IF NOT EXISTS ", 1)
		_, err = tx.Exec(ctx, idempotentDef)
		if err != nil {
			logger.Error("failed to recreate index", "index", idx.IndexName, "error", err)
			tx.Rollback(ctx) //nolint:errcheck
			allRecreated = false
			continue
		}

		if err = tx.Commit(ctx); err != nil {
			logger.Error("failed to commit index recreation transaction", "error", err)
			allRecreated = false
			continue
		}
		logger.Debug("index recreated", "index", idx.IndexName)
	}

	if allRecreated {
		logger.Info("all indexes recreated successfully")
		_, err = s.conn.Exec(ctx, fmt.Sprintf(
			"DELETE FROM %s WHERE table_name = $1",
			tmpTableSavedIndexes,
		), s.table)
		if err != nil {
			logger.Error("failed to clean up saved indexes", "error", err)
		}
	} else {
		logger.Warn("some indexes could not be recreated, keeping saved_indexes entries for next run")
	}

	_, err = s.conn.Exec(ctx, fmt.Sprintf("ANALYZE %s", s.table))
	if err != nil {
		logger.Error("failed to ANALYZE table", "error", err)
	}
	logger.Debug("Table ANALYZEd")
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

	logger.Debug("setup table, drop indexes")

	// For performance reasons, we drop the indexes and recreate them after bulk
	// import

	// We want indexes but NOT primary keys
	rows, err := s.conn.Query(ctx, `
		SELECT i.indexname, i.indexdef
		FROM pg_indexes i
		JOIN pg_class c ON c.relname = i.indexname
		JOIN pg_index idx ON idx.indexrelid = c.oid
		WHERE i.tablename = $1
		AND i.schemaname = current_schema()
		AND NOT idx.indisprimary
	`, s.table)
	if err != nil {
		return fmt.Errorf("failed to retrieve indexes: %w", err)
	}

	indexes, err := pgx.CollectRows(rows, pgx.RowToStructByName[indexInfo])
	if err != nil {
		return fmt.Errorf("failed to collect indexes: %w", err)
	}

	// Save indexes to database so they survive a process crash
	for _, idx := range indexes {
		_, err = s.conn.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (table_name, index_name, index_def) VALUES ($1, $2, $3) ON CONFLICT (table_name, index_name) DO UPDATE SET index_def = EXCLUDED.index_def",
			tmpTableSavedIndexes,
		), s.table, idx.IndexName, idx.IndexDef)
		if err != nil {
			return fmt.Errorf("failed to save index %s to %s: %w", idx.IndexName, tmpTableSavedIndexes, err)
		}
	}

	// Recreate indexes even if an error occurred
	defer s.recreateSavedIndexes(ctx, logger)

	// Dropper les indexes avant l'import en masse
	for _, idx := range indexes {
		logger.Info("dropping index", "index", idx.IndexName)
		_, err = s.conn.Exec(ctx, fmt.Sprintf("DROP INDEX IF EXISTS %s", idx.IndexName))
		if err != nil {
			return fmt.Errorf("failed to drop index %s: %w", idx.IndexName, err)
		}
	}

	logger.Debug("data insertion")

	nInserted := 0

	for tuple := range ch {
		if headers == nil {
			headers = ExtractTableColumns(tuple)
		}

		currentBatch = append(currentBatch, tuple)

		if len(currentBatch) >= BatchSize {

			if err = insertTuples(ctx, currentBatch, s.conn, s.table, headers); err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}

			nInserted += len(currentBatch)

			currentBatch = currentBatch[:0] // Reset currentBatch slice
		}
	}

	// Insert remaining tuples after channel closes
	if len(currentBatch) > 0 {

		if err = insertTuples(ctx, currentBatch, s.conn, s.table, headers); err != nil {
			return fmt.Errorf("failed to execute final batch: %w", err)
		}

		nInserted += len(currentBatch)
	}

	logger.Info("output streaming to PostgreSQL ended successfully", "n_inserted", nInserted)
	if len(s.viewsToRefresh) > 0 {
		logger.Info("update materialized views, some views may take some time... (up to an hour or more)", "views", s.viewsToRefresh)

		for _, view := range s.viewsToRefresh {
			_, err = s.conn.Exec(ctx, fmt.Sprintf(`
      BEGIN;
      SET LOCAL work_mem = '%s';
      REFRESH MATERIALIZED VIEW %s;
      COMMIT;
      `, MaterializedViewsWorkMem, view))
			if err != nil {
				return fmt.Errorf("failed to refresh materialized view %s: %w", view, err)
			}

			logger.Debug("materialized view updated", "view", view)
		}
		logger.Info("materialized view update ended successfully")
	}

	return nil
}

func insertTuples(ctx context.Context, tuples []engine.Tuple, conn db.Pool, tableName string, columns []string) error {
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
		ctx,
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
