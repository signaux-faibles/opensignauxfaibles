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

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren                string     `json:"siren,omitempty" bson:"siren,omitempty"`
	Nic                  string     `json:"nic,omitempty" bson:"nic,omitempty"`
	Siege                bool       `json:"siege,omitempty" bson:"siege,omitempty"`
	ComplementAdresse    string     `json:"complement_adresse,omitempty" bson:"complement_adresse,omitempty"`
	NumVoie              string     `json:"numero_voie,omitempty" bson:"numero_voie,omitempty"`
	IndRep               string     `json:"indrep,omitempty" bson:"indrep,omitempty"`
	TypeVoie             string     `json:"type_voie,omitempty" bson:"type_voie,omitempty"`
	Voie                 string     `json:"voie,omitempty" bson:"voie,omitempty"`
	Commune              string     `json:"commune,omitempty" bson:"commune,omitempty"`
	CommuneEtranger      string     `json:"commune_etranger,omitempty" bson:"commune_etranger,omitempty"`
	DistributionSpeciale string     `json:"distribution_speciale,omitempty" bson:"distribution_speciale,omitempty"`
	CodeCommune          string     `json:"code_commune,omitempty" bson:"code_commune,omitempty"`
	CodeCedex            string     `json:"code_cedex,omitempty" bson:"code_cedex,omitempty"`
	Cedex                string     `json:"cedex,omitempty" bson:"cedex,omitempty"`
	CodePaysEtranger     string     `json:"code_pays_etranger,omitempty" bson:"code_pays_etranger,omitempty"`
	PaysEtranger         string     `json:"pays_etranger,omitempty" bson:"pays_etranger,omitempty"`
	CodePostal           string     `json:"code_postal,omitempty" bson:"code_postal,omitempty"`
	Departement          string     `json:"departement,omitempty" bson:"departement,omitempty"`
	APE                  string     `json:"ape,omitempty" bson:"ape,omitempty"`
	CodeActivite         string     `json:"code_activite,omitempty" bson:"code_activite,omitempty"`
	NomenActivite        string     `json:"nomen_activite,omitempty" bson:"nomen_activite,omitempty"`
	Creation             *time.Time `json:"date_creation,omitempty" bson:"date_creation,omitempty"`
	Longitude            float64    `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Latitude             float64    `json:"latitude,omitempty" bson:"latitude,omitempty"`
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

var typeVoie = map[string]string{
	"ALL":  "ALLÉE",
	"AV":   "AVENUE",
	"BD":   "BOULEVARD",
	"CAR":  "CARREFOUR",
	"CHE":  "CHEMIN",
	"CHS":  "CHAUSSÉE",
	"CITE": "CITÉ",
	"COR":  "CORNICHE",
	"CRS":  "COURS",
	"DOM":  "DOMAINE",
	"DSC":  "DESCENTE",
	"ECA":  "ECART",
	"ESP":  "ESPLANADE",
	"FG":   "FAUBOURG",
	"GR":   "GRANDE RUE",
	"HAM":  "HAMEAU",
	"HLE":  "HALLE",
	"IMP":  "IMPASSE",
	"LD":   "LIEU DIT",
	"LOT":  "LOTISSEMENT",
	"MAR":  "MARCHÉ",
	"MTE":  "MONTÉE",
	"PAS":  "PASSAGE",
	"PL":   "PLACE",
	"PLN":  "PLAINE",
	"PLT":  "PLATEAU",
	"PRO":  "PROMENADE",
	"PRV":  "PARVIS",
	"QUA":  "QUARTIER",
	"QUAI": "QUAI",
	"RES":  "RÉSIDENCE",
	"RLE":  "RUELLE",
	"ROC":  "ROCADE",
	"RPT":  "ROND POINT",
	"RTE":  "ROUTE",
	"RUE":  "RUE",
	"SEN":  "SENTE - SENTIER",
	"SQ":   "SQUARE",
	"TPL":  "TERRE-PLEIN",
	"TRA":  "TRAVERSE",
	"VLA":  "VILLA",
	"VLGE": "VILLAGE",
}

var indRep = map[string]string{
	"B": "BIS",
	"T": "TER",
	"Q": "QUATER",
	"C": "QUINQUIES",
}

// Key id de l'objet",
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
func Parser(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {

	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)

	event := marshal.Event{
		Code:    "sireneParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["sirene"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				marshal.TrackerReports)

			event.Info(path + ": ouverture")
			ParseFile(viper.GetString("APP_DATA")+path, &cache, batch, &tracker, outputChannel)
			event.Info(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}

// ParseFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	filter := marshal.GetSirenFilterFromCache(*cache)
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	parseSireneFile(reader, filter, tracker, outputChannel)
}

func parseSireneFile(reader *csv.Reader, filter map[string]bool, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else {
			validSiren := sfregexp.RegexpDict["siren"].MatchString(row[f["siren"]])
			if !validSiren {
				tracker.Add(errors.New("siren invalide : " + row[f["siren"]]))
			} else {
				filtered, err := marshal.IsFiltered(row[f["siren"]], filter)
				tracker.Add(err)
				if !filtered {
					outputChannel <- parseSireneLine(row, tracker)
				}
			}
		}
		tracker.Next()
	}
}

func parseSireneLine(row []string, tracker *gournal.Tracker) Sirene {
	sirene := Sirene{}
	var err error
	sirene.Siren = row[f["siren"]]
	sirene.Nic = row[f["nic"]]
	sirene.Siege, err = strconv.ParseBool(row[f["etablissementSiege"]])
	tracker.Add(err)

	sirene.ComplementAdresse = row[f["complementAdresseEtablissement"]]
	sirene.NumVoie = row[f["numeroVoieEtablissement"]]
	sirene.IndRep = indRep[row[f["indiceRepetitionEtablissement"]]]
	sirene.TypeVoie = typeVoie[row[f["typeVoieEtablissement"]]]
	sirene.Voie = row[f["libelleVoieEtablissement"]]
	sirene.Commune = row[f["libelleCommuneEtablissement"]]
	sirene.CommuneEtranger = row[f["libelleCommuneEtrangerEtablissement"]]
	sirene.DistributionSpeciale = row[f["distributionSpecialeEtablissement"]]
	sirene.CodeCommune = row[f["codeCommuneEtablissement"]]
	sirene.CodeCedex = row[f["codeCedexEtablissement"]]
	sirene.Cedex = row[f["libelleCedexEtablissement"]]
	sirene.CodePaysEtranger = row[f["codePaysEtrangerEtablissement"]]
	sirene.PaysEtranger = row[f["libellePaysEtrangerEtablissement"]]

	if len(row[f["codePostalEtablissement"]]) > 2 {
		sirene.CodePostal = row[f["codePostalEtablissement"]]
		departement := row[f["codePostalEtablissement"]][0:2]
		// traitement pour les départements de Corse
		if row[f["codePostalEtablissement"]][0:3] == "201" || row[f["codePostalEtablissement"]][0:3] == "200" {
			departement = "2A"
		} else if row[f["codePostalEtablissement"]][0:2] == "20" {
			departement = "2B"
		}
		sirene.Departement = departement
	} else {
		tracker.Add(errors.New("Code postal est manquant ou de format incorrect"))
	}

	if row[f["activitePrincipaleEtablissement"]] != "" {
		if row[f["nomenclatureActivitePrincipaleEtablissement"]] == "NAFRev2" {
			ape := strings.Replace(row[f["activitePrincipaleEtablissement"]], ".", "", -1)
			if matched, err := regexp.MatchString(`^[0-9]{4}[A-Z]$`, ape); err == nil && matched {
				sirene.APE = ape
			}
		} else {
			sirene.CodeActivite = row[f["activitePrincipaleEtablissement"]]
			sirene.NomenActivite = row[f["nomenclatureActivitePrincipaleEtablissement"]]
		}
	}

	loc, _ := time.LoadLocation("Europe/Paris")
	creation, err := time.ParseInLocation("2006-01-02", row[f["dateCreationEtablissement"]], loc)
	if err == nil {
		sirene.Creation = &creation
	}
	tracker.Add(err)

	long, err := strconv.ParseFloat(row[f["longitude"]], 64)
	if err == nil {
		sirene.Longitude = long
	}
	if row[48] != "" {
		tracker.Add(err)
	}

	lat, err := strconv.ParseFloat(row[f["latitude"]], 64)
	if err == nil {
		sirene.Latitude = lat
	}
	if row[49] != "" {
		tracker.Add(err)
	}

	return sirene
}
