package marshal

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/spf13/viper"
)

// IsFiltered determines if the siret must be filtered or not
func IsFiltered(id string, filter map[string]bool) (bool, error) {

	validSiret := sfregexp.RegexpDict["siret"].MatchString(id)
	validSiren := sfregexp.RegexpDict["siren"].MatchString(id)
	if !validSiret && !validSiren {
		return true, errors.New("Le siret/siren est invalide") // TODO: retirer la validation de cette fonction
	}

	// if no filter, then all ids pass
	if filter == nil {
		return false, nil
	}
	return !filter[id[0:9]], nil
}

// GetSirenFilter reads the filter from cache if it cans, or else it reads it
// from input files and stores it in cache
func GetSirenFilter(cache base.Cache, batch *base.AdminBatch) (map[string]bool, error) {
	return getSirenFilter(cache, batch, readFilterFiles)
}

// getSirenFilter reads the filter from cache if it cans, or else it reads it
// from input files and stores it in cache
func getSirenFilter(cache base.Cache, batch *base.AdminBatch, fr filterReader) (map[string]bool, error) {

	value, err := cache.Get("filter")

	if err == nil {
		filter, ok := value.(map[string]bool)
		if ok {
			return filter, nil
		} else {
			return nil, errors.New("Wrong format from existing field filter in cache")
		}
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

type filterReader func(string, []string) (map[string]bool, error)

// openAndReadFilters reads several files, reads their content and concatenate
// it into a map[string]bool
func readFilterFiles(basePath string, filenames []string) (map[string]bool, error) {
	filter := make(map[string]bool)
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
func readFilter(reader io.Reader, filter map[string]bool) error {

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
