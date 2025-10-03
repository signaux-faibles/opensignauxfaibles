package urssaf

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

// Debit Débit – fichier Urssaf
type Debit struct {
	Siret                        string    `                         json:"-"                            sql:"siret"                            csv:"siret"`
	NumeroCompte                 string    `input:"num_cpte"         json:"numero_compte"                sql:"numero_compte"                    csv:"numéro_compte"`
	NumeroEcartNegatif           string    `input:"Num_Ecn"          json:"numero_ecart_negatif"         sql:"numero_ecart_negatif"             csv:"numéro_écart_négatif"`
	DateTraitement               time.Time `input:"Dt_trt_ecn"       json:"date_traitement"              sql:"date_traitement"                  csv:"date_traitement"`
	PartOuvriere                 float64   `input:"Mt_PO"            json:"part_ouvriere"                sql:"part_ouvriere"                    csv:"part_ouvrière"`
	PartPatronale                float64   `input:"Mt_PP"            json:"part_patronale"               sql:"part_patronale"                   csv:"part_patronale"`
	NumeroHistoriqueEcartNegatif *int      `input:"Num_Hist_Ecn"     json:"numero_historique"            sql:"numero_historique_ecart_negatif"  csv:"numéro_historique_écart_négatif"`
	EtatCompte                   *int      `input:"Etat_cpte"        json:"etat_compte"                  sql:"etat_compte"                      csv:"état_compte"`
	CodeProcedureCollective      string    `input:"Cd_pro_col"       json:"code_procedure_collective"    sql:"code_procedure_collective"        csv:"code_procédure_collective"`
	PeriodeDebut                 time.Time `input:"Periode"          json:"periode_debut"                sql:"periode_debut"                    csv:"période_début"`
	PeriodeFin                   time.Time `input:"Periode"          json:"periode_fin"                  sql:"periode_fin"                      csv:"période_fin"`
	CodeOperationEcartNegatif    string    `input:"Cd_op_ecn"        json:"code_operation_ecart_negatif" sql:"code_operation_ecart_negatif"     csv:"code_opération_écart_négatif"`
	CodeMotifEcartNegatif        string    `input:"Motif_ecn"        json:"code_motif_ecart_negatif"     sql:"code_motif_ecart_negatif"         csv:"code_motif_écart_négatif"`
	Recours                      bool      `input:"Recours_en_cours" json:"recours_en_cours"             sql:"recours_en_cours"                 csv:"recours"`
}

// Key _id de l'objet
func (debit Debit) Key() string {
	return debit.Siret
}

// Scope de l'objet
func (debit Debit) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (debit Debit) Type() base.ParserType {
	return base.Debit
}

// ParserDebit fournit une instance utilisable par ParseFilesFromBatch.
var ParserDebit = &debitParser{}

type debitParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes engine.Comptes
	idx     engine.ColMapping
}

func (parser *debitParser) Type() base.ParserType {
	return base.Debit
}

func (parser *debitParser) Close() error {
	return parser.file.Close()
}

func (parser *debitParser) Init(cache *engine.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = engine.GetCompteSiretMapping(*cache, batch, engine.OpenAndReadSiretMapping)
	return err
}

func (parser *debitParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, Debit{})
	}
	return err
}

func (parser *debitParser) ParseLines(parsedLineChan chan engine.ParsedLineResult) {
	engine.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *engine.ParsedLineResult) {
		parser.parseLine(row, parsedLine)
	})

}

func (parser *debitParser) parseLine(row []string, parsedLine *engine.ParsedLineResult) {
	idxRow := parser.idx.IndexRow(row)
	periodeDebut, periodeFin, err := engine.UrssafToPeriod(idxRow.GetVal("Periode"))
	parsedLine.AddRegularError(err)

	if siret, err := engine.GetSiretFromComptesMapping(idxRow.GetVal("num_cpte"), &periodeDebut, parser.comptes); err == nil {
		debit := Debit{
			Siret:                     siret,
			NumeroCompte:              idxRow.GetVal("num_cpte"),
			NumeroEcartNegatif:        idxRow.GetVal("Num_Ecn"),
			CodeProcedureCollective:   idxRow.GetVal("Cd_pro_col"),
			CodeOperationEcartNegatif: idxRow.GetVal("Cd_op_ecn"),
			CodeMotifEcartNegatif:     idxRow.GetVal("Motif_ecn"),
		}

		var err error
		debit.DateTraitement, err = engine.UrssafToDate(idxRow.GetVal("Dt_trt_ecn"))
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

		debit.PeriodeDebut = periodeDebut
		debit.PeriodeFin = periodeFin

		debit.Recours, err = idxRow.GetBool("Recours_en_cours")
		parsedLine.AddRegularError(err)
		// debit.MontantMajorations, err = strconv.ParseFloat(idxRow.GetVal("montantMajorations"), 64)
		// tracker.Error(err)
		// debit.MontantMajorations = debit.MontantMajorations / 100
		parsedLine.AddTuple(debit)

		if len(parsedLine.Errors) > 0 {
			parsedLine.Tuples = []engine.Tuple{}
		}
	} else {
		parsedLine.SetFilterError(err)
	}
}
