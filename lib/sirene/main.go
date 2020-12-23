package sirene

import (
	//"bufio"

	"encoding/csv"
	"errors"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
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

var fields = []string{
	"siren",
	"nic",
	"siret",
	"statutDiffusionEtablissement",
	"dateCreationEtablissement",
	"trancheEffectifsEtablissement",
	"anneeEffectifsEtablissement",
	"activitePrincipaleRegistreMetiersEtablissement",
	"dateDernierTraitementEtablissement",
	"etablissementSiege",
	"nombrePeriodesEtablissement",
	"complementAdresseEtablissement",
	"numeroVoieEtablissement",
	"indiceRepetitionEtablissement",
	"typeVoieEtablissement",
	"libelleVoieEtablissement",
	"codePostalEtablissement",
	"libelleCommuneEtablissement",
	"libelleCommuneEtrangerEtablissement",
	"distributionSpecialeEtablissement",
	"codeCommuneEtablissement",
	"codeCedexEtablissement",
	"libelleCedexEtablissement",
	"codePaysEtrangerEtablissement",
	"libellePaysEtrangerEtablissement",
	"complementAdresse2Etablissement",
	"numeroVoie2Etablissement",
	"indiceRepetition2Etablissement",
	"typeVoie2Etablissement",
	"libelleVoie2Etablissement",
	"codePostal2Etablissement",
	"libelleCommune2Etablissement",
	"libelleCommuneEtranger2Etablissement",
	"distributionSpeciale2Etablissement",
	"codeCommune2Etablissement",
	"codeCedex2Etablissement",
	"libelleCedex2Etablissement",
	"codePaysEtranger2Etablissement",
	"libellePaysEtranger2Etablissement",
	"dateDebut",
	"etatAdministratifEtablissement",
	"enseigne1Etablissement",
	"enseigne2Etablissement",
	"enseigne3Etablissement",
	"denominationUsuelleEtablissement",
	"activitePrincipaleEtablissement",
	"nomenclatureActivitePrincipaleEtablissement",
	"caractereEmployeurEtablissement",
	"longitude",
	"latitude",
	"geo_score",
	"geo_type",
	"geo_adresse",
	"geo_id",
	"geo_ligne",
	"geo_l4",
	"geo_l5",
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

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &sireneParser{}

type sireneParser struct {
	file   *os.File
	reader *csv.Reader
}

func (parser *sireneParser) GetFileType() string {
	return "sirene"
}

func (parser *sireneParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *sireneParser) Close() error {
	return parser.file.Close()
}

func (parser *sireneParser) Open(filePath string) (err error) {
	parser.file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	parser.reader = csv.NewReader(parser.file)
	parser.reader.Comma = ','
	parser.reader.LazyQuotes = true

	// parse header
	row, err := parser.reader.Read()
	if err != nil {
		return err // may be io.EOF
	} else if !reflect.DeepEqual(row, fields) {
		return errors.New("sirene header does not match the parser's expectations")
	}

	return nil
}

func (parser *sireneParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	f := marshal.GetFieldBindings(fields)
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseLine(f, row, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseLine(f map[string]int, row []string, parsedLine *marshal.ParsedLineResult) {
	var err error
	sirene := Sirene{}
	sirene.Siren = row[f["siren"]]
	sirene.Nic = row[f["nic"]]
	sirene.Siege, err = strconv.ParseBool(row[f["etablissementSiege"]])
	parsedLine.AddRegularError(err)

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
		parsedLine.AddRegularError(errors.New("Code postal est manquant ou de format incorrect"))
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

	creation, err := time.Parse("2006-01-02", row[f["dateCreationEtablissement"]]) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sirene.Creation = &creation
	}

	long, err := strconv.ParseFloat(row[f["longitude"]], 64)
	if err == nil {
		sirene.Longitude = long
	}
	if row[48] != "" {
		parsedLine.AddRegularError(err)
	}

	lat, err := strconv.ParseFloat(row[f["latitude"]], 64)
	if err == nil {
		sirene.Latitude = lat
	}
	if row[49] != "" {
		parsedLine.AddRegularError(err)
	}
	parsedLine.AddTuple(sirene)
}
