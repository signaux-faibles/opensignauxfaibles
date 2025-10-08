package urssaf

import (
	"io"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"time"
)

type DelaiParser struct{}

func NewDelaiParser() engine.Parser {
	return &DelaiParser{}
}

func (parser *DelaiParser) Type() base.ParserType { return base.Delai }
func (parser *DelaiParser) New(r io.Reader) engine.ParserInst {
	return &UrssafParserInst{
		parsing.CsvParserInst{
			Reader:     r,
			RowParser:  &delaiRowParser{},
			Comma:      ';',
			LazyQuotes: false,
			DestTuple:  Delai{},
		},
	}
}

type delaiRowParser struct {
	UrssafRowParser
}

func (rp *delaiRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) error {

	idxRow := idx.IndexRow(row)

	date, err := time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))

	if err != nil {
		res.AddRegularError(err)
		return err

	}

	siret, err := rp.GetComptes().GetSiret(idxRow.GetVal("Numero_compte_externe"), &date)
	if err != nil {
		res.SetFilterError(err)
		return err
	}

	delai := Delai{}
	delai.Siret = siret
	delai.NumeroCompte = idxRow.GetVal("Numero_compte_externe")
	delai.NumeroContentieux = idxRow.GetVal("Numero_structure")
	delai.DateCreation, err = time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))
	res.AddRegularError(err)
	delai.DateEcheance, err = time.Parse("02/01/2006", idxRow.GetVal("Date_echeance"))
	res.AddRegularError(err)
	delai.DureeDelai, err = idxRow.GetInt("Duree_delai")
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
	return nil
}
