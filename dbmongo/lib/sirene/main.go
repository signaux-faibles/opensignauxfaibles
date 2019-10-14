package sirene

import (
	//"bufio"
	"encoding/csv"
	"errors"
	"io"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren       string     `json:"siren,omitempty" bson:"siren,omitempty"`
	Nic         string     `json:"nic,omitempty" bson:"nic,omitempty"`
	Siege       bool       `json:"siege,omitempty" bson:"siege,omitempty"`
	NumVoie     string     `json:"numero_voie,omitempty" bson:"numero_voie,omitempty"`
	IndRep      string     `json:"indrep,omitempty" bson:"indrep,omitempty"`
	TypeVoie    string     `json:"type_voie,omitempty" bson:"type_voie,omitempty"`
	CodePostal  string     `json:"code_postal,omitempty" bson:"code_postal,omitempty"`
	Cedex       string     `json:"cedex,omitempty" bson:"cedex,omitempty"`
	Departement string     `json:"departement,omitempty" bson:"departement,omitempty"`
	Commune     string     `json:"commune,omitempty" bson:"commune,omitempty"`
	APE         string     `json:"ape,omitempty" bson:"ape,omitempty"`
	Creation    *time.Time `json:"date_creation,omitempty" bson:"date_creation,omitempty"`
	Longitude   *float64   `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Lattitude   *float64   `json:"lattitude,omitempty" bson:"lattitude,omitempty"`
	Adresse     [6]string  `json:"adresse" bson:"adresse"`
}

// Key id de l'objet
func (sirene Sirene) Key() string {
	return sirene.Siren + sirene.Nic
}

// Type de données
func (sirene Sirene) Type() string {
	return "sirene"
}

// Scope de l'objet
func (sirene Sirene) Scope() string {
	return "etablissement"
}

// Parser produit les données sirene à partir du fichier geosirene
func Parser(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {

	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "sireneParser",
		Channel: eventChannel,
	}

	go func() {
		dateDebut := batch.Params.DateDebut
		for _, path := range batch.Files["sirene"] {
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

			// _, _ = reader.Read()

			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					event.Critical(tracker.Report("fatalError"))
					break
				}
				filtered, err := marshal.IsFiltered(row[0], cache, batch)
				tracker.Error(err)
				if !filtered {
					notFiltered := (row[40] == "A")
					// Est-ce que l'établissement est intéressant ?
					// = Actif ou a été actif depuis dateDebut
					if row[40] == "F" && row[8] != "" {
						date, err := time.Parse("2006-01-02", row[8][0:10])
						tracker.Error(err)
						notFiltered = date.After(dateDebut)
					}

					if notFiltered {
						sirene := readLineEtablissement(row, &tracker)
						outputChannel <- sirene
						tracker.Next()
					}
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

func readLineEtablissement(row []string, tracker *gournal.Tracker) Sirene {
	sirene := Sirene{}

	sirene.Siren = row[0]

	sirene.Nic = row[1]
	sirene.NumVoie = row[12]
	sirene.IndRep = row[13]
	sirene.TypeVoie = row[14]
	sirene.CodePostal = row[16]
	sirene.Cedex = row[21]
	if len(sirene.CodePostal) >= 2 {
		sirene.Departement = sirene.CodePostal[0:2]
	} else {
		tracker.Error(errors.New("Code postal est manquant ou de format incorrect"))
	}
	sirene.Commune = row[17]
	sirene.APE = strings.Replace(row[45], ".", "", -1)

	loc, _ := time.LoadLocation("Europe/Paris")
	creation, err := time.ParseInLocation("2006-01-02", row[4], loc)
	if err == nil {
		sirene.Creation = &creation
	}
	tracker.Error(err)

	sirene.Siege, err = strconv.ParseBool(row[9])
	tracker.Error(err)

	long, err := strconv.ParseFloat(row[48], 64)
	if err == nil {
		sirene.Longitude = &long
	}
	if row[48] != "" {
		tracker.Error(err)
	}

	latt, err := strconv.ParseFloat(row[49], 64)
	if err == nil {
		sirene.Lattitude = &latt
	}
	if row[49] != "" {
		tracker.Error(err)
	}

	sirene.Adresse = [6]string{row[41], row[11], row[15], row[16], row[17], row[52]}

	return sirene
}
