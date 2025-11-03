package filter

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"

	"github.com/jackc/pgx/v5"
)

const (
	TableName = "stg_filter_import"
)

// DBWriter writes the filter to the database table.
// Existing data is truncated before inserting new data.
type DBWriter struct {
	DB db.Pool
}

// Write implements engine.FilterWriter
func (f *DBWriter) Write(filter engine.SirenFilter) error {
	ctx := context.Background()

	// Truncate existing data
	_, err := f.DB.Exec(ctx, fmt.Sprintf("TRUNCATE %s", TableName))
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %w", TableName, err)
	}

	slog.Debug("Truncated filter table", "table", TableName)

	// Get all SIRENs from the filter
	sirens := filter.All()
	if len(sirens) == 0 {
		slog.Warn("No SIRENs in filter to write", "table", TableName)
		return nil
	}

	// Prepare values for bulk insert
	values := make([][]any, 0, len(sirens))
	for siren := range sirens {
		values = append(values, []any{siren})
	}

	// Bulk insert using CopyFrom
	_, err = f.DB.CopyFrom(
		ctx,
		pgx.Identifier{TableName},
		[]string{"siren"},
		pgx.CopyFromRows(values),
	)
	if err != nil {
		return fmt.Errorf("failed to insert into table %s: %w", TableName, err)
	}

	slog.Debug("Filter written to DB", "table", TableName, "n_sirens", len(sirens))
	return nil
}
