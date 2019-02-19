package sirene

import (
	"bufio"
	"dbmongo/lib/engine"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/chrnin/gournal"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren              string     `json,omitempty:"siren" bson,omitempty:"siren"`
	Nic                string     `json,omitempty:"nic" bson,omitempty:"nic"`
	NicSiege           string     `json,omitempty:"nic_siege" bson,omitempty:"nic_siege"`
	RaisonSociale      string     `json,omitempty:"raison_sociale" bson,omitempty:"raison_sociale"`
	NumVoie            string     `json,omitempty:"numero_voie" bson,omitempty:"numero_voie"`
	IndRep             string     `json,omitempty:"indrep" bson,omitempty:"indrep"`
	TypeVoie           string     `json,omitempty:"type_voie" bson,omitempty:"type_voie"`
	CodePostal         string     `json,omitempty:"code_postal" bson,omitempty:"code_postal"`
	Cedex              string     `json,omitempty:"cedex" bson,omitempty:"cedex"`
	Region             string     `json,omitempty:"region" bson,omitempty:"region"`
	Departement        string     `json,omitempty:"departement" bson,omitempty:"departement"`
	Commune            string     `json,omitempty:"commune" bson,omitempty:"commune"`
	APE                string     `json,omitempty:"ape" bson,omitempty:"ape"`
	NatureActivite     string     `json,omitempty:"nature_activite" bson,omitempty:"nature_activite"`
	ActiviteSaisoniere string     `json,omitempty:"activite_saisoniere" bson,omitempty:"activite_sai"`
	ModaliteActivite   string     `json,omitempty:"modalite_activite" bson,omitempty:"modalite_activite"`
	Productif          string     `json,omitempty:"productif" bson,omitempty:"productif"`
	NatureJuridique    string     `json,omitempty:"nature_juridique" bson,omitempty:"nature_juridique"`
	Categorie          string     `json,omitempty:"categorie" bson,omitempty:"categorie"`
	Creation           *time.Time `json,omitempty:"date_creation" bson,omitempty:"date_creation"`
	IndiceMonoactivite *int       `json,omitempty:"indice_monoactivite" bson,omitempty:"indice_monoactivite"`
	TrancheCA          *int       `json,omitempty:"tranche_ca" bson,omitempty:"tranche_ca"`
	Sigle              string     `json,omitempty:"sigle" bson,omitempty:"sigle"`
	Longitude          *float64   `json,omitempty:"longitude" bson:"longitude"`
	Lattitude          *float64   `json,omitempty:"lattitude" bson:"lattitude"`
	Adresse            [7]string  `json:"adresse" bson:"adresse"`
}

// Key id de l'objet
func (sirene Sirene) Key() string {
	return sirene.Siren + sirene.Nic
}

// Type de données
func (sirene Sirene) Type() string {
	return "bdf"
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
		for _, path := range batch.Files["engine"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(path)
			if err != nil {
				tracker.Error(err)
				tracker.Report("fatalError")
			}
			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ','

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
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}
