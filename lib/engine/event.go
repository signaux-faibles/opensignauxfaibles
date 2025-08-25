// Ce fichier est responsable de collecter les messages et de les ajouter
// dans la collection Journal.

package engine

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	logger := slog.With("sink", "postgresql", "table", s.table)

	logger.Debug("stream events to PostgreSQL")

	logger.Debug("events insertion")
	nInserted := 0

	for event := range ch {

		if err := insertEvent(event, s.conn, s.table); err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}

		nInserted++
	}

	logger.Debug("events streaming to PostgreSQL ended successfully", "n_inserted", nInserted)

	return nil
}

func insertEvent(event marshal.Event, conn *pgxpool.Pool, tableName string) error {

	query := fmt.Sprintf(`
		INSERT INTO %s (
			start_date, parser, batch_key, head_fatal, head_rejected,
			is_fatal, lines_parsed, lines_rejected, lines_skipped,
			lines_valid, summary
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, tableName)

	// Auxiliaire pour insérer des données au format Postgres TEXT[]
	toPgArray := func(slice []string) pgtype.FlatArray[string] {
		if slice == nil {
			return pgtype.FlatArray[string](nil)
		}
		return pgtype.FlatArray[string](slice)
	}

	row := []any{
		event.StartDate,
		event.Parser,
		event.Report.BatchKey,
		toPgArray(event.Report.HeadFatal),
		toPgArray(event.Report.HeadRejected),
		event.Report.IsFatal,
		event.Report.LinesParsed,
		event.Report.LinesRejected,
		event.Report.LinesSkipped,
		event.Report.LinesValid,
		event.Report.Summary,
	}
	ctx := context.Background()
	_, err := conn.Exec(ctx, query, row...)
	return err
}
