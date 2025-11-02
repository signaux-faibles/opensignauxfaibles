package filter

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
)

// DBReader reads the filter from the database "filter" table.
type DBReader struct {
	Conn      db.Pool
	TableName string
}

func (f *DBReader) Read() (engine.SirenFilter, error) {
	var filter = make(MapFilter)

	rows, err := f.Conn.Query(context.Background(), fmt.Sprintf("SELECT siren FROM %s", f.TableName))
	if err != nil {
		return nil, fmt.Errorf("error retrieving \"filter\" from DB, query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var siren string
		if err := rows.Scan(&siren); err != nil {
			return nil, fmt.Errorf("error reading \"filter\" from DB, scan failed: %w", err)
		}
		filter[siren] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading \"filter\" from DB, rows iteration failed: %w", err)
	}

	if len(filter) == 0 {
		return nil, fmt.Errorf("error reading \"filter\" from DB: table %s is empty", f.TableName)
	}

	slog.Debug("Filter retrieved from DB")
	return filter, nil
}
