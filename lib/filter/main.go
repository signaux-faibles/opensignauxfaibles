// Package filter manages the import perimeter for Signaux Faibles data.
//
// The perimeter is stored as state between successive imports, to avoid
// the requirement of importing files together every time.
//
// This package provides utilities to create and maintain SIREN filters that
// determine which companies should be included in the data import. Filters
// are typically derived from effectif_ent (employee count) data, selecting
// companies that meet minimum employee thresholds over a specified time
// period.
//
// Note that a subsequent more fine-grained filtering (e.g. on juridic nature)
// happens at a later stage, thanks to SQL queries, between the "stg_..." and
// the "clean_..." layers.
//
// The package provides functions to:
// - Create filters from effectif_ent files based on configurable criteria
// - Check if valid filtering conditions are met before import
// - Read filters from multiple sources (files, database). Filters provided as
// an explicit file have precedence over the database stored filter.
// - Update filter state in the database when appropriate (effectif_ent file
// is present, and no explicit filter has been provided in the batch).
package filter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"opensignauxfaibles/lib/engine"
	"regexp"
	"strconv"
)

// Writer writes a filter
// Implementations may write to e.g. a file, a database table.
type Writer interface {
	Write(engine.SirenFilter) error
}

// Reader retrieves a SirenFilter for a given batch.
// Implementations may read from files, databases, or other sources.
type Reader interface {
	Read() (engine.SirenFilter, error)
}

// DefaultNbMois is the default number of the most recent months during which the effectif of the company must reach the threshold.
const DefaultNbMois = 120

// DefaultMinEffectif is the default effectif threshold, expressed in number of employees.
const DefaultMinEffectif = 10

// DefaultNbIgnoredCols is the default number of rightmost columns that don't contain effectif data.
const DefaultNbIgnoredCols = 1 // column name: "rais_soc"

// NbLeadingColsToSkip is the number of leftmost columns that don't contain effectif data.
const NbLeadingColsToSkip = 1 // column name: "siren"

// Create generates a "filter" from an "effectif_ent" file.
func Create(effectifEntFile engine.BatchFile, nbMois, minEffectif int, nIgnoredCols int) (engine.SirenFilter, error) {
	last, err := guessLastNMissing(effectifEntFile, nIgnoredCols)
	if err != nil {
		return nil, err
	}

	perimeter, err := getImportPerimeter(effectifEntFile, nbMois, minEffectif, nIgnoredCols+last)
	if err != nil {
		return nil, err
	}

	// Convert to MapFilter
	mapFilter := make(MapFilter)
	for siren := range perimeter {
		mapFilter[siren] = true
	}

	return mapFilter, nil
}

// Check checks whether the conditions for filtering are met, as we
// do not want to import all data by accident.
//
// It checks whether :
// - a  non-empty filter can be read from the provided reader
// - OR an "effectif_ent" file is provided.
//
// If a nil interface is provided fails.
// Note however that a nil *Reader pointer is properly handled and accepted.
func Check(r Reader, batchFiles engine.BatchFiles) error {
	var err error

	effectifEntFile := batchFiles.GetEffectifEntFile()

	if r == nil {
		return errors.New("please provide a supported filter : nil interface is not supported")
	}

	// check if a filter can be read
	_, err = r.Read()

	validFiltering := (err == nil || effectifEntFile != nil)

	if !validFiltering {
		return errors.New("filter is missing: a filter or one effectif_ent file should be provided")
	} else {
		slog.Debug("filter can be retrieved or created from effectif_ent file")
	}

	return nil
}

