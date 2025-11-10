package filter

import (
	"errors"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/engine"
)

// StandardFilterResolver performs full filter resolution with Check, Update, and Read.
// It orchestrates the complete filter lifecycle: checking requirements, updating state
// based on effectif files, and loading the final filter.
type StandardFilterResolver struct {
	Reader Reader
	Writer Writer
}

// Resolve implements engine.FilterResolver.
// It performs the complete filter resolution workflow:
// 1. Check if filtering requirements are met (filter exists or effectif file present)
// 2. Update filter state in database if appropriate (effectif file provided, no explicit filter)
// 3. Read and return the final SirenFilter
func (r *StandardFilterResolver) Resolve(files engine.BatchFiles) (engine.SirenFilter, error) {
	slog.Info("resolving filter...")

	// Check filter requirements
	if err := Check(r.Reader, files); err != nil {
		return nil, fmt.Errorf("filter check failed: %w", err)
	}
	slog.Debug("filter check passed: conditions are met for filtering input data")

	// Update filter state if needed
	if err := UpdateState(r.Writer, files); err != nil {
		return nil, fmt.Errorf("filter update failed: %w", err)
	}

	// Read the filter
	sirenFilter, err := r.Reader.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to read filter: %w", err)
	}
	if sirenFilter == nil {
		return nil, errors.New(`filter is required but missing.
When the filter is missing, it must be initialized by importing an 'effectif' file,
or by placing a filter file (prefixed with 'filter_') in the data import directory.
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
