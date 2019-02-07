package urssaf

import (
	"bufio"
	"dbmongo/lib/engine"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func getCompteSiretMapping(batch *engine.AdminBatch) (map[string]string, error) {
	compteSiretMapping := make(map[string]string)
	compteSiretLast := make(map[string]int)

	path := batch.Files["admin_urssaf"]
	basePath := viper.GetString("APP_DATA")

	for _, p := range path {
		file, err := os.Open(basePath + p)
		if err != nil {
			return map[string]string{}, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
		}

		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'

		// discard header row
		reader.Read()

		siretIndex := 3
		compteIndex := 0
		fermetureIndex := 5

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return map[string]string{}, err
			}

			_, err1 := strconv.Atoi(row[siretIndex])
			fermeture, err2 := strconv.Atoi(row[fermetureIndex])
			if err2 != nil {
				if row[fermetureIndex] == "" {
					fermeture = 1
				} else {
					return map[string]string{}, err2
				}
			}
			derniereFermetureLue, ok := compteSiretLast[row[compteIndex]]
			if err1 == nil &&
				len(row[siretIndex]) == 14 &&
				(!ok ||
					(derniereFermetureLue != 0 && derniereFermetureLue < fermeture) ||
					fermeture == 0) {

				compteSiretMapping[row[compteIndex]] = row[siretIndex]
				compteSiretLast[row[compteIndex]] = fermeture
			}
		}
		file.Close()
	}
	return compteSiretMapping, nil
}
