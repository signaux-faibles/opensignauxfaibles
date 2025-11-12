package urssaf

import (
	"io"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type DebitParser struct{}

func NewDebitParser() engine.Parser {
	return &DebitParser{}
}

func (parser *DebitParser) Type() engine.ParserType { return engine.Debit }
func (parser *DebitParser) New(r io.Reader) engine.ParserInst {
	return &UrssafParserInst{
		parsing.CsvParserInst{
			Reader:        r,
			RowParser:     &debitRowParser{},
			Comma:         ';',
			CaseSensitive: false,
			LazyQuotes:    false,
			DestTuple:     Debit{},
		},
	}
}

type debitRowParser struct {
	UrssafRowParser
}

func (rp *debitRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {

	idxRow := idx.IndexRow(row)

	periodeDebut, periodeFin, err := UrssafToPeriod(idxRow.GetVal("Periode"))
	res.AddRegularError(err)

	if siret, err := rp.GetComptes().GetSiret(idxRow.GetVal("num_cpte"), &periodeDebut); err == nil {
		debit := Debit{
			Siret:                     siret,
			NumeroCompte:              idxRow.GetVal("num_cpte"),
			NumeroEcartNegatif:        idxRow.GetVal("Num_Ecn"),
			CodeProcedureCollective:   idxRow.GetVal("Cd_pro_col"),
			CodeOperationEcartNegatif: idxRow.GetVal("Cd_op_ecn"),
			CodeMotifEcartNegatif:     idxRow.GetVal("Motif_ecn"),
		}

		debit.DateTraitement, err = UrssafToDate(idxRow.GetVal("Dt_trt_ecn"))
		res.AddRegularError(err)
		partOuvriere, err := idxRow.GetFloat64("Mt_PO")
		res.AddRegularError(err)
		debit.PartOuvriere = *partOuvriere / 100
		partPatronale, err := idxRow.GetFloat64("Mt_PP")
		res.AddRegularError(err)
		debit.PartPatronale = *partPatronale / 100
		debit.NumeroHistoriqueEcartNegatif, err = idxRow.GetInt("Num_Hist_Ecn")
		res.AddRegularError(err)
		debit.EtatCompte, err = idxRow.GetInt("Etat_cpte")
		res.AddRegularError(err)

		debit.PeriodeDebut = periodeDebut
		debit.PeriodeFin = periodeFin

		debit.Recours, err = idxRow.GetBool("Recours_en_cours")
		res.AddRegularError(err)
		// debit.MontantMajorations, err = strconv.ParseFloat(idxRow.GetVal("montantMajorations"), 64)
		// tracker.Error(err)
		// debit.MontantMajorations = debit.MontantMajorations / 100
		res.AddTuple(debit)

		if len(res.Errors) > 0 {
			res.Tuples = []engine.Tuple{}
		}
	} else {
		res.SetFilterError(err)
	}
}
