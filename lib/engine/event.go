// Ce fichier est responsable de collecter les messages et de les ajouter
// dans la collection Journal.

package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ReportTable string = "import_logs"

type ReportSink interface {
	Process(ch chan Report) error
}

type PostgresReportSink struct {
	conn *pgxpool.Pool

	// Name of the table to which to write
	table string
}

func NewPostgresReportSink(conn *pgxpool.Pool) ReportSink {
	return &PostgresReportSink{conn, ReportTable}
}

func (s *PostgresReportSink) Process(ch chan Report) error {
	logger := slog.With("sink", "postgresql", "table", s.table)

	logger.Debug("stream reports/logs to PostgreSQL")

	logger.Debug("reports insertion")
	nInserted := 0

	for report := range ch {

		if err := insertReport(report, s.conn, s.table); err != nil {
			return fmt.Errorf("failed to insert report: %w", err)
		}

		nInserted++
	}

	logger.Debug("reports streaming to PostgreSQL ended successfully", "n_inserted", nInserted)

	return nil
}

func insertReport(report Report, conn *pgxpool.Pool, tableName string) error {

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
		report.StartDate,
		report.Parser,
		report.BatchKey,
		toPgArray(report.HeadFatal),
		toPgArray(report.HeadRejected),
		report.IsFatal,
		report.LinesParsed,
		report.LinesRejected,
		report.LinesSkipped,
		report.LinesValid,
		report.Summary,
	}
	ctx := context.Background()
	_, err := conn.Exec(ctx, query, row...)
	return err
}

type StdoutReportSink struct{}

func (s *StdoutReportSink) Process(ch chan Report) error {
	for report := range ch {
		jsonReport, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		slog.Info(string(jsonReport))
	}
	return nil
}
