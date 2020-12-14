package marshal

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sfregexp"

	"github.com/spf13/viper"
)

// SirenFilter liste les numéros SIREN d'entreprise et établissements à exclure des traitements.
type SirenFilter map[string]bool

// Skips retourne `false` si le numéro SIREN/SIRET peut être traité,
// car il est inclus dans le Filtre, ou car il n'y a pas de filtre.
func (filter SirenFilter) Skips(siretOrSiren string) bool {
	return filter != nil && !filter.includes(siretOrSiren)
}

// includes retourne `true` si le numéro SIREN/SIRET est inclus, c.a.d. à traiter.
func (filter SirenFilter) includes(siretOrSiren string) bool {
	if len(siretOrSiren) >= 9 {
		return filter[siretOrSiren[0:9]]
	}
	return false
}

// GetSirenFilterFromCache reads the filter from cache.
func GetSirenFilterFromCache(cache Cache) SirenFilter {
	value, err := cache.Get("filter")
	if err == nil {
		filter, ok := value.(SirenFilter)
		if ok {
			return filter
		}
	}
	return nil
}

// GetSirenFilter reads the filter from cache if it cans, or else it reads it
// from input files and stores it in cache
func GetSirenFilter(cache Cache, batch *base.AdminBatch) (SirenFilter, error) {
	return getSirenFilter(cache, batch, readFilterFiles)
}

// getSirenFilter reads the filter from cache if it cans, or else it reads it
// from input files and stores it in cache
func getSirenFilter(cache Cache, batch *base.AdminBatch, fr filterReader) (SirenFilter, error) {

	value, err := cache.Get("filter")

	if err == nil {
		filter, ok := value.(SirenFilter)
		if ok {
			return filter, nil
		}
		return nil, errors.New("Wrong format from existing field filter in cache")
	}

	paths := batch.Files["filter"]
	if len(paths) == 0 {
		// No filter
		return nil, nil
	}

	basePath := viper.GetString("APP_DATA")
	filter, err := fr(basePath, paths)
	if err != nil {
		return nil, err
	}
	cache.Set("filter", filter)
	return filter, nil
}

type filterReader func(string, []string) (SirenFilter, error)

// openAndReadFilters reads several files, reads their content and concatenate
// it into a SirenFilter
func readFilterFiles(basePath string, filenames []string) (SirenFilter, error) {
	filter := make(SirenFilter)
	for _, p := range filenames {
		file, err := os.Open(filepath.Join(basePath, p))
		if err != nil {
			return nil, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
		}
		defer file.Close()
		err = readFilter(bufio.NewReader(file), filter)
		if err != nil {
			return nil, errors.New("Erreur à la lecture du fichier, " + err.Error())
		}
	}
	return filter, nil
}

// readFilter reads the content of a io.Reader and adds it to an existing
// filter
func readFilter(reader io.Reader, filter SirenFilter) error {

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

		if sfregexp.RegexpDict["siren"].MatchString(siren) {
			filter[siren] = true
		} else {
			return errors.New("Format de siren incorrect trouvé : " + siren)
		}
	}
	return nil
}
