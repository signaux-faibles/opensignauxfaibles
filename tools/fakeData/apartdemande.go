package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func readAndRandomApartDemande(fileName string, outputFileName string, mapping map[string]string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}
	// destination
	outputFile := xlsx.NewFile()
	outputFile.AddSheet("DEMANDE")
	outputFile.Sheets[0].Cols = file.Sheets[0].Cols
	outputFile.Sheet["DEMANDE"].Rows = append(outputFile.Sheet["DEMANDE"].Rows, file.Sheets[0].Rows[0])
	for _, row := range file.Sheets[0].Rows[1:] {
		l := len(row.Cells)
		for i := 0; i < l; i++ {
			if !contains([]int{2, 3, 14, 15, 16, 20, 21, 22, 24, 26, 30, 31}, i) {
				row.Cells[i].Value = ""
			}
		}

		if l > 3 {
			siret := strings.Replace(row.Cells[3].Value, " ", "", -1)
			row.Cells[3].Value = mapping[siret]
		}

		for _, i := range []int{14, 15, 22, 24, 26, 30, 31} {
			if l > 31 && row.Cells[i].Value != "" {
				v, err := strconv.ParseFloat(row.Cells[i].Value, 64)
				if err != nil {
					panic(err)
				}
				row.Cells[i].Value = strconv.Itoa(int(v * rand.Float64() * 2))
			}
		}

		if l > 3 && row.Cells[3].Value != "" {
			outputFile.Sheet["DEMANDE"].Rows = append(outputFile.Sheet["DEMANDE"].Rows, row)
		}

	}

	outputFile.Save(outputFileName)

	return nil
}
