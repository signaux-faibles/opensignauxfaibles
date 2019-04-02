package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TomOnTime/utfutil"
	"golang.org/x/text/encoding/unicode"
)

func readAndRandomDiane(fileName string, outputFileName string, mapping map[string]string) error {
	rand.Seed(time.Now().UTC().UnixNano())
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder()
	knownSiren := make(map[string]struct{})

	// source
	file, err := utfutil.OpenFile(fileName, utfutil.UTF16LE)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	sirens := make(map[string]string)
	for k, v := range mapping {
		sirens[k[0:9]] = v[0:9]
	}

	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// ligne de titre
	row, err := reader.Read()
	outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
	titleRow, err := encoder.String(outputRow)
	if err != nil {
		panic(err)
	}
	_, err = outputFile.WriteString(titleRow)
	if err != nil {
		return err
	}

	newRow := make([]string, 452)

	ints := []int{
		5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27,
		34, 35, 36, 50, 51, 52, 53, 54, 55, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113,
		114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 140, 141, 142,
		143, 144, 145, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 188, 189,
		190, 191, 192, 193,
	}
	ints = append(ints, iter(290, 451)...)

	immutable := []int{0, 3, 4, 37, 38, 39, 40, 41, 42, 43, 44, 45, 16, 47, 48, 49}
	floats := []int{28, 29, 30, 31, 32, 33, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67,
		68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
		89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 128, 129, 130, 131,
		132, 133, 134, 135, 136, 137, 138, 139, 146, 147, 148, 149, 150, 151, 152, 153, 154,
		155, 156, 157, 158, 159, 160, 161, 162, 163, 176, 177, 178, 179, 180, 181, 182, 183,
		184, 185, 186, 187,
	}
	floats = append(floats, iter(194, 289)...)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if _, ok := knownSiren[row[2]]; !ok && row[1] != "Nom de l'entreprise" {
			knownSiren[row[2]] = struct{}{}

			for _, i := range ints {
				newRow[i] = randomInt(row[i])
			}

			for _, i := range immutable {
				newRow[i] = row[i]
			}

			for _, i := range floats {
				newRow[i] = randomFloat(row[i])
			}

			newRow[1] = ""
			newRow[2] = sirens[row[2]]

			if newRow[2] != "" {
				outRow := "\"" + strings.Join(newRow, "\";\"") + "\"\n"
				encodedRow, err := encoder.String(outRow)
				_, err = outputFile.WriteString(encodedRow)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func randomInt(i string) string {
	if i == "" {
		return ""
	}
	v, err := strconv.ParseFloat(i, 64)
	if err != nil {
		panic(err)
	}
	v = v * (rand.Float64() * 2)
	return strconv.Itoa(int(v))
}

func randomFloat(i string) string {
	if i == "" {
		return ""
	}
	j := strings.Replace(i, ",", ".", -1)

	v, err := strconv.ParseFloat(j, 64)
	if err != nil {
		panic(err)
	}
	v = v * (rand.Float64() * 2)

	return strings.Replace(fmt.Sprintf("%-6.2f", v), ".", ",", -1)
}

func iter(a, b int) []int {
	var s []int
	for i := a; i <= b; i++ {
		s = append(s, i)
	}
	return s
}
