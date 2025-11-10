package filter

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/sfregexp"
)

// CsvReader reads the filter from a CSV file.
// Implements filterReader
type CsvReader struct {
	BatchFile engine.BatchFile
}

func (f *CsvReader) Read() (engine.SirenFilter, error) {
	if f.BatchFile == nil {
		return nil, nil
	}

	file, err := f.BatchFile.Open()
	if err != nil {
		return nil, errors.New("error opening file: " + err.Error())
	}
	defer file.Close()

	filter := make(MapFilter)
	err = parseCSVFilter(bufio.NewReader(file), filter)
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}

	slog.Debug("Filter retrieved from csv")
	return filter, nil
}

// parseCSVFilter reads the content of a io.Reader and adds it to an existing
// filter
func parseCSVFilter(reader io.Reader, filter MapFilter) error {

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
			return errors.New("invalid SIREN format found: " + siren)
		}
	}
	return nil
}
