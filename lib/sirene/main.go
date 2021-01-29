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

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren                string     `col:"siren" json:"siren,omitempty" bson:"siren,omitempty"`
	Nic                  string     `col:"nic" json:"nic,omitempty" bson:"nic,omitempty"`
	Siege                bool       `col:"etablissementSiege" json:"siege,omitempty" bson:"siege,omitempty"`
	ComplementAdresse    string     `col:"complementAdresseEtablissement" json:"complement_adresse,omitempty" bson:"complement_adresse,omitempty"`
	NumVoie              string     `col:"numeroVoieEtablissement" json:"numero_voie,omitempty" bson:"numero_voie,omitempty"`
	IndRep               string     `col:"indiceRepetitionEtablissement" json:"indrep,omitempty" bson:"indrep,omitempty"`
	TypeVoie             string     `col:"typeVoieEtablissement" json:"type_voie,omitempty" bson:"type_voie,omitempty"`
	Voie                 string     `col:"libelleVoieEtablissement" json:"voie,omitempty" bson:"voie,omitempty"`
	Commune              string     `col:"libelleCommuneEtablissement" json:"commune,omitempty" bson:"commune,omitempty"`
	CommuneEtranger      string     `col:"libelleCommuneEtrangerEtablissement" json:"commune_etranger,omitempty" bson:"commune_etranger,omitempty"`
	DistributionSpeciale string     `col:"distributionSpecialeEtablissement" json:"distribution_speciale,omitempty" bson:"distribution_speciale,omitempty"`
	CodeCommune          string     `col:"codeCommuneEtablissement" json:"code_commune,omitempty" bson:"code_commune,omitempty"`
	CodeCedex            string     `col:"codeCedexEtablissement" json:"code_cedex,omitempty" bson:"code_cedex,omitempty"`
	Cedex                string     `col:"libelleCedexEtablissement" json:"cedex,omitempty" bson:"cedex,omitempty"`
	CodePaysEtranger     string     `col:"codePaysEtrangerEtablissement" json:"code_pays_etranger,omitempty" bson:"code_pays_etranger,omitempty"`
	PaysEtranger         string     `col:"libellePaysEtrangerEtablissement" json:"pays_etranger,omitempty" bson:"pays_etranger,omitempty"`
	CodePostal           string     `col:"codePostalEtablissement" json:"code_postal,omitempty" bson:"code_postal,omitempty"`
	Departement          string     `json:"departement,omitempty" bson:"departement,omitempty"`
	APE                  string     `json:"ape,omitempty" bson:"ape,omitempty"`
	CodeActivite         string     `col:"activitePrincipaleEtablissement" json:"code_activite,omitempty" bson:"code_activite,omitempty"`
	NomenActivite        string     `col:"nomenclatureActivitePrincipaleEtablissement" json:"nomen_activite,omitempty" bson:"nomen_activite,omitempty"`
	Creation             *time.Time `col:"dateCreationEtablissement" json:"date_creation,omitempty" bson:"date_creation,omitempty"`
	Longitude            float64    `col:"longitude" json:"longitude,omitempty" bson:"longitude,omitempty"`
	Latitude             float64    `col:"latitude" json:"latitude,omitempty" bson:"latitude,omitempty"`
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
	file     *os.File
	reader   *csv.Reader
	colIndex marshal.ColMapping
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
	parser.colIndex, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Sirene{})
	return err
}

func (parser *sireneParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseLine(parser.colIndex, row, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseLine(idx marshal.ColMapping, row []string, parsedLine *marshal.ParsedLineResult) {
	var err error
	idxRow := idx.IndexRow(row)
	sirene := Sirene{}
	sirene.Siren = idxRow.GetVal("siren")
	sirene.Nic = idxRow.GetVal("nic")
	sirene.Siege, err = strconv.ParseBool(idxRow.GetVal("etablissementSiege"))
	parsedLine.AddRegularError(err)

	sirene.ComplementAdresse = idxRow.GetVal("complementAdresseEtablissement")
	sirene.NumVoie = idxRow.GetVal("numeroVoieEtablissement")
	sirene.IndRep = indRep[idxRow.GetVal("indiceRepetitionEtablissement")]
	sirene.TypeVoie = typeVoie[idxRow.GetVal("typeVoieEtablissement")]
	sirene.Voie = idxRow.GetVal("libelleVoieEtablissement")
	sirene.Commune = idxRow.GetVal("libelleCommuneEtablissement")
	sirene.CommuneEtranger = idxRow.GetVal("libelleCommuneEtrangerEtablissement")
	sirene.DistributionSpeciale = idxRow.GetVal("distributionSpecialeEtablissement")
	sirene.CodeCommune = idxRow.GetVal("codeCommuneEtablissement")
	sirene.CodeCedex = idxRow.GetVal("codeCedexEtablissement")
	sirene.Cedex = idxRow.GetVal("libelleCedexEtablissement")
	sirene.CodePaysEtranger = idxRow.GetVal("codePaysEtrangerEtablissement")
	sirene.PaysEtranger = idxRow.GetVal("libellePaysEtrangerEtablissement")

	if len(idxRow.GetVal("codePostalEtablissement")) > 2 {
		sirene.CodePostal = idxRow.GetVal("codePostalEtablissement")
		departement := idxRow.GetVal("codePostalEtablissement")[0:2]
		// traitement pour les départements de Corse
		if idxRow.GetVal("codePostalEtablissement")[0:3] == "201" || idxRow.GetVal("codePostalEtablissement")[0:3] == "200" {
			departement = "2A"
		} else if idxRow.GetVal("codePostalEtablissement")[0:2] == "20" {
			departement = "2B"
		}
		sirene.Departement = departement
	} else {
		parsedLine.AddRegularError(errors.New("Code postal est manquant ou de format incorrect"))
	}

	if idxRow.GetVal("activitePrincipaleEtablissement") != "" {
		if idxRow.GetVal("nomenclatureActivitePrincipaleEtablissement") == "NAFRev2" {
			ape := strings.Replace(idxRow.GetVal("activitePrincipaleEtablissement"), ".", "", -1)
			if matched, err := regexp.MatchString(`^[0-9]{4}[A-Z]$`, ape); err == nil && matched {
				sirene.APE = ape
			}
		} else {
			sirene.CodeActivite = idxRow.GetVal("activitePrincipaleEtablissement")
			sirene.NomenActivite = idxRow.GetVal("nomenclatureActivitePrincipaleEtablissement")
		}
	}

	creation, err := time.Parse("2006-01-02", idxRow.GetVal("dateCreationEtablissement")) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sirene.Creation = &creation
	}

	long, err := strconv.ParseFloat(idxRow.GetVal("longitude"), 64)
	if err == nil {
		sirene.Longitude = long
	}
	if row[48] != "" {
		parsedLine.AddRegularError(err)
	}

	lat, err := strconv.ParseFloat(idxRow.GetVal("latitude"), 64)
	if err == nil {
		sirene.Latitude = lat
	}
	if row[49] != "" {
		parsedLine.AddRegularError(err)
	}
	parsedLine.AddTuple(sirene)
}
