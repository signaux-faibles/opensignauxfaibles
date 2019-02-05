package main

import (
  "github.com/kshedden/datareader"
//  "github.com/davecgh/go-spew/spew"
  "fmt"
	"os"
  "time"
	"github.com/spf13/viper"
  "regexp"
	"github.com/cnf/structhash"
)

// Interim Interim – fichier DARES 
type Interim struct {
	Siret        string    `json:"siret" bson:"siret"`
	Periode      time.Time `json:"periode" bson:"periode"`
	ETP          float64   `json:"etp" bson:"etp"`
}

func parseInterim(paths []string) chan *Interim {
	outputChannel := make(chan *Interim)

	field := map[string]int{
		"Siret": 0,
		"Periode": 1,
		"ETP": 4,
	}

	go func() {
		for _, path := range paths {
			// get current file name
      file, err := os.Open(viper.GetString("APP_DATA") + path)

			if err != nil {
				journal(critical, "importInterim", "Erreur à l'ouverture du fichier "+path+": "+err.Error())
      } else {

        e:= 0
        n := 0
        journal(info, "importInterim", "Ouverture du fichier "+path)

        reader, _ := datareader.NewSAS7BDATReader(file)

        row, err := reader.Read(-1)
        if err != nil{
          journal(critical, "importInterim", "Erreur à la lecture du fichier "+path+": "+err.Error())
        }

        sirets, missing, _ := row[field["Siret"]].AsStringSlice()
        periode, _, _ :=  row[field["Periode"]].AsFloat64Slice()
        etp, _, _ := row[field["ETP"]].AsFloat64Slice()
        for i := 0; i < len(sirets); i++ {
          interim := Interim{}
          validSiret, _ := regexp.MatchString("[0-9]{14}", sirets[i])
          n++
          if !missing[i] && validSiret {
            interim.Siret = sirets[i][:14]
            interim.Periode, _ = time.Parse("20060102", fmt.Sprintf("%6.0f", periode[i]) + "01")
            interim.ETP =  etp[i]
            outputChannel <- &interim
          } else {
            e++
          }

        }
        journal(debug, "importInterim", "Import du fichier interim "+path+" terminé. "+fmt.Sprint(n)+" lignes traitée(s), "+fmt.Sprint(e)+" rejet(s)")
        file.Close()
      }
    }

    close(outputChannel)
  }()
  return outputChannel
}

func importInterim(batch *AdminBatch) error {
  journal(info, "importInterim", "Import du batch "+batch.ID.Key+": Interim")

  for interim := range parseInterim(batch.Files["interim"]) {
    hash := fmt.Sprintf("%x", structhash.Md5(interim, 1))

    value := Value{
      Value: Data{
        Scope: "etablissement",
        Key: interim.Siret,
        Batch: map[string]Batch{
          batch.ID.Key: Batch{
            Interim: map[string]*Interim{
              hash: interim,
            },
          },
        },
      },
    }
    db.ChanData <- &value
  }
  db.ChanData <- &Value{}
  journal(info, "importInterim", "Fin de l'import du batch "+batch.ID.Key+": Interim")

  return nil
}
