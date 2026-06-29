package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// RefreshUnpopulatedMatviews refreshes every materialized view in the current
// schema that is still in its initial `WITH NO DATA` state.
//
// A `SELECT` (including the one used by `COPY ... TO STDOUT`) on such a view
// fails with SQLSTATE 55000. We must therefore populate them before any code
// path that reads them — typically before `export`.
//
// An unpopulated MV may itself depend on another unpopulated MV (e.g.
// `clean_debit` -> `stg_tmp_debits_simplified`), so we loop until either all
// are populated or no progress is made in a pass.
func RefreshUnpopulatedMatviews(ctx context.Context, pool Pool) error {
	for {
		names, err := listUnpopulatedMatviews(ctx, pool)
		if err != nil {
			return fmt.Errorf("failed to list unpopulated matviews: %w", err)
		}
		if len(names) == 0 {
			return nil
		}

		var refreshed int
		var lastErr error
		for _, name := range names {
			_, err := pool.Exec(ctx, fmt.Sprintf("REFRESH MATERIALIZED VIEW %s", pgx.Identifier{name}.Sanitize()))
			if err != nil {
				lastErr = err
				continue
			}
			slog.Info("refreshed unpopulated materialized view", "view", name)
			refreshed++
		}

		if refreshed == 0 {
			return fmt.Errorf("could not refresh remaining unpopulated matviews %v: %w", names, lastErr)
		}
	}
}

func listUnpopulatedMatviews(ctx context.Context, pool Pool) ([]string, error) {
	rows, err := pool.Query(ctx, `
		SELECT matviewname
		FROM pg_matviews
		WHERE schemaname = current_schema()
		  AND ispopulated = false
		ORDER BY matviewname`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}
