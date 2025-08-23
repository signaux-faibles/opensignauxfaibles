// Ce fichier est responsable de collecter les messages et de les ajouter
// dans la collection Journal.

package engine

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"opensignauxfaibles/lib/marshal"
)

const EventTable string = "import_logs"

type EventSink interface {
	Process(ch chan marshal.Event) error
}

type PostgresEventSink struct {
	conn *pgxpool.Pool

	// Name of the table to which to write
	table string

	// Attach the processed events to a specific command
	command string

	// time of creation
	timestamp time.Time
}

func NewPostgresEventSink(conn *pgxpool.Pool, command string) EventSink {
	return &PostgresEventSink{conn, EventTable, command, time.Now()}
}

func (s *PostgresEventSink) Process(ch chan marshal.Event) error {
	logger := slog.With("table", s.table)

	ctx := context.Background()
	logger.Debug("stream events to PostgreSQL")

	logger.Debug("begin transaction")
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx) // Will be no-op if tx is already committed

	logger.Debug("events insertion")

	nInserted := 0

	var currentBatch []marshal.Event
	var headers []string

	for event := range ch {
		if headers == nil {
			headers = []string{"start_date", "parser", "batch_key", "head_fatal", "head_rejected", "is_fatal", "lines_parsed", "lines_rejected", "lines_skipped", "lines_valid", "summary"}
		}

		currentBatch = append(currentBatch, event)

		if len(currentBatch) >= BatchSize {

			if err := insertEvents(currentBatch, tx, s.table, headers); err != nil {
				return fmt.Errorf("failed to execute batch insert: %w", err)
			}

			nInserted += len(currentBatch)

			currentBatch = currentBatch[:0] // Reset currentBatch slice
		}
	}

	// Insert remaining tuples after channel closes
	if len(currentBatch) > 0 {

		if err := insertEvents(currentBatch, tx, s.table, headers); err != nil {
			return fmt.Errorf("failed to execute final batch: %w", err)
		}

		nInserted += len(currentBatch)
	}

	tx.Commit(context.Background())

	logger.Debug("output streaming to PostgreSQL ended successfully", "n_inserted", nInserted)

	return nil
}

func insertEvents(events []marshal.Event, tx pgx.Tx, tableName string, columns []string) error {
	if len(events) == 0 {
		return nil
	}

	values := make([][]any, 0, len(events))

	// TODO rather than construct values
	// implement CopyFromSource interface
	for _, event := range events {
		row := []any{
			event.StartDate,
			event.Parser,
			event.Report.BatchKey,
			// event.Report.HeadFatal,
			// event.Report.HeadRejected,
			"abc",
			"def",
			event.Report.IsFatal,
			event.Report.LinesParsed,
			event.Report.LinesRejected,
			event.Report.LinesSkipped,
			event.Report.LinesValid,
			event.Report.Summary,
		}
		values = append(values, row)
	}
	lowerColumns := make([]string, len(columns))
	for i, c := range columns {
		lowerColumns[i] = strings.ToLower(c)
	}

	// Batch insertion
	_, err := tx.CopyFrom(
		context.Background(),
		pgx.Identifier{tableName},
		lowerColumns,
		pgx.CopyFromRows(values),
	)

	return err
}
