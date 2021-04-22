package urssaf

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// Debit Débit – fichier Urssaf
type Debit struct {
	key                          string       `                       hash:"-"`
	NumeroCompte                 string       `col:"num_cpte"         json:"numero_compte"                bson:"numero_compte"`
	NumeroEcartNegatif           string       `col:"Num_Ecn"          json:"numero_ecart_negatif"         bson:"numero_ecart_negatif"`
	DateTraitement               time.Time    `col:"Dt_trt_ecn"       json:"date_traitement"              bson:"date_traitement"`
	PartOuvriere                 float64      `col:"Mt_PO"            json:"part_ouvriere"                bson:"part_ouvriere"`
	PartPatronale                float64      `col:"Mt_PP"            json:"part_patronale"               bson:"part_patronale"`
	NumeroHistoriqueEcartNegatif *int         `col:"Num_Hist_Ecn"     json:"numero_historique"            bson:"numero_historique"`
	EtatCompte                   *int         `col:"Etat_cpte"        json:"etat_compte"                  bson:"etat_compte"`
	CodeProcedureCollective      string       `col:"Cd_pro_col"       json:"code_procedure_collective"    bson:"code_procedure_collective"`
	Periode                      misc.Periode `col:"Periode"          json:"periode"                      bson:"periode"`
	CodeOperationEcartNegatif    string       `col:"Cd_op_ecn"        json:"code_operation_ecart_negatif" bson:"code_operation_ecart_negatif"`
	CodeMotifEcartNegatif        string       `col:"Motif_ecn"        json:"code_motif_ecart_negatif"     bson:"code_motif_ecart_negatif"`
	Recours                      bool         `col:"Recours_en_cours" json:"recours_en_cours"             bson:"recours_en_cours"`
	DebitSuivant                 string       `                       json:"debit_suivant,omitempty"      bson:"debit_suivant,omitempty"` // généré par traitement map-reduce
	// MontantMajorations        float64      `                       json:"montant_majorations"          bson:"montant_majorations"`  // TODO: montant_majorations n'est pas fourni par les fichiers debit de l'urssaf pour l'instant, mais on aimerait y avoir accès un jour.
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

func (parser *debitParser) GetFileType() string {
	return "debit"
}

func (parser *debitParser) Close() error {
	return parser.file.Close()
}

func (parser *debitParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *debitParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ';', false)
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
