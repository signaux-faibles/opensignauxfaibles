package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
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

// ParserDebit expose le parseur et le type de fichier qu'il supporte.
var ParserDebit = marshal.Parser{FileType: "debit", FileParser: ParseDebitFile}

// ParseDebitFile extrait les tuples depuis un fichier "débit" de l'URSSAF.
func ParseDebitFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker) marshal.ParsedLineChan {
	comptes, err := marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	if err != nil {
		tracker.Add(err)
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return nil
	}
	// defer file.Close() // TODO: à réactiver

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	// ligne de titre
	fields, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return nil
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
		tracker.Add(errors.New("CSV non conforme"))
		return nil
	}

	var lineNumber = 0 // starting with the header
	stopProgressLogger := marshal.LogProgress(&lineNumber)
	defer stopProgressLogger()

	parsedLineChan := make(marshal.ParsedLineChan)
	go func() {
		for {
			tuples := []marshal.Tuple{}
			lineNumber++
			row, err := reader.Read()
			if err == io.EOF {
				close(parsedLineChan)
				break
			} else if err != nil {
				tracker.Add(err)
			} else {
				period, _ := marshal.UrssafToPeriod(row[idx["periode"]])
				date := period.Start

				if siret, err := marshal.GetSiretFromComptesMapping(row[idx["numeroCompte"]], &date, comptes); err == nil {
					debit := parseDebitLine(siret, row, tracker, idx)
					if !tracker.HasErrorInCurrentCycle() {
						tuples = []marshal.Tuple{debit}
					}
				} else {
					tracker.Add(base.NewFilterError(err))
				}
			}
			parsedLineChan <- marshal.ParsedLineResult{Tuples: tuples, Errors: []marshal.ParseError{}}
		}
	}()
	return parsedLineChan
}

func parseDebitLine(siret string, row []string, tracker *gournal.Tracker, idx colMapping) Debit {

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
	tracker.Add(err)
	debit.PartOuvriere, err = strconv.ParseFloat(row[idx["partOuvriere"]], 64)
	tracker.Add(err)
	debit.PartOuvriere = debit.PartOuvriere / 100
	debit.PartPatronale, err = strconv.ParseFloat(row[idx["partPatronale"]], 64)
	tracker.Add(err)
	debit.PartPatronale = debit.PartPatronale / 100
	debit.NumeroHistoriqueEcartNegatif, err = strconv.Atoi(row[idx["numeroHistoriqueEcartNegatif"]])
	tracker.Add(err)
	debit.EtatCompte, err = strconv.Atoi(row[idx["etatCompte"]])
	tracker.Add(err)
	debit.Periode, err = marshal.UrssafToPeriod(row[idx["periode"]])
	tracker.Add(err)
	debit.Recours, err = strconv.ParseBool(row[idx["recours"]])
	tracker.Add(err)
	// debit.MontantMajorations, err = strconv.ParseFloat(row[idx["montantMajorations"]], 64)
	// tracker.Error(err)
	// debit.MontantMajorations = debit.MontantMajorations / 100
	return debit
}
