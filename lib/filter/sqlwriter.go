package filter

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"

	"github.com/jackc/pgx/v5"
)

// DBWriter writes the filter to the database table.
// Existing data is truncated before inserting new data.
type DBWriter struct {
	Conn      db.Pool
	TableName string
}

// Write implements engine.FilterWriter
func (f *DBWriter) Write(filter engine.SirenFilter) error {
	ctx := context.Background()

	// Truncate existing data
	_, err := f.Conn.Exec(ctx, fmt.Sprintf("TRUNCATE %s", f.TableName))
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %w", f.TableName, err)
	}

	slog.Debug("Truncated filter table", "table", f.TableName)

	// Get all SIRENs from the filter
	sirens := filter.All()
	if len(sirens) == 0 {
		slog.Warn("No SIRENs in filter to write", "table", f.TableName)
		return nil
	}

	// Prepare values for bulk insert
	values := make([][]any, 0, len(sirens))
	for siren := range sirens {
		values = append(values, []any{siren})
	}

	// Bulk insert using CopyFrom
	_, err = f.Conn.CopyFrom(
		ctx,
		pgx.Identifier{f.TableName},
		[]string{"siren"},
		pgx.CopyFromRows(values),
	)
	if err != nil {
		return fmt.Errorf("failed to insert into table %s: %w", f.TableName, err)
	}

	slog.Debug("Filter written to DB", "table", f.TableName, "n_sirens", len(sirens))
	return nil
}
