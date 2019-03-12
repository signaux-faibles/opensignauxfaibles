package urssaf

import (
  "bufio"
  "dbmongo/lib/engine"
  "encoding/csv"
  "errors"
  "io"
  "os"
  //"strconv"
  "time"
  "sort"
  "github.com/spf13/viper"
)

type SiretDate struct{
  Siret string
  Date time.Time
}
type Comptes map[string][]SiretDate


func (c *Comptes) GetSiret(compte string, date time.Time) (string, error) {

  for _, sd := range((*c)[compte]) {
    if date.Before(sd.Date){
      return sd.Siret , nil
    }

  }
    return "", errors.New("Pas de siret associé au compte " + compte + " à cette période")
}

func getCompteSiretMapping(batch *engine.AdminBatch) (Comptes, error) {

  compteSiretMapping := make(map[string][]SiretDate)

  path := batch.Files["admin_urssaf"]
  basePath := viper.GetString("APP_DATA")

  for _, p := range path {
    file, err := os.Open(basePath + p)
    if err != nil {
      return map[string][]SiretDate{}, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
    }

    reader := csv.NewReader(bufio.NewReader(file))
    reader.Comma = ';'

    // discard header row
    reader.Read()

    compteIndex := 0
    siretIndex := 3
    fermetureIndex := 5

    for {
      row, err := reader.Read()
      if err == io.EOF {
        break
      } else if err != nil {
        return map[string][]SiretDate{}, err
      }

      maxTime := "9990101"

      if row[fermetureIndex] == "" {row[fermetureIndex] = "0"} // date de fermeture manquante
      if (row[fermetureIndex] == "0") { row[fermetureIndex] = maxTime } // compte non fermé

      fermeture, err := urssafToDate(row[fermetureIndex])
      if  err != nil {
        return map[string][]SiretDate{}, err // fermeture n'a pas pu être lue ou convertie en date
      }

      compte := row[compteIndex]
      siret := row[siretIndex]
      if len(siret) == 14 {
        //siret valide
        compteSiretMapping[compte] = append(compteSiretMapping[compte], SiretDate{siret, fermeture})
        // Tri des sirets pour chaque compte par ordre croissant de date de fermeture
        // TODO pour être exact, trier également selon que le compte est ouvert ou fermé. Comptes ouverts d'abord dans la liste.
        // Permettrait d'éviter de sélectionner des comptes fermés mais dont la date de fermeture n'a pas encore été renseignée
        sort.Slice(compteSiretMapping[compte],
        func(i, j int) bool {return(
          compteSiretMapping[compte][i].Date.Before(compteSiretMapping[compte][j].Date))})
        }
      }
      file.Close()

    }
    return compteSiretMapping, nil
  }
