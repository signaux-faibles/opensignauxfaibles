package filter

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/sfregexp"
	"os"
)

// FileReader reads the filter from a CSV file.
// Implements filterReader
type FileReader struct {
	BatchFile base.BatchFile
}

func (f *FileReader) Read() (engine.SirenFilter, error) {
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