// UpdateState udpates (or creates) the filter if appropriate.
// Providing a `nil` writer will result in no update.
//
// It updates (or creates if none exists) the filter if the following conditions are met :
// - An "effectif" file is provided
// - AND the filter is not explicitely provided in the batchfile
//
// The rationale behind this last point is that a user-provided filter is
// usually used solely for tests and should not affect the saved perimeter in
// the database.
func UpdateState(w Writer, batchFiles engine.BatchFiles) error {
	// Guard clause 1: the import filter is based uniquely on the effectif_ent file.
	// If no effectif_ent file is provided, there is nothing to update.
	effectifEntFile := batchFiles.GetEffectifEntFile()

	if effectifEntFile == nil {
		slog.Info("no effectif_ent file provided, filter is not updated")
		return nil
	}

	// Guard clause 2: Check if filter has been explicitely provided in the batch
	// In this case, we do not update the filter state.
	filterFile := batchFiles.GetFilterFile()
	filterIsExplicit := (filterFile != nil)

	if filterIsExplicit {
		slog.Info("explicit filter file provided, filter is not updated")
		return nil
	}

	// Guard clause 3: if no writer is provided, don't update
	if w == nil {
		slog.Warn("no filter writer provided, filter is not updated")
		return nil
	}

	slog.Info("update filter...")

	// Create the filter
	sirenFilter, err := Create(
		effectifEntFile,
		DefaultNbMois,
		DefaultMinEffectif,
		DefaultNbIgnoredCols,
	)

	if err != nil {
		return err
	}

	// Write the filter
	err = w.Write(sirenFilter)

	if err != nil {
		return err
	}

	slog.Info("updated filter written with success")
	return nil
}

func newCsvReader(reader io.Reader) *csv.Reader {
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	r.Comma = ';'
	return r
}

// getImportPerimeter makes a perimeter on effectif criterias alone
func getImportPerimeter(effectifEntFile engine.BatchFile, nbMois, minEffectif, nIgnoredCols int) (map[string]struct{}, error) {
	f, err := effectifEntFile.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := newCsvReader(f)

	detectedSirens := map[string]struct{}{} // smaller memory footprint than map[string]bool
	if _, err = r.Read(); err != nil {      // header
		return nil, err
	}

	lineNumber, skippedLines := 0, 0
	for {
		lineNumber++
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		siren := record[0]
		shouldKeep := len(siren) == 9 &&
			isInsidePerimeter(record[NbLeadingColsToSkip:len(record)-nIgnoredCols], nbMois, minEffectif)

		if len(siren) == 9 {
			_, alreadyDetected := detectedSirens[siren]
			if shouldKeep && !alreadyDetected {
				detectedSirens[siren] = struct{}{}
			}
		} else {
			skippedLines++
		}
	}
	if skippedLines > 0 {
		slog.Info(fmt.Sprintf("%d lines with bad siret/siren skipped in the effectif_ent file at filter creation", skippedLines))
	}
	return detectedSirens, nil
}

func isInsidePerimeter(record []string, nbMois, minEffectif int) bool {
	for i := len(record) - 1; i >= len(record)-nbMois && i >= 0; i-- {
		if record[i] == "" {
			continue
		}
		reg := regexp.MustCompile("[^0-9]")

		effectif, err := strconv.Atoi(reg.ReplaceAllString(record[i], ""))
		if err != nil {
			slog.Error(fmt.Sprintf("%v", record))
			log.Panic(err)
		}
		if effectif >= minEffectif {
			return true
		}
	}
	return false
}

// guessLastNMissingFromReader returns the number of rightmost columns
// (on top of nIgnoredCols columns) that never have a value.
func guessLastNMissing(file engine.BatchFile, nIgnoredCols int) (int, error) {
	f, err := file.Open()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := newCsvReader(f)

	if _, err = r.Read(); err != nil { // en tête
		return 0, err
	}

	var lastConsideredCol int // index of the rightmost column of the last read row
	lastColWithValue := -1    // index of the rightmost column which had a value at least once
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
		lastConsideredCol = len(record) - 1 - nIgnoredCols
		for i := lastConsideredCol; i > lastColWithValue; i-- {
			if record[i] != "" {
				lastColWithValue = i
			}
		}
	}
	return lastConsideredCol - lastColWithValue, nil
}
