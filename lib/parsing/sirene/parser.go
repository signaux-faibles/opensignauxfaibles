// Package siren fournit les parser pour extraire les données des fichiers
// Sirene
package sirene

import (
	//"bufio"

	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type SireneParser struct{}

func NewSireneParser() engine.Parser {
	return &SireneParser{}
}

func (p *SireneParser) Type() base.ParserType { return base.Sirene }
func (p *SireneParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:     r,
		RowParser:  &sireneRowParser{},
		Comma:      ',',
		LazyQuotes: false,
		DestTuple:  Sirene{},
	}
}

type sireneRowParser struct{}

func (rp *sireneRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) error {
	var err error

	idxRow := idx.IndexRow(row)
	sirene := Sirene{}
	sirene.Siren = idxRow.GetVal("siren")
	sirene.Nic = idxRow.GetVal("nic")
	sirene.Siret = sirene.Siren + sirene.Nic

	sirene.Siege, err = idxRow.GetBool("etablissementSiege")
	res.AddRegularError(err)

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
		res.AddRegularError(errors.New("code postal est manquant ou de format incorrect"))
	}

	if idxRow.GetVal("activitePrincipaleEtablissement") != "" {
		if idxRow.GetVal("nomenclatureActivitePrincipaleEtablissement") == "NAFRev2" {
			ape := strings.ReplaceAll(idxRow.GetVal("activitePrincipaleEtablissement"), ".", "")
			if matched, matchErr := regexp.MatchString(`^[0-9]{4}[A-Z]$`, ape); matchErr == nil && matched {
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
		res.AddRegularError(err)
	}

	if val, ok := idxRow.GetOptionalVal("latitude"); ok && val != "" {
		sirene.Latitude, err = strconv.ParseFloat(val, 64)
		res.AddRegularError(err)
	}

	res.AddTuple(sirene)
	return nil
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
