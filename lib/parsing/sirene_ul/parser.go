// Package sireneul exporte un parser pour extraire les données d'un fichier
// Sirene Unités légales
package sireneul

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type SireneULParser struct{}

func NewSireneULParser() engine.Parser {
	return &SireneULParser{}
}

func (p *SireneULParser) Type() engine.ParserType { return engine.SireneUl }
func (p *SireneULParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:     r,
		RowParser:  &sireneULRowParser{},
		Comma:      ',',
		LazyQuotes: false,
		DestTuple:  SireneUL{},
	}
}

type sireneULRowParser struct{}

func (rp *sireneULRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {
	var err error

	idxRow := idx.IndexRow(row)
	sireneul := SireneUL{}
	sireneul.Siren = idxRow.GetVal("siren")
	sireneul.RaisonSociale = idxRow.GetVal("denominationUniteLegale")
	sireneul.Prenom1UniteLegale = idxRow.GetVal("prenom1UniteLegale")
	sireneul.Prenom2UniteLegale = idxRow.GetVal("prenom2UniteLegale")
	sireneul.Prenom3UniteLegale = idxRow.GetVal("prenom3UniteLegale")
	sireneul.Prenom4UniteLegale = idxRow.GetVal("prenom4UniteLegale")
	sireneul.NomUniteLegale = idxRow.GetVal("nomUniteLegale")
	sireneul.NomUsageUniteLegale = idxRow.GetVal("nomUsageUniteLegale")
	sireneul.CategorieJuridique = idxRow.GetVal("categorieJuridiqueUniteLegale")

	if idxRow.GetVal("activitePrincipaleUniteLegale") != "" {
		nomenclature := idxRow.GetVal("nomenclatureActivitePrincipaleUniteLegale")
		if nomenclature != "NAFRev2" {
			res.SetFilterError(fmt.Errorf("nomenclature activité non NAFRev2 : %s", nomenclature))
			return
		}
		ape := idxRow.GetVal("activitePrincipaleUniteLegale")
		if matched, matchErr := regexp.MatchString(`^[0-9]{2}\.[0-9]{2}[A-Z]$`, ape); matchErr == nil && matched {
			sireneul.APE = ape
		}
	}

	creation, err := time.Parse("2006-01-02", idxRow.GetVal("dateCreationUniteLegale")) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sireneul.Creation = &creation
	}

	// etatAdministratifUniteLegale: "A" pour actif, "F" pour fermé
	etatAdministratif := idxRow.GetVal("etatAdministratifUniteLegale")
	sireneul.EstActif = (etatAdministratif == "A")

	res.AddTuple(sireneul)
}
