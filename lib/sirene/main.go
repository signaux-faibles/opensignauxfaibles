package sirene

import (
	//"bufio"

	"encoding/csv"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren                string     `input:"siren"                                       json:"siren,omitempty"                 csv:"siren"`
	Nic                  string     `input:"nic"                                         json:"nic,omitempty"                   csv:"nic"`
	Siege                bool       `input:"etablissementSiege"                          json:"siege,omitempty"                 csv:"siege"`
	ComplementAdresse    string     `input:"complementAdresseEtablissement"              json:"complement_adresse,omitempty"    csv:"complement_adresse"`
	NumVoie              string     `input:"numeroVoieEtablissement"                     json:"numero_voie,omitempty"           csv:"numéro_voie"`
	IndRep               string     `input:"indiceRepetitionEtablissement"               json:"indrep,omitempty"                csv:"indrep"`
	TypeVoie             string     `input:"typeVoieEtablissement"                       json:"type_voie,omitempty"             csv:"type_voie"`
	Voie                 string     `input:"libelleVoieEtablissement"                    json:"voie,omitempty"                  csv:"voie"`
	Commune              string     `input:"libelleCommuneEtablissement"                 json:"commune,omitempty"               csv:"commune"`
	CommuneEtranger      string     `input:"libelleCommuneEtrangerEtablissement"         json:"commune_etranger,omitempty"      csv:"commune_étranger"`
	DistributionSpeciale string     `input:"distributionSpecialeEtablissement"           json:"distribution_speciale,omitempty" csv:"distribution_speciale"`
	CodeCommune          string     `input:"codeCommuneEtablissement"                    json:"code_commune,omitempty"          csv:"code_commune"`
	CodeCedex            string     `input:"codeCedexEtablissement"                      json:"code_cedex,omitempty"            csv:"code_cedex"`
	Cedex                string     `input:"libelleCedexEtablissement"                   json:"cedex,omitempty"                 csv:"cedex"`
	CodePaysEtranger     string     `input:"codePaysEtrangerEtablissement"               json:"code_pays_etranger,omitempty"    csv:"code_pays_étranger"`
	PaysEtranger         string     `input:"libellePaysEtrangerEtablissement"            json:"pays_etranger,omitempty"         csv:"pays_étranger"`
	CodePostal           string     `input:"codePostalEtablissement"                     json:"code_postal,omitempty"           csv:"code_postal"`
	Departement          string     `                                                    json:"departement,omitempty"           csv:"département"`
	APE                  string     `                                                    json:"ape,omitempty"                   csv:"ape"`
	CodeActivite         string     `input:"activitePrincipaleEtablissement"             json:"code_activite,omitempty"         csv:"code_activité"`
	NomenActivite        string     `input:"nomenclatureActivitePrincipaleEtablissement" json:"nomen_activite,omitempty"        csv:"nomenclature_activité"`
	Creation             *time.Time `input:"dateCreationEtablissement"                   json:"date_creation,omitempty"         csv:"création"`
	Longitude            float64    `input:"longitude"                                   json:"longitude,omitempty"             csv:"longitude"`
	Latitude             float64    `input:"latitude"                                    json:"latitude,omitempty"              csv:"latitude"`
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

func (parser *sireneParser) Type() string {
	return "sirene"
}

func (parser *sireneParser) Init(_ *marshal.Cache, _ *base.AdminBatch) error {
	return nil
}

func (parser *sireneParser) Close() error {
	return parser.file.Close()
}

func (parser *sireneParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ',', true)
	if err == nil {
		parser.colIndex, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Sirene{})
	}
	return err
}

func (parser *sireneParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseLine(parser.colIndex, row, parsedLine)
	})
}

func parseLine(idx marshal.ColMapping, row []string, parsedLine *marshal.ParsedLineResult) {
	var err error
	idxRow := idx.IndexRow(row)
	sirene := Sirene{}
	sirene.Siren = idxRow.GetVal("siren")
	sirene.Nic = idxRow.GetVal("nic")
	sirene.Siege, err = idxRow.GetBool("etablissementSiege")
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
		parsedLine.AddRegularError(errors.New("code postal est manquant ou de format incorrect"))
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

	if val, ok := idxRow.GetOptionalVal("longitude"); ok && val != "" {
		sirene.Longitude, err = strconv.ParseFloat(val, 64)
		parsedLine.AddRegularError(err)
	}

	if val, ok := idxRow.GetOptionalVal("latitude"); ok && val != "" {
		sirene.Latitude, err = strconv.ParseFloat(val, 64)
		parsedLine.AddRegularError(err)
	}

	parsedLine.AddTuple(sirene)
}
