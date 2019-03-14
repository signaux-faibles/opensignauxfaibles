package sirene

import (
	//"bufio"
	"dbmongo/lib/engine"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/chrnin/gournal"
	"github.com/spf13/viper"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren              string     `json:"siren,omitempty" bson:"siren,omitempty"`
	Nic                string     `json:"nic,omitempty" bson:"nic,omitempty"`
	NicSiege           string     `json:"nic_siege,omitempty" bson:"nic_siege,omitempty"`
	RaisonSociale      string     `json:"raison_sociale,omitempty" bson:"raison_sociale,omitempty"`
	NumVoie            string     `json:"numero_voie,omitempty" bson:"numero_voie,omitempty"`
	IndRep             string     `json:"indrep,omitempty" bson:"indrep,omitempty"`
	TypeVoie           string     `json:"type_voie,omitempty" bson:"type_voie,omitempty"`
	CodePostal         string     `json:"code_postal,omitempty" bson:"code_postal,omitempty"`
	Cedex              string     `json:"cedex,omitempty" bson:"cedex,omitempty"`
	Region             string     `json:"region,omitempty" bson:"region,omitempty"`
	Departement        string     `json:"departement,omitempty" bson:"departement,omitempty"`
	Commune            string     `json:"commune,omitempty" bson:"commune,omitempty"`
	APE                string     `json:"ape,omitempty" bson:"ape,omitempty"`
	NatureActivite     string     `json:"nature_activite,omitempty" bson:"nature_activite,omitempty"`
	ActiviteSaisoniere string     `json:"activite_saisoniere,omitempty" bson:"activite_saisoniere,omitempty"`
	ModaliteActivite   string     `json:"modalite_activite,omitempty" bson:"modalite_activite,omitempty"`
	Productif          string     `json:"productif,omitempty" bson:"productif,omitempty"`
	NatureJuridique    string     `json:"nature_juridique,omitempty" bson:"nature_juridique,omitempty"`
	Categorie          string     `json:"categorie,omitempty" bson:"categorie,omitempty"`
	Creation           *time.Time `json:"date_creation,omitempty" bson:"date_creation,omitempty"`
	IndiceMonoactivite *int       `json:"indice_monoactivite,omitempty" bson:"indice_monoactivite,omitempty"`
	TrancheCA          *int       `json:"tranche_ca,omitempty" bson:"tranche_ca,omitempty"`
	Sigle              string     `json:"sigle,omitempty" bson:"sigle,omitempty"`
	Longitude          *float64   `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Lattitude          *float64   `json:"lattitude,omitempty" bson:"lattitude,omitempty"`
	Adresse            [7]string  `json:"adresse" bson:"adresse"`
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
func Parser(batch engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "sireneParser",
		Channel: eventChannel,
	}

	go func() {
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

				if len(row) >= 102 {
					sirene := Sirene{}
					sirene.Siren = row[0]

					sirene.Nic = row[1]
					sirene.NicSiege = row[65]
					sirene.RaisonSociale = row[2]
					sirene.NumVoie = row[16]
					sirene.IndRep = row[17]
					sirene.TypeVoie = row[19]
					sirene.CodePostal = row[20]
					sirene.Cedex = row[21]
					sirene.Region = row[23]
					sirene.Departement = row[24]
					sirene.Commune = row[28]
					sirene.APE = row[42]
					sirene.NatureActivite = row[52]
					sirene.ActiviteSaisoniere = row[55]
					sirene.ModaliteActivite = row[56]
					sirene.Productif = row[57]
					sirene.NatureJuridique = row[71]
					sirene.Categorie = row[82]

					if i, err := time.Parse("20060102", row[50]); err == nil {
						sirene.Creation = &i
					}
					if i, err := strconv.Atoi(row[85]); err == nil {
						sirene.IndiceMonoactivite = &i
					}
					if i, err := strconv.Atoi(row[89]); err == nil {
						sirene.TrancheCA = &i
					}

					sirene.Sigle = row[61]

					if i, err := strconv.ParseFloat(row[100], 64); err == nil {
						sirene.Longitude = &i
					}
					if i, err := strconv.ParseFloat(row[101], 64); err == nil {
						sirene.Lattitude = &i
					}
					sirene.Adresse = [7]string{row[2], row[3], row[4], row[5], row[6], row[7], row[8]}

					outputChannel <- sirene
				} else {
					tracker.Error(errors.New("la ligne ne comporte pas assez de champs"))
					event.Debug(tracker.Report("invalidLine"))
				}
				tracker.Next()
			}
			file.Close()
      event.Info(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}
