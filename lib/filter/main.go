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
	"opensignauxfaibles/lib/sinks"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/viper"
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
const DefaultNbMois = 36

// DefaultMinEffectif is the default effectif threshold, expressed in number of employees.
const DefaultMinEffectif = 10

// DefaultNbMoisRecentData is the number of most recent months that must have at least one effectif data point.
const DefaultNbMoisRecentData = 3

// effColRegex is a regular expression used to extract data related to
// headcounts
var effColRegex = regexp.MustCompile(`^eff[0-9]+$`)

// sirenColumn is the name of the column that holds the SIREN number
const sirenColumn = "siren"

// FilterStats tracks the number of companies at each stage of the filtering process
type FilterStats struct {
	TotalCompanies              int
	AfterEffectifFilter         int
	AfterRecentDataFilter       int
	ExcludedByEffectifFilter    int
	ExcludedByRecentDataFilter  int
}

// Create generates a "filter" from an "effectif_ent" file.
func Create(effectifEntFile engine.BatchFile, nbMois, minEffectif int) (engine.SirenFilter, error) {
	extractor, err := newEffectifDataExtractor(effectifEntFile)
	if err != nil {
		return nil, err
	}

	perimeter, stats, err := getImportPerimeter(effectifEntFile, nbMois, minEffectif, extractor)
	if err != nil {
		return nil, err
	}

	// Generate detailed filter report
	if err := generateFilterReport(stats); err != nil {
		slog.Warn("failed to generate filter report", "error", err)
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

// newEffectifDataExtractor builds an effectifDataExtractor by reading the
// header of "file". It identifies the "siren" column and the effectif
// columns (matching "^eff[0-9]+$") from the header, then uses
// guessLastNMissing to detect and exclude trailing entirely-empty columns
func newEffectifDataExtractor(file engine.BatchFile) (*effectifDataExtractor, error) {
	header, err := readHeader(file)
	if err != nil {
		return nil, err
	}

	e := effectifDataExtractor{sirenColIndex: -1}

	for i, colName := range header {
		if colName == sirenColumn {
			e.sirenColIndex = i
		} else if effColRegex.MatchString(colName) {
			e.effectifColIndexes = append(e.effectifColIndexes, i)
		}
	}

	if e.sirenColIndex == -1 {
		return nil, fmt.Errorf("no \"%s\" column found in header: %v", sirenColumn, header)
	}
	if len(e.effectifColIndexes) == 0 {
		return nil, fmt.Errorf("no effectif column (matching ^eff[0-9]+$) found in header: %v", header)
	}

	// compute how many columns to ignore for guessLastNMissing
	lastEffIdx := e.effectifColIndexes[len(e.effectifColIndexes)-1]
	nIgnoredCols := len(header) - 1 - lastEffIdx

	nTrailingEmpty, err := guessLastNMissing(file, nIgnoredCols)
	if err != nil {
		return nil, err
	}
	e.effectifColIndexes = e.effectifColIndexes[:len(e.effectifColIndexes)-nTrailingEmpty]

	return &e, nil
}

func readHeader(file engine.BatchFile) ([]string, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := newCsvReader(f)
	return r.Read()
}

type effectifDataExtractor struct {
	sirenColIndex      int
	effectifColIndexes []int
}

func (e *effectifDataExtractor) GetSiren(record []string) string {
	return record[e.sirenColIndex]
}

func (e *effectifDataExtractor) GetEffectif(record []string) []string {
	effectifData := make([]string, len(e.effectifColIndexes))
	for j, i := range e.effectifColIndexes {
		effectifData[j] = record[i]
	}
	return effectifData
}

// getImportPerimeter makes a perimeter on effectif criterias alone
func getImportPerimeter(effectifEntFile engine.BatchFile, nbMois, minEffectif int, extractor *effectifDataExtractor) (map[string]struct{}, FilterStats, error) {
	f, err := effectifEntFile.Open()
	if err != nil {
		return nil, FilterStats{}, err
	}
	defer f.Close()

	r := newCsvReader(f)

	stats := FilterStats{}
	detectedSirens := map[string]struct{}{} // smaller memory footprint than map[string]bool
	if _, err := r.Read(); err != nil {     // skip header
		return nil, FilterStats{}, err
	}

	skippedLines := 0
	for {
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, FilterStats{}, err
		}

		siren := extractor.GetSiren(record)
		if len(siren) != 9 {
			skippedLines++
			continue
		}

		stats.TotalCompanies++
		effData := extractor.GetEffectif(record)

		// Filter 1: Effectif threshold over nbMois
		if !isInsidePerimeter(effData, nbMois, minEffectif) {
			stats.ExcludedByEffectifFilter++
			continue
		}
		stats.AfterEffectifFilter++

		// Filter 2: Recent data availability
		if !hasRecentEffectifData(effData, DefaultNbMoisRecentData) {
			stats.ExcludedByRecentDataFilter++
			continue
		}
		stats.AfterRecentDataFilter++

		detectedSirens[siren] = struct{}{}
	}
	if skippedLines > 0 {
		slog.Info(fmt.Sprintf("%d lines with bad siren skipped in the effectif_ent file at filter creation", skippedLines))
	}
	return detectedSirens, stats, nil
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

// hasRecentEffectifData checks if the company has at least one non-empty effectif value
// in the last nbMoisRecent months
func hasRecentEffectifData(record []string, nbMoisRecent int) bool {
	for i := len(record) - 1; i >= len(record)-nbMoisRecent && i >= 0; i-- {
		if record[i] != "" {
			return true
		}
	}
	return false
}

// generateFilterReport creates a detail-filtres.md file with detailed filter statistics
func generateFilterReport(stats FilterStats) error {
	exportPath := filepath.Join(sinks.DefaultExportPath, viper.GetString("batch"))
	filename := filepath.Join(exportPath, "detail-filtres.md")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create filter report file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "# Détail des filtres appliqués\n\n")
	fmt.Fprintf(f, "Généré le : %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "## Paramètres\n\n")
	fmt.Fprintf(f, "- **Fenêtre d'observation effectifs** : %d mois\n", DefaultNbMois)
	fmt.Fprintf(f, "- **Seuil minimum d'effectif** : %d employés\n", DefaultMinEffectif)
	fmt.Fprintf(f, "- **Nombre de mois avec données récentes requis** : %d mois\n\n", DefaultNbMoisRecentData)

	// Write filter statistics
	fmt.Fprintf(f, "## Résultats du filtrage\n\n")
	fmt.Fprintf(f, "| Étape | Entreprises | Exclues | Taux d'exclusion |\n")
	fmt.Fprintf(f, "|-------|-------------|---------|------------------|\n")

	// Initial count
	fmt.Fprintf(f, "| **Entreprises initiales** | %d | - | - |\n", stats.TotalCompanies)

	// Filter 1: Effectif threshold
	effectifExclusionRate := 0.0
	if stats.TotalCompanies > 0 {
		effectifExclusionRate = float64(stats.ExcludedByEffectifFilter) / float64(stats.TotalCompanies) * 100
	}
	fmt.Fprintf(f, "| Suppression des entreprises ne dépassant pas %d employés sur les %d derniers mois | %d | %d | %.2f%% |\n",
		DefaultMinEffectif, DefaultNbMois, stats.AfterEffectifFilter, stats.ExcludedByEffectifFilter, effectifExclusionRate)

	// Filter 2: Recent data
	recentDataExclusionRate := 0.0
	if stats.AfterEffectifFilter > 0 {
		recentDataExclusionRate = float64(stats.ExcludedByRecentDataFilter) / float64(stats.AfterEffectifFilter) * 100
	}
	fmt.Fprintf(f, "| Suppression des entreprises n'ayant pas de données d'effectif sur les %d derniers mois | %d | %d | %.2f%% |\n",
		DefaultNbMoisRecentData, stats.AfterRecentDataFilter, stats.ExcludedByRecentDataFilter, recentDataExclusionRate)

	// Note about subsequent filters
	fmt.Fprintf(f, "\n## Filtres appliqués ultérieurement (en SQL)\n\n")
	fmt.Fprintf(f, "Les filtres suivants sont appliqués au niveau de la base de données (entre les couches `stg_*` et `clean_*`) :\n\n")
	fmt.Fprintf(f, "1. **Suppression de certaines formes juridiques**\n")
	fmt.Fprintf(f, "   - Établissements publics nationaux et locaux\n")
	fmt.Fprintf(f, "   - Communes, départements, associations\n")
	fmt.Fprintf(f, "   - Syndicats et groupements de collectivités territoriales\n\n")
	fmt.Fprintf(f, "2. **Suppression de certaines activités principales** (codes NAF)\n")
	fmt.Fprintf(f, "   - Administration publique (84.XX)\n")
	fmt.Fprintf(f, "   - Enseignement (85.XX)\n")
	fmt.Fprintf(f, "   - Organisations associatives (94.XX)\n")
	fmt.Fprintf(f, "   - Services financiers et assurance (64.XX, 65.XX, 66.XX)\n")
	fmt.Fprintf(f, "   - Organisations extraterritoriales (99.XX)\n\n")
	fmt.Fprintf(f, "3. **Suppression des entreprises domiciliées hors de France**\n")
	fmt.Fprintf(f, "   - Exclusion des sièges avec département vide\n\n")

	// Final note
	fmt.Fprintf(f, "## Périmètre final\n\n")
	fmt.Fprintf(f, "Le périmètre final est disponible dans :\n")
	fmt.Fprintf(f, "- **Base de données** : table `sfdata.clean_filter`\n")
	fmt.Fprintf(f, "- **Export CSV** : fichier `clean_filter.csv`\n")

	slog.Info("filter report generated", "filename", filename)
	return nil
}

// SQLFilterStats contains statistics about SQL-based filtering
type SQLFilterStats struct {
	BeforeSQLFilter    int
	AfterSQLFilter     int
	ExcludedBySQL      int
	ExcludedByJuridic  int
	ExcludedByNAF      int
	ExcludedByForeign  int
}

// UpdateFilterReportWithSQLStats appends SQL filter statistics to the existing report
func UpdateFilterReportWithSQLStats(connGetter func() (interface{}, error), batchKey string) error {
	conn, err := connGetter()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Type assertion to get the actual connection type
	type queryExecutor interface {
		QueryRow(ctx context.Context, query string, args ...interface{}) interface{ Scan(...interface{}) error }
	}

	db, ok := conn.(queryExecutor)
	if !ok {
		return fmt.Errorf("connection does not support QueryRow")
	}

	ctx := context.Background()

	// Query SQL filter statistics
	var stats SQLFilterStats
	
	// Count before SQL filters (stg_filter_import)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM stg_filter_import").Scan(&stats.BeforeSQLFilter)
	if err != nil {
		return fmt.Errorf("failed to count stg_filter_import: %w", err)
	}

	// Count excluded by SQL filters (siren_blacklist)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM siren_blacklist").Scan(&stats.ExcludedBySQL)
	if err != nil {
		return fmt.Errorf("failed to count siren_blacklist: %w", err)
	}

	// Count after SQL filters (clean_filter)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM clean_filter").Scan(&stats.AfterSQLFilter)
	if err != nil {
		return fmt.Errorf("failed to count clean_filter: %w", err)
	}

	// Count detailed exclusions by category
	detailQuery := `
		WITH excluded_categories AS (
			SELECT ARRAY[
				'4110', '4120', '4140', '4160', '7210', '7220', '7346', '7348', 
				'7366', '7373', '7379', '7383', '7389', '7410', '7430', '7470', '7490'
			] AS categories
		),
		excluded_activities AS (
			SELECT ARRAY[
				'84.11Z', '84.12Z', '84.13Z', '84.21Z', '84.22Z', '84.23Z', '84.24Z', '84.25Z',
				'84.30A', '84.30B', '84.30C', '85.10Z', '85.20Z', '85.31Z', '85.32Z', '85.41Z',
				'85.42Z', '94.11Z', '94.12Z', '94.20Z', '94.91Z', '94.92Z', '64.11Z', '64.19Z',
				'64.30Z', '64.91Z', '64.92Z', '64.99Z', '65.11Z', '65.12Z', '65.20Z', '65.30Z',
				'66.11Z', '66.12Z', '66.19A', '66.19B', '66.21Z', '66.22Z', '66.29Z', '66.30Z', '99.00Z'
			] AS activities
		)
		SELECT 
			COUNT(*) FILTER (WHERE sirene_ul.categorie_juridique::text = ANY (ec.categories)) as juridic_count,
			COUNT(*) FILTER (WHERE sirene_ul.activite_principale::text = ANY (ea.activities)) as naf_count,
			COUNT(*) FILTER (WHERE sirene.siren IS NULL) as foreign_count
		FROM stg_filter_import fp
		LEFT JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren::text = fp.siren::text
		LEFT JOIN stg_sirene sirene ON sirene.siren = fp.siren 
			AND sirene.siege = true
			AND sirene.departement <> ''
		CROSS JOIN excluded_categories ec
		CROSS JOIN excluded_activities ea
		WHERE fp.siren IN (SELECT siren FROM siren_blacklist)
	`
	
	err = db.QueryRow(ctx, detailQuery).Scan(&stats.ExcludedByJuridic, &stats.ExcludedByNAF, &stats.ExcludedByForeign)
	if err != nil {
		return fmt.Errorf("failed to get detailed exclusion counts: %w", err)
	}

	// Open report file in append mode
	exportPath := filepath.Join(sinks.DefaultExportPath, batchKey)
	filename := filepath.Join(exportPath, "detail-filtres.md")
	
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open filter report for update: %w", err)
	}
	defer f.Close()

	// Append SQL filter statistics
	fmt.Fprintf(f, "\n---\n\n")
	fmt.Fprintf(f, "**Mise à jour après application des filtres SQL** (généré le : %s)\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	sqlExclusionRate := 0.0
	if stats.BeforeSQLFilter > 0 {
		sqlExclusionRate = float64(stats.ExcludedBySQL) / float64(stats.BeforeSQLFilter) * 100
	}
	
	fmt.Fprintf(f, "| Suppression via filtres SQL (formes juridiques, NAF, domiciliation) | %d | %d | %.2f%% |\n",
		stats.AfterSQLFilter, stats.ExcludedBySQL, sqlExclusionRate)
	
	fmt.Fprintf(f, "\n### Détail des exclusions SQL\n\n")
	fmt.Fprintf(f, "- **Formes juridiques exclues** : %d entreprises\n", stats.ExcludedByJuridic)
	fmt.Fprintf(f, "- **Codes NAF exclus** : %d entreprises\n", stats.ExcludedByNAF)
	fmt.Fprintf(f, "- **Domiciliation hors France** : %d entreprises\n", stats.ExcludedByForeign)
	fmt.Fprintf(f, "\n*Note : Une entreprise peut être comptée dans plusieurs catégories si elle remplit plusieurs critères d'exclusion.*\n")
	
	fmt.Fprintf(f, "\n### Périmètre final\n\n")
	fmt.Fprintf(f, "**%d entreprises** dans le périmètre final après tous les filtres.\n\n", stats.AfterSQLFilter)
	
	totalExcluded := stats.BeforeSQLFilter - stats.AfterSQLFilter
	finalExclusionRate := 0.0
	if stats.BeforeSQLFilter > 0 {
		finalExclusionRate = float64(totalExcluded) / float64(stats.BeforeSQLFilter) * 100
	}
	fmt.Fprintf(f, "**Taux d'exclusion global** : %.2f%% (%d entreprises exclues sur %d initiales)\n",
		finalExclusionRate, totalExcluded, stats.BeforeSQLFilter)

	slog.Info("filter report updated with SQL statistics", "filename", filename)
	return nil
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
