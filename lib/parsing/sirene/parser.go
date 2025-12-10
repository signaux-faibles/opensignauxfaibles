// Package sirene fournit les parser pour extraire les données des fichiers
// Sirene
package sirene

import (
	//"bufio"

	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
)

type SireneParser struct{}

func NewSireneParser() engine.Parser {
	return &SireneParser{}
}

func (p *SireneParser) Type() engine.ParserType { return engine.Sirene }
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

func (rp *sireneRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {
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
	sirene.CodeCedex = idxRow.GetVal("codeCedexEtablissement")
	sirene.Cedex = idxRow.GetVal("libelleCedexEtablissement")
	sirene.CodePaysEtranger = idxRow.GetVal("codePaysEtrangerEtablissement")
	sirene.PaysEtranger = idxRow.GetVal("libellePaysEtrangerEtablissement")

	codePostal := idxRow.GetVal("codePostalEtablissement")
	codeCommune := idxRow.GetVal("codeCommuneEtablissement")

	if sfregexp.ValidCodePostal(codePostal) {
		sirene.CodePostal = idxRow.GetVal("codePostalEtablissement")
	}

	if sfregexp.ValidCodeCommune(codeCommune) {
		sirene.CodeCommune = codeCommune
		// on extrait le département pour codeCommune qui est mieux renseigné que
		// codePostal
		sirene.Departement = extractDepartement(codeCommune)
	}

	if idxRow.GetVal("activitePrincipaleEtablissement") != "" {
		nomenclature := idxRow.GetVal("nomenclatureActivitePrincipaleEtablissement")
		if nomenclature != "NAFRev2" {
			res.SetFilterError(fmt.Errorf("nomenclature activité non NAFRev2 : %s", nomenclature))
			return
		}
		ape := idxRow.GetVal("activitePrincipaleEtablissement")
		if matched, matchErr := regexp.MatchString(`^[0-9]{2}\.[0-9]{2}[A-Z]$`, ape); matchErr == nil && matched {
			sirene.APE = ape
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

	// etatAdministratifEtablissement: "A" pour actif, "F" pour fermé
	etatAdministratif := idxRow.GetVal("etatAdministratifEtablissement")
	sirene.EstActif = (etatAdministratif == "A")

	res.AddTuple(sirene)
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

// extractDemartement extrait le département d'un code commune INSEE.
// Il est attendu sans vérification que le code postal est composé de 5
// chiffres, ou qu'il commence par 2A ou 2B (corse)
func extractDepartement(codeCommune string) string {
	var departement string
	// Départements et territoires d'outre-mer (codes à 3 chiffres)
	if codeCommune[0:2] == "97" || codeCommune[0:2] == "98" {
		departement = codeCommune[0:3]
	} else {
		// Départements de métropole (codes à 2 chiffres)
		// Les codes INSEE pour la corse commencent bien par 2A ou 2B
		departement = codeCommune[0:2]
	}
	return departement

}
