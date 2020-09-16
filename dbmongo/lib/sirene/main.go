package sirene

import (
	//"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren                string     `json:"siren,omitempty" bson:"siren,omitempty"`
	Nic                  string     `json:"nic,omitempty" bson:"nic,omitempty"`
	Siege                bool       `json:"siege,omitempty" bson:"siege,omitempty"`
	NumVoie              *string    `json:"numero_voie,omitempty" bson:"numero_voie,omitempty"`
	IndRep               *string    `json:"indrep,omitempty" bson:"indrep,omitempty"`
	TypeVoie             *string    `json:"type_voie,omitempty" bson:"type_voie,omitempty"`
	Voie                 *string    `json:"voie,omitempty" bson:"voie,omitempty"`
	CodePostal           *string    `json:"code_postal,omitempty" bson:"code_postal,omitempty"`
	Cedex                *string    `json:"cedex,omitempty" bson:"cedex,omitempty"`
	Departement          *string    `json:"departement,omitempty" bson:"departement,omitempty"`
	Commune              *string    `json:"commune,omitempty" bson:"commune,omitempty"`
	DistributionSpeciale *string    `json:"distribution_speciale,omitempty" bson:"distribution_speciale,omitempty"`
	APE                  *string    `json:"ape,omitempty" bson:"ape,omitempty"`
	CodeActivite         *string    `json:"code_activite,omitempty" bson:"code_activite,omitempty"`
	NomenActivite        *string    `json:"nomen_activite,omitempty" bson:"nomen_activite,omitempty"`
	Creation             *time.Time `json:"date_creation,omitempty" bson:"date_creation,omitempty"`
	Longitude            *float64   `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Latitude             *float64   `json:"latitude,omitempty" bson:"latitude,omitempty"`
}

var f = map[string]int{
	"siren":                         0,
	"nic":                           1,
	"siret":                         2,
	"statutDiffusionEtablissement":  3,
	"dateCreationEtablissement":     4,
	"trancheEffectifsEtablissement": 5,
	"anneeEffectifsEtablissement":   6,
	"activitePrincipaleRegistreMetiersEtablissement": 7,
	"dateDernierTraitementEtablissement":             8,
	"etablissementSiege":                             9,
	"nombrePeriodesEtablissement":                    10,
	"complementAdresseEtablissement":                 11,
	"numeroVoieEtablissement":                        12,
	"indiceRepetitionEtablissement":                  13,
	"typeVoieEtablissement":                          14,
	"libelleVoieEtablissement":                       15,
	"codePostalEtablissement":                        16,
	"libelleCommuneEtablissement":                    17,
	"libelleCommuneEtrangerEtablissement":            18,
	"distributionSpecialeEtablissement":              19,
	"codeCommuneEtablissement":                       20,
	"codeCedexEtablissement":                         21,
	"libelleCedexEtablissement":                      22,
	"codePaysEtrangerEtablissement":                  23,
	"libellePaysEtrangerEtablissement":               24,
	"complementAdresse2Etablissement":                25,
	"numeroVoie2Etablissement":                       26,
	"indiceRepetition2Etablissement":                 27,
	"typeVoie2Etablissement":                         28,
	"libelleVoie2Etablissement":                      29,
	"codePostal2Etablissement":                       30,
	"libelleCommune2Etablissement":                   31,
	"libelleCommuneEtranger2Etablissement":           32,
	"distributionSpeciale2Etablissement":             33,
	"codeCommune2Etablissement":                      34,
	"codeCedex2Etablissement":                        35,
	"libelleCedex2Etablissement":                     36,
	"codePaysEtranger2Etablissement":                 37,
	"libellePaysEtranger2Etablissement":              38,
	"dateDebut":                                      39,
	"etatAdministratifEtablissement":                 40,
	"enseigne1Etablissement":                         41,
	"enseigne2Etablissement":                         42,
	"enseigne3Etablissement":                         43,
	"denominationUsuelleEtablissement":               44,
	"activitePrincipaleEtablissement":                45,
	"nomenclatureActivitePrincipaleEtablissement":    46,
	"caractereEmployeurEtablissement":                47,
	"longitude":                                      48,
	"latitude":                                       49,
	"geo_score":                                      50,
	"geo_type":                                       51,
	"geo_adresse":                                    52,
	"geo_id":                                         53,
	"geo_ligne":                                      54,
	"geo_l4":                                         55,
	"geo_l5":                                         56,
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
					sirene := readLineEtablissement(row, &tracker)
					outputChannel <- sirene
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

func readLineEtablissement(row []string, tracker *gournal.Tracker) Sirene {
	sirene := Sirene{}

	sirene.Siren = row[f["siren"]]
	sirene.Nic = row[f["nic"]]

	if row[f["numeroVoieEtablissement"]] != "" {
		sirene.NumVoie = &row[f["numeroVoieEtablissement"]]
	}

	if row[f["indiceRepetitionEtablissement"]] != "" {
		sirene.IndRep = &row[f["indiceRepetitionEtablissement"]]
	}

	if row[f["typeVoieEtablissement"]] != "" {
		sirene.TypeVoie = &row[f["typeVoieEtablissement"]]
	}

	if row[f["libelleVoieEtablissement"]] != "" {
		sirene.Voie = &row[f["libelleVoieEtablissement"]]
	}

	if len(row[f["codePostalEtablissement"]]) > 2 {
		sirene.CodePostal = &row[f["codePostalEtablissement"]]
		departement := row[f["codePostalEtablissement"]][0:2]
		if row[f["codePostalEtablissement"]][0:3] == "201" || row[f["codePostalEtablissement"]][0:3] == "200" {
			departement = "2A"
		} else if row[f["codePostalEtablissement"]][0:2] == "20" {
			departement = "2B"
		}
		sirene.Departement = &departement
	} else {
		tracker.Error(errors.New("Code postal est manquant ou de format incorrect"))
	}

	if row[f["codeCedexEtablissement"]] != "" {
		sirene.Cedex = &row[f["codeCedexEtablissement"]]
	}

	if row[f["distributionSpecialeEtablissement"]] != "" {
		sirene.DistributionSpeciale = &row[f["distributionSpecialeEtablissement"]]
	}

	if row[f["libelleCommuneEtablissement"]] != "" {
		sirene.Commune = &row[f["libelleCommuneEtablissement"]]
	}

	if row[f["activitePrincipaleEtablissement"]] != "" {
		if row[f["nomenclatureActivitePrincipaleEtablissement"]] == "NAFRev2" {
			ape := strings.Replace(row[f["activitePrincipaleEtablissement"]], ".", "", -1)
			if matched, err := regexp.MatchString(`^[0-9]{4}[A-Z]$`, ape); err == nil && matched {
				sirene.APE = &ape
			}
		} else {
			sirene.CodeActivite = &row[f["activitePrincipaleEtablissement"]]
			sirene.NomenActivite = &row[f["nomenclatureActivitePrincipaleEtablissement"]]
		}
	}

	loc, _ := time.LoadLocation("Europe/Paris")
	creation, err := time.ParseInLocation("2006-01-02", row[f["dateCreationEtablissement"]], loc)
	if err == nil {
		sirene.Creation = &creation
	}
	tracker.Error(err)

	sirene.Siege, err = strconv.ParseBool(row[f["etablissementSiege"]])
	tracker.Error(err)

	long, err := strconv.ParseFloat(row[f["longitude"]], 64)
	if err == nil {
		sirene.Longitude = &long
	}
	if row[48] != "" {
		tracker.Error(err)
	}

	lat, err := strconv.ParseFloat(row[f["latitude"]], 64)
	if err == nil {
		sirene.Latitude = &lat
	}
	if row[49] != "" {
		tracker.Error(err)
	}

	return sirene
}
