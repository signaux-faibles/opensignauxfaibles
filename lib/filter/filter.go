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

// Get retrieves the SIREN filter using a priority-based approach:
// 1. Batch filter file (if available)
// 2. Database (fallback)
//
// This is a convenience wrapper that uses default dependencies.
func Get(batch *engine.AdminBatch) (engine.SirenFilter, error) {
	filterFile, _ := batch.Files.GetFilterFile()

	readers := []engine.FilterReader{
		&CsvReader{filterFile},
		&DBReader{db.DB},
	}

	return GetFromReaders(readers)
}

// GetFromReaders tries each reader in order until one succeeds.
// The first successful filter is cached and returned.
func GetFromReaders(readers []engine.FilterReader) (engine.SirenFilter, error) {
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
