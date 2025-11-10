package filter

import (
	"errors"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"reflect"
)

// StandardReader implements engine.FilterReader
// It retrieves SIREN filters using a priority-based approach:
// 1. Batch filter file (if available, e.g. provided by user)
// 2. Database "stg_filter_import" table
type StandardReader struct {
	Batch *engine.AdminBatch
	DB    db.Pool
}

// Read implements engine.FilterReader
func (r *StandardReader) Read() (engine.SirenFilter, error) {
	if r == nil {
		return engine.NoFilter, nil
	}

	filterFile := r.Batch.Files.GetFilterFile()

	readers := []Reader{
		&CsvReader{filterFile},
		&DBReader{r.DB, db.TableStgFilterImport},
	}

	return trySeveralReaders(readers)
}

// trySeveralReaders tries each reader in order until one succeeds.
// The first successful filter is returned.
func trySeveralReaders(readers []Reader) (engine.SirenFilter, error) {
	var filter engine.SirenFilter
	var lastErr error

	for _, reader := range readers {
		var err error
		filter, err = reader.Read()

		if err != nil {
			// try next source
			slog.Debug("filter reader attempt failed", "reader_type", reflect.TypeOf(reader).String(), "error", err)
			lastErr = err
			continue
		}

		if filter != nil {
			return filter, nil
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to load filter: %w", lastErr)
	}

	return nil, errors.New("no filter found from any source")
}
