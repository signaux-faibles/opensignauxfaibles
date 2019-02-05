package main

import (
  "bufio"
  "encoding/csv"
  "fmt"
  "io"
  "os"
  "strconv"
  "time"

  "regexp"
  "github.com/cnf/structhash"
  "github.com/spf13/viper"
)

// DPAE Déclaration préalabre à l'embauche
type DPAE struct {
  Siret    string    `json:"-" bson:"-"`
  Date     time.Time `json:"date" bson:"date"`
  CDI      float64   `json:"cdi" bson:"cdi"`
  CDDLong  float64   `json:"cdd_long" bson:"cdd_long"`
  CDDCourt float64   `json:"cdd_court" bson:"cdd_court"`
}

func parseDPAE(path string) chan *DPAE {
  outputChannel := make(chan *DPAE)

  file, err := os.Open(viper.GetString("APP_DATA") + path)
  if err != nil {
    journal(critical, "importDPAE", "Erreur à l'ouverture du fichier "+path+": "+err.Error())
  }


  reader := csv.NewReader(bufio.NewReader(file))
  reader.Comma = ';'
  reader.Read()

  go func() {
    e:= 0
    n := 0
    for {
      row, error := reader.Read()
      journal(info, "importDPAE", "Ouverture du fichier "+path)
      if error == io.EOF {
        break
      } else if error != nil {
        journal(critical, "importInterim", "Erreur à la lecture du fichier "+path+": "+err.Error())
      }

      date, err := time.Parse("20060102", row[1]+row[2]+"01")

      siret := row[0]
      validSiret, _ := regexp.MatchString("[0-9]{14}", siret)
      n++
      if (validSiret && err == nil){
        dpae := DPAE{
          Siret: row[0],
          Date:  date,
        }
        dpae.CDI, _ = strconv.ParseFloat(row[3], 64)
        dpae.CDDLong, _ = strconv.ParseFloat(row[4], 64)
        dpae.CDDCourt, _ = strconv.ParseFloat(row[5], 64)

        outputChannel <- &dpae

      } else {
        e++
      }
    }
    file.Close()
    close(outputChannel)
  }()
  return outputChannel
}

func importDPAE(batch *AdminBatch) error {
  journal(info, "importDPAE", "Import du batch "+batch.ID.Key+": DPAE")
  for _, dpaeFile := range batch.Files["dpae"] {
    for dpae := range parseDPAE(dpaeFile) {
      hash := fmt.Sprintf("%x", structhash.Md5(dpae, 1))

      value := Value{
        Value: Data{
          Scope: "etablissement",
          Key:   dpae.Siret,
          Batch: map[string]Batch{
            batch.ID.Key: Batch{
              DPAE: map[string]*DPAE{
                hash: dpae,
              },
            },
          },
        },
      }
      db.ChanData <- &value
    }
  }
  db.ChanData <- &Value{}
  journal(info, "importDPAE", "Fin de l'import du batch "+batch.ID.Key+": DPAE")
  return nil

}
