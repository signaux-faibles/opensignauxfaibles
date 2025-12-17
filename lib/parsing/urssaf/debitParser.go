package urssaf

import (
	"fmt"
	"io"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
)

type DebitParser struct{}

func NewDebitParser() engine.Parser {
	return &DebitParser{}
}

func (parser *DebitParser) Type() engine.ParserType { return engine.Debit }
func (parser *DebitParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:        r,
		RowParser:     &debitRowParser{},
		Comma:         ';',
		CaseSensitive: false,
		LazyQuotes:    false,
		DestTuple:     Debit{},
	}
}

type debitRowParser struct{}

func (rp *debitRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {

	idxRow := idx.IndexRow(row)

	periodeDebut, periodeFin, err := UrssafToPeriod(idxRow.GetVal("Periode"))
	res.AddRegularError(err)

	siret := idxRow.GetVal("Siret")
	if sfregexp.RegexpDict["acossInternal"].MatchString(siret) {
		res.SetFilterError(fmt.Errorf("acoss internal id: %s", siret))
		return
	}

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

	// Calcul de la période de prise en compte :
	// - Si date_traitement <= 20 du mois : période en cours (1er du mois)
	// - Si date_traitement > 20 du mois : période suivante (1er du mois suivant)
	if debit.DateTraitement.Day() <= 20 {
		debit.PeriodePriseEnCompte = time.Date(
			debit.DateTraitement.Year(),
			debit.DateTraitement.Month(),
			1, 0, 0, 0, 0,
			debit.DateTraitement.Location(),
		)
	} else {
		debit.PeriodePriseEnCompte = time.Date(
			debit.DateTraitement.Year(),
			debit.DateTraitement.Month(),
			1, 0, 0, 0, 0,
			debit.DateTraitement.Location(),
		).AddDate(0, 1, 0)
	}

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
}
