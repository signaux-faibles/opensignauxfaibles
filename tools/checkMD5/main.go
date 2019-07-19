package main

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: md5check etablissement.csv")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file: " + err.Error())
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	line := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			line++
			var testRow string
			if row[7] == "" {
				testRow = strings.Join(row[0:7], ";")
			} else {
				testRow = strings.Join(row[0:8], ";")
			}
			h := md5.New()
			io.WriteString(h, testRow)
			sum := fmt.Sprintf("%x", h.Sum(nil))
			if strings.ToUpper(sum) != row[8] {
				fmt.Println("erreur à la ligne " + strconv.Itoa(line) + ": " + strings.ToUpper(sum) + " calculé contre " + row[8])
			}
		}
	}
}
