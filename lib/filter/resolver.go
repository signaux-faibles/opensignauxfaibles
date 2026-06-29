package filter

import (
	"errors"
	"fmt"
	"log/slog"

	"opensignauxfaibles/lib/engine"
)

// StandardFilterResolver reads an existing filter from the database or batch files.
// It does NOT compute or update the perimeter — use the "computePerimeter" command for that.
type StandardFilterResolver struct {
	Reader Reader
}

// Resolve implements engine.FilterResolver.
// It checks that a filter exists, then reads and returns it.
func (r *StandardFilterResolver) Resolve(files engine.BatchFiles) (engine.SirenFilter, error) {
	slog.Info("resolving filter...")

	if err := CheckFilterExists(r.Reader); err != nil {
		return nil, fmt.Errorf("filter check failed: %w", err)
	}
	slog.Debug("filter check passed: a filter is available")

	// Read the filter
	sirenFilter, err := r.Reader.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to read filter: %w", err)
	}
	if sirenFilter == nil {
		return nil, errors.New(`filter is required but missing.
Run "computePerimeter" to compute the perimeter from an effectif_ent file,
or place a filter file (prefixed with 'filter_') in the data import directory.
If you wish to import without a filter, use the "--no-filter" option`)
	}

	slog.Info("filter resolution ended successfully")
	return sirenFilter, nil
}

// NoFilterResolver bypasses all filtering logic and returns NoFilter directly.
type NoFilterResolver struct{}

// Resolve implements engine.FilterResolver.
// It immediately returns NoFilter without performing any checks or updates.
func (r *NoFilterResolver) Resolve(_ engine.BatchFiles) (engine.SirenFilter, error) {
	slog.Info("no-filter mode: all data will be imported without filtering")
	return engine.NoFilter, nil
}
