// Package filter provides helper functions for providing a Filter, defining
// the perimeter of the import
package filter

import (
	"errors"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
)

// filterReader defines the interface for reading SIREN filters from various sources.
// This is a private interface used internally by the filter package.
type filterReader interface {
	Read() (engine.SirenFilter, error)

	// SuccessStr returns a string to be displayed to the user in case of success
	SuccessStr() string
}

// Provider implements engine.FilterProvider.
// It retrieves SIREN filters using a priority-based approach:
// 1. Batch filter file (if available)
// 2. Database (fallback)
type Provider struct {
	DB db.Pool
}

// Get implements engine.FilterProvider.
func (p *Provider) Get(batch *engine.AdminBatch) (engine.SirenFilter, error) {
	filterFile, _ := batch.Files.GetFilterFile()

	readers := []filterReader{
		&CsvReader{filterFile},
		&DBReader{p.DB},
	}

	return getFromReaders(readers)
}

// Get retrieves the SIREN filter using a priority-based approach:
// 1. Batch filter file (if available)
// 2. Database (fallback)
//
// This is a convenience wrapper that uses default dependencies.
// Deprecated: Use Provider.Get() with dependency injection instead.
func Get(batch *engine.AdminBatch) (engine.SirenFilter, error) {
	p := &Provider{DB: db.DB}
	return p.Get(batch)
}

// getFromReaders tries each reader in order until one succeeds.
// The first successful filter is cached and returned.
func getFromReaders(readers []filterReader) (engine.SirenFilter, error) {
	var filter engine.SirenFilter
	var lastErr error

	for _, reader := range readers {
		var err error
		filter, err = reader.Read()

		if err != nil {
			// try next source
			lastErr = err
			continue
		}

		if filter != nil {
			slog.Debug(reader.SuccessStr())
			return filter, nil
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to load filter: %w", lastErr)
	}

	return nil, errors.New("no filter found from any source")
}
