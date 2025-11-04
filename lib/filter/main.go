package filter

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"opensignauxfaibles/lib/engine"
	"os"
	"regexp"
	"strconv"
)

// Usage: $ ./create_filter --path testData/test_data.csv

// DefaultNbMois is the default number of the most recent months during which the effectif of the company must reach the threshold.
const DefaultNbMois = 100

// DefaultMinEffectif is the default effectif threshold, expressed in number of employees.
const DefaultMinEffectif = 10

// DefaultNbIgnoredCols is the default number of rightmost columns that don't contain effectif data.
const DefaultNbIgnoredCols = 2

// NbLeadingColsToSkip is the number of leftmost columns that don't contain effectif data.
const NbLeadingColsToSkip = 5 // column names: "compte", "siret", "rais_soc", "ape_ins" and "dep"

// Implementation of the create_filter command.
func main() {

	var path = flag.String("path", "", "Chemin d'accès au fichier effectif")
	var nbMois = flag.Int(
		"nbMois",
		DefaultNbMois,
		"Nombre de mois observés (avec effectif connu) pour déterminer si l'entreprise dépasse 10 salariés",
	)
	var minEffectif = flag.Int(
		"minEffectif",
		DefaultMinEffectif,
		"Si une entreprise atteint ou dépasse 'minEffectif' dans les 'nbMois' derniers mois, elle est inclue dans le périmètre du filtre.",
	)
	var nIgnoredCols = flag.Int(
		"nIgnoredCols",
		DefaultNbIgnoredCols,
		"Nombre de colonnes à ignorer à la fin du fichier effectif",
	)
	flag.Parse()

	// create filter
	filter, err := Create(engine.NewBatchFile(*path), *nbMois, *minEffectif, *nIgnoredCols)
	if err != nil {
		log.Panic(err)
	}

	// write filter
	csvWriter := NewCsvWriter(os.Stdout)
	if err := csvWriter.Write(filter); err != nil {
		log.Panic(err)
	}
}

// Create generates a "filter" from an "effectif" file.
func Create(effectifFile engine.BatchFile, nbMois, minEffectif int, nIgnoredCols int) (engine.SirenFilter, error) {
	last, err := guessLastNMissing(effectifFile, nIgnoredCols)
	if err != nil {
		return nil, err
	}

	perimeter, err := getImportPerimeter(effectifFile, nbMois, minEffectif, nIgnoredCols+last)
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
// - OR an "effectif" file is provided.
//
// If a nil interface is provided fails.
// Note however that a nil *Reader pointer is properly handled and accepted.
func Check(r engine.FilterReader, batchFiles engine.BatchFiles) error {
	var err error

	effectifFile := batchFiles.GetEffectifFile()

	if r == nil {
		return errors.New("please provided a supported filter : nil interface is not supported")
	}

	// check if a filter can be read
	_, err = r.Read()

	validFiltering := (err == nil || effectifFile != nil)

	if !validFiltering {
		return errors.New("filter is missing: a filter or one effectif file should be provided")
	} else {
		slog.Debug("filter can be retrieved or created from effectif file")
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
func UpdateState(w engine.FilterWriter, batchFiles engine.BatchFiles) error {
	// Guard clause 1: the import filter is based uniquely on the effectif file.
	// If no effectif file is provided, there is nothing to update.
	effectifFile := batchFiles.GetEffectifFile()

	if effectifFile == nil {
		return nil
	}

	// Guard clause 2: Check if filter has been explicitely provided in the batch
	// In this case, we do not update the filter state.
	filterFile := batchFiles.GetFilterFile()
	filterIsExplicit := (filterFile != nil)

	if filterIsExplicit {
		return nil
	}

	// Guard clause 3: if no writer is provided, don't update
	if w == nil {
		slog.Debug("No filter writer provided, filter is not updated")
		return nil
	}

	slog.Debug("Writing filter file")

	// Create the filter
	sirenFilter, err := Create(
		effectifFile, // input: the effectif file
		DefaultNbMois,
		DefaultMinEffectif,
		DefaultNbIgnoredCols,
	)

	if err != nil {
		return err
	}

	// Write the filter
	return w.Write(sirenFilter)
}

func newCsvReader(reader io.Reader) *csv.Reader {
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	r.Comma = ';'
	return r
}

// getImportPerimeter makes a perimeter on effectif criterias alone
// This perimeter is used for efficient imports, and is further refined with
// SQL for the "clean_data" layer
func getImportPerimeter(effectifFile engine.BatchFile, nbMois, minEffectif, nIgnoredCols int) (map[string]struct{}, error) {
	f, err := effectifFile.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := newCsvReader(f)

	detectedSirens := map[string]struct{}{} // smaller memory footprint than map[string]bool
	if _, err = r.Read(); err != nil {      // en tête
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
		siret := record[1]
		shouldKeep := len(siret) == 14 &&
			isInsidePerimeter(record[NbLeadingColsToSkip:len(record)-nIgnoredCols], nbMois, minEffectif)

		var siren string
		if len(siret) >= 9 {
			siren = siret[0:9] // trim siret into a siren
			_, alreadyDetected := detectedSirens[siren]
			if shouldKeep && !alreadyDetected {
				detectedSirens[siren] = struct{}{}
			}
		} else {
			skippedLines++
			fmt.Printf("%d digit siret encountered, skipping line %d \n", len(siret), lineNumber)
		}
	}
	if skippedLines > 0 {
		fmt.Printf("%d lines with bad siret/siren skipped :( \n", skippedLines)
	}
	return detectedSirens, nil
}

func isInsidePerimeter(record []string, nbMois, minEffectif int) bool {
	for i := len(record) - 1; i >= len(record)-nbMois && i >= 0; i-- {
		if record[i] == "" {
			continue
		}
		reg, err := regexp.Compile("[^0-9]")
		if err != nil {
			log.Fatal(err)
		}
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
