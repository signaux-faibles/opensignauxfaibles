package urssaf

import (
	"fmt"
	"io"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
	"strings"
	"time"
)

type DelaiParser struct{}

func NewDelaiParser() engine.Parser {
	return &DelaiParser{}
}

func (parser *DelaiParser) Type() engine.ParserType { return engine.Delai }
func (parser *DelaiParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:        r,
		RowParser:     &delaiRowParser{},
		Comma:         ';',
		CaseSensitive: false,
		LazyQuotes:    false,
		DestTuple:     Delai{},
	}
}

type delaiRowParser struct{}

func (rp *delaiRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {

	idxRow := idx.IndexRow(row)

	var err error

	delai := Delai{}
	delai.Siret = idxRow.GetVal("Siret")

	if strings.TrimSpace(delai.Siret) == "" ||
		delai.Siret == "@" ||
		sfregexp.RegexpDict["acossInternal"].MatchString(delai.Siret) ||
		sfregexp.RegexpDict["delaiInvalid"].MatchString(delai.Siret) {
		res.SetFilterError(fmt.Errorf("acoss internal id: %s", delai.Siret))
		return
	}

	delai.NumeroCompte = idxRow.GetVal("Numero_compte_externe")
	delai.NumeroContentieux = idxRow.GetVal("Numero_structure")
	delai.DateCreation, err = time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))
	res.AddRegularError(err)
	delai.DateEcheance, err = time.Parse("02/01/2006", idxRow.GetVal("Date_echeance"))
	res.AddRegularError(err)
	delai.DureeDelai, _ = idxRow.GetInt("Duree_delai")
	delai.Denomination = idxRow.GetVal("Denomination_premiere_ligne")
	delai.Indic6m = idxRow.GetVal("Indic_6M")
	delai.AnneeCreation, err = idxRow.GetInt("Annee_creation")
	res.AddRegularError(err)
	delai.MontantEcheancier, err = idxRow.GetCommaFloat64("Montant_global_echeancier")
	res.AddRegularError(err)
	delai.Stade = idxRow.GetVal("Code_externe_stade")
	delai.Action = idxRow.GetVal("Code_externe_action")

	if len(res.Errors) == 0 {
		res.AddTuple(delai)
	}
}
