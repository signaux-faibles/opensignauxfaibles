package engine

import (
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"sort"
	"time"

	"github.com/spf13/viper"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/sfregexp"
)

// SiretDate holds a pair of a siret and a date
type SiretDate struct {
	Siret string
	Date  time.Time
}

// Comptes associates a SiretDate to an urssaf account number
type Comptes map[string][]SiretDate

// GetSortedKeys retourne la liste classée des numéros de Comptes.
func (comptes *Comptes) GetSortedKeys() []string {
	keys := make([]string, len(*comptes))
	i := 0
	for k := range *comptes {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// GetSiretFromComptesMapping gets the siret related to a specific compte at a
// given point in time
func (comptes *Comptes) GetSiret(compte string, date *time.Time) (string, error) {
	for _, sd := range (*comptes)[compte] {
		if date.Before(sd.Date) {
			return sd.Siret, nil
		}
	}
	return "", errors.New("Pas de siret associé au compte " + compte + " à la période " + date.String())
}

// GetCompteSiretMapping returns the siret mapping in cache if available, else
// reads the file and save it in cache. Lazy loaded.
func GetCompteSiretMapping(cache Cache, batch *base.AdminBatch, filter SirenFilter, mr mappingReader) (Comptes, error) {
	value, err := cache.Get("comptes")
	slog.Debug("associe les siret et les numéros de compte URSSAF", slog.Any("AdminBatch", *batch))
	if err == nil {
		comptes, ok := value.(Comptes)
		if ok {
			return comptes, nil
		}
		return nil, errors.New("wrong format from existing field comptes in cache")
	}

	compteSiretMapping := make(Comptes)

	path := batch.Files["admin_urssaf"]
	basePath := viper.GetString("APP_DATA")

	if len(path) == 0 {
		return nil, errors.New("no admin_urssaf mapping found")
	}
	for _, p := range path {
		compteSiretMapping, err = mr(basePath, p, compteSiretMapping, cache, batch, filter)
		if err != nil {
			slog.Error("erreur pendant le mapping siret <-> compte", slog.Any("error", err))
			return nil, err
		}
		slog.Debug("mapping siret <-> compte", slog.Any("comptes", compteSiretMapping.GetSortedKeys()))
	}
	cache.Set("comptes", compteSiretMapping)
	slog.Debug("Chargement des comptes URSSAF terminé", slog.Any("comptesValue", value))
	return compteSiretMapping, nil
}

type mappingReader func(string, base.BatchFile, Comptes, Cache, *base.AdminBatch, SirenFilter) (Comptes, error)

// OpenAndReadSiretMapping opens files and reads their content
func OpenAndReadSiretMapping(
	basePath string,
	batchFile base.BatchFile,
	compteSiretMapping Comptes,
	cache Cache,
	batch *base.AdminBatch,
	filter SirenFilter,
) (Comptes, error) {

	file, fileReader, err := OpenFileReader(batchFile)
	if err != nil {
		return nil, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
	}
	defer file.Close()

	addSiretMapping, err := readSiretMapping(fileReader, cache, batch, filter)
	slog.Debug("lecture du mapping des sirets", slog.Any("mapping", addSiretMapping))
	if err != nil {
		return nil, err
	}
	for key := range addSiretMapping {
		compteSiretMapping[key] = addSiretMapping[key]
	}
	slog.Debug("nouveau mapping des sirest", slog.Any("mapping", compteSiretMapping))
	return compteSiretMapping, nil
}

// readSiretMapping reads a admin_urssaf file
func readSiretMapping(
	reader io.Reader,
	cache Cache,
	batch *base.AdminBatch,
	filter SirenFilter,
) (Comptes, error) {

	var addSiretMapping = make(map[string][]SiretDate)

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	// parse header row
	fields, err := csvReader.Read() // => Urssaf_gestion;Dep;Compte;Etat_compte;Siren;Siret;Date_crea_siret;Date_disp_siret;Cle_md5
	if err != nil {
		return nil, err
	}
	idx := indexFields(LowercaseFields(fields))
	requiredFields := []string{"compte", "siret", "date_disp_siret"}
	if _, err := idx.HasFields(requiredFields); err != nil {
		return nil, err
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		idxRow := idx.IndexRow(row)

		maxTime := "9990101"

		fermetureRaw := idxRow.GetVal("date_disp_siret")

		if fermetureRaw == "" {
			fermetureRaw = maxTime
		} // compte non fermé

		// fermeture, err := UrssafToDate(fermetureRaw)
		fermeture, err := UrssafToDate(fermetureRaw)
		if err != nil {
			return nil, err // fermeture n'a pas pu être lue ou convertie en date
		}

		compte := idxRow.GetVal("compte")
		siret := idxRow.GetVal("siret")

		if sfregexp.ValidSiret(siret) && !filter.ShouldSkip(siret) {
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
	slog.Debug("lecture du mapping des siret terminée", slog.Any("nombre de mapping", len(addSiretMapping)))
	return addSiretMapping, nil
}
