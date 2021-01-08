package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

func readAndRandomComptes(fileName string, outputFileName string) (map[string]string, error) {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	mapping := make(map[string]string)
	sirens := make(map[string]string)

	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	// ligne de titre
	row, err := reader.Read()
	outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		siren := row[3][0:9]

		var newSiren string
		if _, ok := sirens[siren]; ok {
			newSiren = sirens[siren]
		} else {
			for {
				newSiren = randStringBytesRmndr(9)
				if _, ok := sirens[newSiren]; !ok && newSiren != siren {
					break
				}
			}
		}

		siret := row[3]

		compte := row[0]
		var newSiret, newCompte string

		for {
			newSiret = newSiren + randStringBytesRmndr(5)
			if _, ok := mapping[newSiret]; !ok && newSiret != siret {
				break
			}
		}
		for {
			newCompte = randStringBytesRmndr(len(compte))
			if _, ok := mapping[newCompte]; !ok && newCompte != compte {
				break
			}
		}
		mapping[compte] = newCompte
		mapping[siret] = newSiret
		sirens[siren] = newSiren

		row[0] = newCompte
		row[2] = newSiren
		row[3] = newSiret
		// row[4] = date
		// row[5] = date

		outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
		_, err = outputFile.WriteString(outputRow)
		if err != nil {
			return nil, err
		}
	}
	return mapping, nil
}
