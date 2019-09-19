package engine

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"

	//"strconv"
	"github.com/spf13/viper"
)

func getSirenFilter(batch *AdminBatch) (map[string]bool, error) {

	filter := make(map[string]bool)

	path := batch.Files["filter"]
	basePath := viper.GetString("APP_DATA")

	for _, p := range path {
		file, err := os.Open(basePath + p)
		if err != nil {
			file.Close()
			return nil, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
		}

		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'

		sirenIndex := 0

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				file.Close()
				return filter, err
			}

			siren := row[sirenIndex]
			if len(siren) == 9 {
				//siret valide
				filter[siren] = true
			} else {
				file.Close()
				return nil, errors.New("Format de siren incorrect trouvé")
			}
		}
		file.Close()
	}
	return filter, nil
}
