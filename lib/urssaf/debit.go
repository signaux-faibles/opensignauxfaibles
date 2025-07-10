package urssaf

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/misc"
)

// Debit Débit – fichier Urssaf
type Debit struct {
	key                          string       `                                                             csv:"siret"`
	NumeroCompte                 string       `input:"num_cpte"         json:"numero_compte"                csv:"numéro_compte"`
	NumeroEcartNegatif           string       `input:"Num_Ecn"          json:"numero_ecart_negatif"         csv:"numéro_écart_négatif"`
	DateTraitement               time.Time    `input:"Dt_trt_ecn"       json:"date_traitement"              csv:"date_traitement"`
	PartOuvriere                 float64      `input:"Mt_PO"            json:"part_ouvriere"                csv:"part_ouvrière"`
	PartPatronale                float64      `input:"Mt_PP"            json:"part_patronale"               csv:"part_patronale"`
	NumeroHistoriqueEcartNegatif *int         `input:"Num_Hist_Ecn"     json:"numero_historique"            csv:"numéro_historique_écart_négatif"`
	EtatCompte                   *int         `input:"Etat_cpte"        json:"etat_compte"                  csv:"état_compte"`
	CodeProcedureCollective      string       `input:"Cd_pro_col"       json:"code_procedure_collective"    csv:"code_procédure_collective"`
	Periode                      misc.Periode `input:"Periode"          json:"periode"                      csv:"période"`
	CodeOperationEcartNegatif    string       `input:"Cd_op_ecn"        json:"code_operation_ecart_negatif" csv:"code_opération_écart_négatif"`
	CodeMotifEcartNegatif        string       `input:"Motif_ecn"        json:"code_motif_ecart_negatif"     csv:"code_motif_écart_négatif"`
	Recours                      bool         `input:"Recours_en_cours" json:"recours_en_cours"             csv:"recours"`
}

func (debit Debit) Headers() []string {
	return marshal.ExtractCSVTags(debit)
}

func (debit Debit) Values() []string {
	return []string{
		debit.key,
		debit.NumeroCompte,
		debit.NumeroEcartNegatif,
		marshal.TimeToCSV(&debit.DateTraitement),
		marshal.FloatToCSV(&debit.PartOuvriere),
		marshal.FloatToCSV(&debit.PartPatronale),
		marshal.IntToCSV(debit.NumeroHistoriqueEcartNegatif),
		marshal.IntToCSV(debit.EtatCompte),
		debit.CodeProcedureCollective,
		debit.Periode.String(),
		debit.CodeOperationEcartNegatif,
		debit.CodeMotifEcartNegatif,
		marshal.BoolToCSV(&debit.Recours),
	}
}

// Key _id de l'objet
func (debit Debit) Key() string {
	return debit.key
}

// Scope de l'objet
func (debit Debit) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (debit Debit) Type() string {
	return "debit"
}

// ParserDebit fournit une instance utilisable par ParseFilesFromBatch.
var ParserDebit = &debitParser{}

type debitParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
	idx     marshal.ColMapping
}

func (parser *debitParser) Type() string {
	return "debit"
}

func (parser *debitParser) Close() error {
	return parser.file.Close()
}

func (parser *debitParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *debitParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Debit{})
	}
	return err
}

func (parser *debitParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parser.parseLine(row, parsedLine)
	})

}

func (parser *debitParser) parseLine(row []string, parsedLine *marshal.ParsedLineResult) {
	idxRow := parser.idx.IndexRow(row)
	period, _ := marshal.UrssafToPeriod(idxRow.GetVal("Periode"))
	date := period.Start
	if siret, err := marshal.GetSiretFromComptesMapping(idxRow.GetVal("num_cpte"), &date, parser.comptes); err == nil {
		parseDebitLine(siret, idxRow, parsedLine)
		if len(parsedLine.Errors) > 0 {
			parsedLine.Tuples = []marshal.Tuple{}
		}
	} else {
		parsedLine.SetFilterError(err)
	}
}

func parseDebitLine(siret string, idxRow marshal.IndexedRow, parsedLine *marshal.ParsedLineResult) {

	debit := Debit{
		key:                       siret,
		NumeroCompte:              idxRow.GetVal("num_cpte"),
		NumeroEcartNegatif:        idxRow.GetVal("Num_Ecn"),
		CodeProcedureCollective:   idxRow.GetVal("Cd_pro_col"),
		CodeOperationEcartNegatif: idxRow.GetVal("Cd_op_ecn"),
		CodeMotifEcartNegatif:     idxRow.GetVal("Motif_ecn"),
	}

	var err error
	debit.DateTraitement, err = marshal.UrssafToDate(idxRow.GetVal("Dt_trt_ecn"))
	parsedLine.AddRegularError(err)
	partOuvriere, err := idxRow.GetFloat64("Mt_PO")
	parsedLine.AddRegularError(err)
	debit.PartOuvriere = *partOuvriere / 100
	partPatronale, err := idxRow.GetFloat64("Mt_PP")
	parsedLine.AddRegularError(err)
	debit.PartPatronale = *partPatronale / 100
	debit.NumeroHistoriqueEcartNegatif, err = idxRow.GetInt("Num_Hist_Ecn")
	parsedLine.AddRegularError(err)
	debit.EtatCompte, err = idxRow.GetInt("Etat_cpte")
	parsedLine.AddRegularError(err)
	debit.Periode, err = marshal.UrssafToPeriod(idxRow.GetVal("Periode"))
	parsedLine.AddRegularError(err)
	debit.Recours, err = idxRow.GetBool("Recours_en_cours")
	parsedLine.AddRegularError(err)
	// debit.MontantMajorations, err = strconv.ParseFloat(idxRow.GetVal("montantMajorations"), 64)
	// tracker.Error(err)
	// debit.MontantMajorations = debit.MontantMajorations / 100
	parsedLine.AddTuple(debit)
}
