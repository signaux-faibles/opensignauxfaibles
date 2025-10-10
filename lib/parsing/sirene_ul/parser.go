// Package sireneul exporte un parser pour extraire les données d'un fichier
// Sirene Unités légales
package sireneul

import (
	"io"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type SireneULParser struct{}

func NewSireneULParser() engine.Parser {
	return &SireneULParser{}
}

func (p *SireneULParser) Type() base.ParserType { return base.SireneUl }
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

func (rp *sireneULRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) error {
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
	sireneul.CodeStatutJuridique = idxRow.GetVal("categorieJuridiqueUniteLegale")

	creation, err := time.Parse("2006-01-02", idxRow.GetVal("dateCreationUniteLegale")) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sireneul.Creation = &creation
	}

	res.AddTuple(sireneul)
	return nil
}
