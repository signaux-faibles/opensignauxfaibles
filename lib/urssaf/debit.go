package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// Debit Débit – fichier Urssaf
type Debit struct {
	key                          string       `hash:"-"`
	NumeroCompte                 string       `json:"numero_compte" bson:"numero_compte"`
	NumeroEcartNegatif           string       `json:"numero_ecart_negatif" bson:"numero_ecart_negatif"`
	DateTraitement               time.Time    `json:"date_traitement" bson:"date_traitement"`
	PartOuvriere                 float64      `json:"part_ouvriere" bson:"part_ouvriere"`
	PartPatronale                float64      `json:"part_patronale" bson:"part_patronale"`
	NumeroHistoriqueEcartNegatif int          `json:"numero_historique" bson:"numero_historique"`
	EtatCompte                   int          `json:"etat_compte" bson:"etat_compte"`
	CodeProcedureCollective      string       `json:"code_procedure_collective" bson:"code_procedure_collective"`
	Periode                      misc.Periode `json:"periode" bson:"periode"`
	CodeOperationEcartNegatif    string       `json:"code_operation_ecart_negatif" bson:"code_operation_ecart_negatif"`
	CodeMotifEcartNegatif        string       `json:"code_motif_ecart_negatif" bson:"code_motif_ecart_negatif"`
	DebitSuivant                 string       `json:"debit_suivant,omitempty" bson:"debit_suivant,omitempty"`
	Recours                      bool         `json:"recours_en_cours" bson:"recours_en_cours"`
	// MontantMajorations           float64   `json:"montant_majorations" bson:"montant_majorations"`
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

type colMapping map[string]int

// ParserDebit fournit une instance utilisable par ParseFilesFromBatch.
var ParserDebit = &debitParser{}

type debitParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
	idx     colMapping
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
	parser.file, parser.reader, err = openDebitFile(filePath)
	if err == nil {
		parser.idx, err = parseDebitColMapping(parser.reader)
	}
	return err
}

func openDebitFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	return file, reader, err
}

func parseDebitColMapping(reader *csv.Reader) (colMapping, error) {
	fields, err := reader.Read()
	if err != nil {
		return nil, err
	}
	var idx = colMapping{
		"dateTraitement":               misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Dt_trt_ecn" }),
		"partOuvriere":                 misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PO" }),
		"partPatronale":                misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PP" }),
		"numeroHistoriqueEcartNegatif": misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Hist_Ecn" }),
		"periode":                      misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Periode" }),
		"etatCompte":                   misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Etat_cpte" }),
		"numeroCompte":                 misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "num_cpte" }),
		"numeroEcartNegatif":           misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Ecn" }),
		"codeProcedureCollective":      misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_pro_col" }),
		"codeOperationEcartNegatif":    misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_op_ecn" }),
		"codeMotifEcartNegatif":        misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Motif_ecn" }),
		"recours":                      misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Recours_en_cours" }),
	}
	// montantMajorationsIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Montant majorations de retard en centimes" })
	if misc.SliceMin(idx["dateTraitement"], idx["partOuvriere"], idx["partPatronale"], idx["numeroHistoriqueEcartNegatif"], idx["periode"], idx["etatCompte"], idx["numeroCompte"], idx["numeroEcartNegatif"], idx["codeProcedureCollective"], idx["codeOperationEcartNegatif"], idx["codeMotifEcartNegatif"]) < 0 {
		return nil, errors.New("CSV non conforme")
	}
	return idx, nil
}

func (parser *debitParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	var lineNumber = 0                                     // starting with the header
	stopProgressLogger := marshal.LogProgress(&lineNumber) // TODO: move this call to runParserWithSirenFilter()
	defer stopProgressLogger()

	for {
		parsedLine := marshal.ParsedLineResult{}
		lineNumber++
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			period, _ := marshal.UrssafToPeriod(row[parser.idx["periode"]])
			date := period.Start

			if siret, err := marshal.GetSiretFromComptesMapping(row[parser.idx["numeroCompte"]], &date, parser.comptes); err == nil {
				parseDebitLine(siret, row, parser.idx, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []marshal.Tuple{}
				}
			} else {
				parsedLine.SetFilterError(err)
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseDebitLine(siret string, row []string, idx colMapping, parsedLine *marshal.ParsedLineResult) {

	debit := Debit{
		key:                       siret,
		NumeroCompte:              row[idx["numeroCompte"]],
		NumeroEcartNegatif:        row[idx["numeroEcartNegatif"]],
		CodeProcedureCollective:   row[idx["codeProcedureCollective"]],
		CodeOperationEcartNegatif: row[idx["codeOperationEcartNegatif"]],
		CodeMotifEcartNegatif:     row[idx["codeMotifEcartNegatif"]],
	}

	var err error
	debit.DateTraitement, err = marshal.UrssafToDate(row[idx["dateTraitement"]])
	parsedLine.AddRegularError(err)
	debit.PartOuvriere, err = strconv.ParseFloat(row[idx["partOuvriere"]], 64)
	parsedLine.AddRegularError(err)
	debit.PartOuvriere = debit.PartOuvriere / 100
	debit.PartPatronale, err = strconv.ParseFloat(row[idx["partPatronale"]], 64)
	parsedLine.AddRegularError(err)
	debit.PartPatronale = debit.PartPatronale / 100
	debit.NumeroHistoriqueEcartNegatif, err = strconv.Atoi(row[idx["numeroHistoriqueEcartNegatif"]])
	parsedLine.AddRegularError(err)
	debit.EtatCompte, err = strconv.Atoi(row[idx["etatCompte"]])
	parsedLine.AddRegularError(err)
	debit.Periode, err = marshal.UrssafToPeriod(row[idx["periode"]])
	parsedLine.AddRegularError(err)
	debit.Recours, err = strconv.ParseBool(row[idx["recours"]])
	parsedLine.AddRegularError(err)
	// debit.MontantMajorations, err = strconv.ParseFloat(row[idx["montantMajorations"]], 64)
	// tracker.Error(err)
	// debit.MontantMajorations = debit.MontantMajorations / 100
	parsedLine.AddTuple(debit)
}
