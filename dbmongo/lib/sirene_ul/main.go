package sirene_ul

import (
  //"bufio"
  "dbmongo/lib/engine"
  "encoding/csv"
  "io"
  "os"

  "github.com/chrnin/gournal"
  "github.com/spf13/viper"
)

// Sirene informations sur les entreprises
type SireneUL struct {
  Siren              string     `json:"siren,omitempty" bson:"siren,omitempty"`
  Nic                string     `json:"nic,omitempty" bson:"nic,omitempty"`
  RaisonSociale      string     `json:"raison_sociale" bson:"raison_sociale"`
  CodeStatutJuridique    string     `json:"statut_juridique" bson:"statut_juridique"`
}

// Key id de l'objet
func (sirene_ul SireneUL) Key() string {
  return sirene_ul.Siren
}

// Type de données
func (sirene_ul SireneUL) Type() string {
  return "sirene_ul"
}

// Scope de l'objet
func (sirene_ul SireneUL) Scope() string {
  return "entreprise"
}

// Parser produit les données sirene à partir du fichier geosirene
func Parser(batch engine.AdminBatch, filter map[string]bool) (chan engine.Tuple, chan engine.Event) {
  outputChannel := make(chan engine.Tuple)
  eventChannel := make(chan engine.Event)

  event := engine.Event{
    Code:    "sireneULParser",
    Channel: eventChannel,
  }


  go func() {
    for _, path := range batch.Files["sirene_ul"] {
      tracker := gournal.NewTracker(
        map[string]string{"path": path},
        engine.TrackerReports)

        file, err := os.Open(viper.GetString("APP_DATA") + path)
        if err != nil {
          tracker.Error(err)
          tracker.Report("fatalError")
        }
        event.Info(path + ": ouverture")
        reader := csv.NewReader(file)
        reader.Comma = ','
        reader.LazyQuotes = true

        _, _ = reader.Read()

        for {
          row, err := reader.Read()
          if err == io.EOF {
            break
          } else if err != nil {
            tracker.Error(err)
            event.Critical(tracker.Report("fatalError"))
            break
          }

          if (filter[row[0]]){
            sirene_ul := readLineEtablissement(row, &tracker)
            outputChannel <- sirene_ul
            tracker.Next()
          }
        }
        file.Close()
        event.Info(tracker.Report("abstract"))
      }
      close(outputChannel)
      close(eventChannel)
    }()

    return outputChannel, eventChannel
  }


func readLineEtablissement(row []string, tracker *gournal.Tracker)(SireneUL){
  sirene_ul := SireneUL{}
  sirene_ul.Siren = row[0]
  sirene_ul.RaisonSociale = row[23]
  sirene_ul.CodeStatutJuridique = row[27]
  return sirene_ul
}
