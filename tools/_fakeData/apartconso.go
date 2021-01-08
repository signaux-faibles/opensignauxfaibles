package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func readAndRandomApartConso(fileName string, outputFileName string, mapping map[string]string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}
	// destination
	outputFile := xlsx.NewFile()
	outputFile.AddSheet("INDEMNITE")
	outputFile.Sheets[0].Cols = file.Sheets[0].Cols
	outputFile.Sheet["INDEMNITE"].Rows = append(outputFile.Sheet["INDEMNITE"].Rows, file.Sheets[0].Rows[0])
	for _, row := range file.Sheets[0].Rows[1:] {
		l := len(row.Cells)
		for i := 0; i < l; i++ {
			if !contains([]int{1, 2, 15, 16, 17, 18}, i) {
				row.Cells[i].Value = ""
			}
		}

		if l > 2 {
			siret := strings.Replace(row.Cells[2].Value, " ", "", -1)
			row.Cells[2].Value = mapping[siret]
		}

		for _, i := range []int{16, 17, 18} {
			if l > 31 && row.Cells[i].Value != "" {
				v, err := strconv.ParseFloat(row.Cells[i].Value, 64)
				if err != nil {
					panic(err)
				}
				row.Cells[i].Value = strconv.Itoa(int(v * rand.Float64() * 2))
			}
		}

		if l > 3 && row.Cells[2].Value != "" {
			outputFile.Sheet["INDEMNITE"].Rows = append(outputFile.Sheet["INDEMNITE"].Rows, row)
		}
	}

	outputFile.Save(outputFileName)

	return nil
}
