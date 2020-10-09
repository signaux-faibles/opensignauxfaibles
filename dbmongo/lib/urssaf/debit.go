package urssaf

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
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

// ParserDebit retourne les entrées lues depuis un fichier "débit" de l'URSSAF.
func ParserDebit(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)

	event := marshal.Event{
		Code:    "debitParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["debit"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key, "MaxParsingErrors": strconv.Itoa(engine.MaxParsingErrors)},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Add(err)
				event.Critical(tracker.Report("fatalError"))
			} else {
				event.Info(path + ": ouverture")
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			// ligne de titre
			fields, err := reader.Read()
			if err != nil {
				tracker.Add(err)
				event.Critical(tracker.Report("fatalError"))
			}

			dateTraitementIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Dt_trt_ecn" })
			partOuvriereIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PO" })
			partPatronaleIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PP" })
			numeroHistoriqueEcartNegatifIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Hist_Ecn" })
			periodeIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Periode" })
			etatCompteIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Etat_cpte" })
			numeroCompteIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "num_cpte" })
			numeroEcartNegatifIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Ecn" })
			codeProcedureCollectiveIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_pro_col" })
			codeOperationEcartNegatifIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_op_ecn" })
			codeMotifEcartNegatifIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Motif_ecn" })
			recoursIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Recours_en_cours" })
			// montantMajorationsIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Montant majorations de retard en centimes" })
			if misc.SliceMin(dateTraitementIndex, partOuvriereIndex, partPatronaleIndex, numeroHistoriqueEcartNegatifIndex, periodeIndex, etatCompteIndex, numeroCompteIndex, numeroEcartNegatifIndex, codeProcedureCollectiveIndex, codeOperationEcartNegatifIndex, codeMotifEcartNegatifIndex) < 0 {
				event.Critical(path + ": CSV non conforme")
				continue
			}

			var shouldBreak = false
			var lineNumber = 0 // starting with the header
			var maxParsingErrors = engine.MaxParsingErrors

			ctx, stopTheTimer := context.WithCancel(context.Background())
			go func(ctx context.Context) {
				for range time.Tick(time.Second * 2) {
					select {
					case <-ctx.Done():
						return
					default:
					}
					shouldBreak = maxParsingErrors > 0 && engine.ShouldBreak(tracker, maxParsingErrors)
					if shouldBreak {
						fmt.Printf("Reached %d parsing errors => stopping.\n", maxParsingErrors)
					} else {
						fmt.Printf("Reading csv line %d\n", lineNumber)
					}
				}
			}(ctx)
			defer stopTheTimer()

			for {
				lineNumber++
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Add(err)
					_, ok := err.(*csv.ParseError)
					if !ok {
						// we tolerate CSV parsing errors, but we generate a fatalError report for others, in order to interrupt the whole import process
						event.Critical(tracker.Report("fatalError"))
					}
					tracker.Next()
					continue
				}

				justParse := true
				if !justParse {
					period, _ := marshal.UrssafToPeriod(row[periodeIndex])
					date := period.Start

					if siret, err := marshal.GetSiret(row[numeroCompteIndex], &date, cache, batch); err == nil {

						debit := Debit{
							key:                       siret,
							NumeroCompte:              row[numeroCompteIndex],
							NumeroEcartNegatif:        row[numeroEcartNegatifIndex],
							CodeProcedureCollective:   row[codeProcedureCollectiveIndex],
							CodeOperationEcartNegatif: row[codeOperationEcartNegatifIndex],
							CodeMotifEcartNegatif:     row[codeMotifEcartNegatifIndex],
						}

						debit.DateTraitement, err = marshal.UrssafToDate(row[dateTraitementIndex])
						tracker.Add(err)
						debit.PartOuvriere, err = strconv.ParseFloat(row[partOuvriereIndex], 64)
						tracker.Add(err)
						debit.PartOuvriere = debit.PartOuvriere / 100
						debit.PartPatronale, err = strconv.ParseFloat(row[partPatronaleIndex], 64)
						tracker.Add(err)
						debit.PartPatronale = debit.PartPatronale / 100
						debit.NumeroHistoriqueEcartNegatif, err = strconv.Atoi(row[numeroHistoriqueEcartNegatifIndex])
						tracker.Add(err)
						debit.EtatCompte, err = strconv.Atoi(row[etatCompteIndex])
						tracker.Add(err)
						debit.Periode, err = marshal.UrssafToPeriod(row[periodeIndex])
						tracker.Add(err)
						debit.Recours, err = strconv.ParseBool(row[recoursIndex])
						tracker.Add(err)
						// debit.MontantMajorations, err = strconv.ParseFloat(row[montantMajorationsIndex], 64)
						// tracker.Error(err)
						// debit.MontantMajorations = debit.MontantMajorations / 100

						if !tracker.HasErrorInCurrentCycle() {
							outputChannel <- debit
						}
					} else {
						tracker.Add(base.NewFilterError(err))
						continue
					}
				}

				if shouldBreak {
					break
				}
				tracker.Next()
			}

			event.Debug(tracker.Report("abstract"))
			file.Close()
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}
