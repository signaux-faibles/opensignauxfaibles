package marshal

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	//"strconv"
	"sort"
	"time"

	"github.com/spf13/viper"
)

// GetSiret gets the siret related to a specific compte at a given point in
// time
func GetSiret(compte string, date *time.Time, cache base.Cache, batch *base.AdminBatch) (string, error) {
	comptes, err := GetCompteSiretMapping(cache, batch, OpenAndReadSiretMapping)

	if err != nil {
		return "", err
	}

	for _, sd := range comptes[compte] {
		if date.Before(sd.Date) {
			return sd.Siret, nil
		}
	}
	return "", errors.New("Pas de siret associé au compte " + compte + " à la période " + date.String())
}

// SiretDate holds a pair of a siret and a date
type SiretDate struct {
	Siret string
	Date  time.Time
}

// Comptes associates a SiretDate to an urssaf account number
type Comptes map[string][]SiretDate

// GetCompteSiretMapping returns the siret mapping in cache if available, else
// reads the file and save it in cache
func GetCompteSiretMapping(cache base.Cache, batch *base.AdminBatch, mr mappingReader) (Comptes, error) {

	value, err := cache.Get("comptes")
	if err == nil {
		comptes, ok := value.(Comptes)
		if ok {
			return comptes, nil
		} else {
			return nil, errors.New("Wrong format from existing field comptes in cache")
		}
	}

	fmt.Println("Chargement des comptes urssaf")

	compteSiretMapping := make(Comptes)

	path := batch.Files["admin_urssaf"]
	basePath := viper.GetString("APP_DATA")

	if len(path) == 0 {
		return nil, errors.New("No admin_urssaf mapping found")
	}
	for _, p := range path {
		compteSiretMapping, err = mr(basePath, p, compteSiretMapping, cache, batch)
		if err != nil {
			return nil, err
		}
	}
	cache.Set("comptes", compteSiretMapping)
	return compteSiretMapping, nil
}

type mappingReader func(string, string, Comptes, base.Cache, *base.AdminBatch) (Comptes, error)

// OpenAndReadSiretMapping opens files and reads their content
func OpenAndReadSiretMapping(
	basePath string,
	endPath string,
	compteSiretMapping Comptes,
	cache base.Cache,
	batch *base.AdminBatch,
) (Comptes, error) {

	file, err := os.Open(basePath + endPath)
	if err != nil {
		return nil, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
	}
	defer file.Close()

	addSiretMapping, err := readSiretMapping(bufio.NewReader(file), cache, batch)
	if err != nil {
		return nil, err
	}
	for key := range addSiretMapping {
		compteSiretMapping[key] = addSiretMapping[key]
	}
	return compteSiretMapping, nil
}

//readSiretMapping reads a admin_urssaf file
func readSiretMapping(
	reader io.Reader,
	cache base.Cache,
	batch *base.AdminBatch,
) (Comptes, error) {

	var addSiretMapping = make(map[string][]SiretDate)

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	// discard header row
	csvReader.Read()

	compteIndex := 2
	siretIndex := 5
	fermetureIndex := 7

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		maxTime := "9990101"

		if row[fermetureIndex] == "" {
			row[fermetureIndex] = maxTime
		} // compte non fermé

		// fermeture, err := urssafToDate(row[fermetureIndex])
		fermeture, err := urssafToDate(row[fermetureIndex])
		if err != nil {
			return nil, err // fermeture n'a pas pu être lue ou convertie en date
		}

		compte := row[compteIndex]
		siret := row[siretIndex]

		if !sfregexp.RegexpDict["siret"].MatchString(siret) {
			continue
		}

		filter, err := GetSirenFilter(cache, batch)
		if err != nil {
			return nil, err
		}

		filtered, err := IsFiltered(siret, filter)
		if err != nil {
			return nil, err
		}

		if sfregexp.RegexpDict["siret"].MatchString(siret) && !filtered {
			//siret valide
			addSiretMapping[compte] = append(addSiretMapping[compte], SiretDate{siret, fermeture})
			// Tri des sirets pour chaque compte par ordre croissant de date de fermeture
			// TODO pour être exact, trier également selon que le compte est ouvert ou fermé. Comptes ouverts d'abord dans la liste.
			// Permettrait d'éviter de sélectionner des comptes fermés mais dont la date de fermeture n'a pas encore été renseignée
			sort.Slice(
				addSiretMapping[compte],
				func(i, j int) bool {
					return (addSiretMapping[compte][i].Date.Before(addSiretMapping[compte][j].Date))
				},
			)
		}
	}
	return addSiretMapping, nil
}
