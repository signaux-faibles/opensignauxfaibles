package engine

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/sfregexp"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	cacheKeyFilter = "filter"
	sirenLength    = 9
)

// SirenFilter décrit le périmètre d'import.
// Les clés sont des numéros SIREN.
// Seules les clés présentes et dont la valeur est `true` sont dans le périmètre.
type SirenFilter map[string]bool

// ShouldSkip retourne `true` si le numéro SIREN/SIRET est hors périmètre.
// Si aucun filtre n'est défini, renvoi `false` par défaut.
func (f SirenFilter) ShouldSkip(siretOrSiren string) bool {
	if f == nil {
		return false
	}

	siren := siretOrSiren

	if len(siretOrSiren) >= sirenLength {
		siren = siretOrSiren[:sirenLength]
	}

	return !f[siren]
}

func (f SirenFilter) Add(siren string) error {
	if !sfregexp.RegexpDict["siren"].MatchString(siren) {
		return fmt.Errorf("format SIREN invalide: %s", siren)
	}
	f[siren] = true
	return nil
}

// FilterReader defines the interface for reading SIREN filters from various sources.
type FilterReader interface {
	Read() (SirenFilter, error)
	SuccessStr() string
}

// GetSirenFilter retrieves the SIREN filter using a priority-based approach:
// 1. Cache (fastest)
// 2. Batch filter file (if available)
// 3. Database (fallback)
//
// This is a convenience wrapper that uses default dependencies.
func GetSirenFilter(cache Cache, batch *base.AdminBatch) (SirenFilter, error) {
	filterFile, _ := batch.Files.GetFilterFile()

	readers := []FilterReader{
		&CacheReader{cache},
		&FileReader{filterFile},
		&DBReader{Db.PostgresDB},
	}
	return getSirenFilterFromReaders(cache, readers)
}

// getSirenFilterFromReaders tries each reader in order until one succeeds.
// The first successful filter is cached and returned.
func getSirenFilterFromReaders(cache Cache, readers []FilterReader) (SirenFilter, error) {
	var filter SirenFilter
	var lastErr error

	// Cache the filter when successfully retrieved from any source
	defer func() {
		if filter != nil {
			cache.Set(cacheKeyFilter, filter)
		}
	}()

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

// ----------------------------------------------------------------

// CacheReader reads the filter from cache.
type CacheReader struct {
	Cache Cache
}

// CacheReader reads the filter from cache.
func (f *CacheReader) Read() (SirenFilter, error) {
	value, err := f.Cache.Get("filter")

	if err != nil {
		// value not found, do not return an error
		return nil, nil
	}

	filter, ok := value.(SirenFilter)

	if !ok {
		// This should not happen. Return an error
		return nil, errors.New("error retrieving \"filter\" from cache: value exists but in wrong format")
	}

	return filter, nil
}

func (f *CacheReader) SuccessStr() string {
	return "Filter retrieved from cache"
}

// ----------------------------------------------------------------

// FileReader reads the filter from a CSV file.
// Implements filterReader
type FileReader struct {
	BatchFile base.BatchFile
}

func (f *FileReader) Read() (SirenFilter, error) {
	if f.BatchFile == nil {
		return nil, nil
	}

	p := f.BatchFile.Path()

	file, err := os.Open(p)
	if err != nil {
		return nil, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
	}
	defer file.Close()

	filter := make(SirenFilter)
	err = parseCSVFilter(bufio.NewReader(file), filter)
	if err != nil {
		return nil, errors.New("Erreur à la lecture du fichier, " + err.Error())
	}
	return filter, nil
}

func (f *FileReader) SuccessStr() string {
	return "Filter retrieved from file"
}

// parseCSVFilter reads the content of a io.Reader and adds it to an existing
// filter
func parseCSVFilter(reader io.Reader, filter SirenFilter) error {

	csvreader := csv.NewReader(reader)
	csvreader.Comma = ';'

	sirenIndex := 0

	for {
		row, err := csvreader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		siren := row[sirenIndex]

		if siren == "siren" {
			continue
		}

		if sfregexp.RegexpDict["siren"].MatchString(siren) {
			filter[siren] = true
		} else {
			return errors.New("Format de siren incorrect trouvé : " + siren)
		}
	}
	return nil
}

// ----------------------------------------------------------------

// DBReader reads the filter from the database "filter" table.
type DBReader struct {
	Conn *pgxpool.Pool
}

func (f *DBReader) Read() (SirenFilter, error) {
	var filter = make(SirenFilter)

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
