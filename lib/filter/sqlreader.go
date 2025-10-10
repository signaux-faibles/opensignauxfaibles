package filter

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
)

// DBReader reads the filter from the database "filter" table.
type DBReader struct {
	Conn db.Pool
}

func (f *DBReader) Read() (engine.SirenFilter, error) {
	var filter = make(MapFilter)

	rows, err := f.Conn.Query(context.Background(), "SELECT siren FROM filter")
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

	return filter, nil
}

func (f *DBReader) SuccessStr() string {
	return "Filter retrieved from DB"
}
