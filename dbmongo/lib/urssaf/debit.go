package urssaf

import (
  "bufio"
  "dbmongo/lib/engine"
  "dbmongo/lib/misc"
  "encoding/csv"
  "io"
  "os"
  "strconv"
  "time"

  "github.com/chrnin/gournal"
  "github.com/spf13/viper"
)

// Debit Débit – fichier Urssaf
type Debit struct {
  key                          string    `hash:"-"`
  NumeroCompte                 string    `json:"numero_compte" bson:"numero_compte"`
  NumeroEcartNegatif           string    `json:"numero_ecart_negatif" bson:"numero_ecart_negatif"`
  DateTraitement               time.Time `json:"date_traitement" bson:"date_traitement"`
  PartOuvriere                 float64   `json:"part_ouvriere" bson:"part_ouvriere"`
  PartPatronale                float64   `json:"part_patronale" bson:"part_patronale"`
  NumeroHistoriqueEcartNegatif int       `json:"numero_historique" bson:"numero_historique"`
  EtatCompte                   int       `json:"etat_compte" bson:"etat_compte"`
  CodeProcedureCollective      string    `json:"code_procedure_collective" bson:"code_procedure_collective"`
  Periode                      Periode   `json:"periode" bson:"periode"`
  CodeOperationEcartNegatif    string    `json:"code_operation_ecart_negatif" bson:"code_operation_ecart_negatif"`
  CodeMotifEcartNegatif        string    `json:"code_motif_ecart_negatif" bson:"code_motif_ecart_negatif"`
  DebitSuivant                 string    `json:"debit_suivant,omitempty" bson:"debit_suivant,omitempty"`
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

func parseDebit(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event) {
  outputChannel := make(chan engine.Tuple)
  eventChannel := make(chan engine.Event)

  event := engine.Event{
    Code:    "debitParser",
    Channel: eventChannel,
  }

  go func() {
    for _, path := range batch.Files["debit"] {
      tracker := gournal.NewTracker(
        map[string]string{"path": path},
        engine.TrackerReports)

      file, err := os.Open(viper.GetString("APP_DATA") + path)
      if err != nil {
        tracker.Error(err)
        event.Critical(tracker.Report("fatalError"))
      }

      reader := csv.NewReader(bufio.NewReader(file))
      reader.Comma = ';'
      // ligne de titre
      fields, err := reader.Read()
      if err != nil {
        tracker.Error(err)
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
      // montantMajorationsIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "Montant majorations de retard en centimes" })
      if misc.SliceMin(dateTraitementIndex, partOuvriereIndex, partPatronaleIndex, numeroHistoriqueEcartNegatifIndex, periodeIndex, etatCompteIndex, numeroCompteIndex, numeroEcartNegatifIndex, codeProcedureCollectiveIndex, codeOperationEcartNegatifIndex, codeMotifEcartNegatifIndex) < 0 {
        event.Critical(path + ": CSV non conforme")
        continue
      }

      for {
        row, err := reader.Read()
        if err == io.EOF {
          break
        } else if err != nil {
          tracker.Error(err)
          event.Critical(tracker.Report("fatalError"))
          break
        }

        date, err := urssafToDate(row[periodeIndex])
        if err != nil { date = time.Now() }
        if siret, err := mapping.GetSiret(row[numeroCompteIndex], date); err == nil {

          debit := Debit{
            key:                       siret,
            NumeroCompte:              row[numeroCompteIndex],
            NumeroEcartNegatif:        row[numeroEcartNegatifIndex],
            CodeProcedureCollective:   row[codeProcedureCollectiveIndex],
            CodeOperationEcartNegatif: row[codeOperationEcartNegatifIndex],
            CodeMotifEcartNegatif:     row[codeMotifEcartNegatifIndex],
          }

          debit.DateTraitement, err = urssafToDate(row[dateTraitementIndex])
          tracker.Error(err)
          debit.PartOuvriere, err = strconv.ParseFloat(row[partOuvriereIndex], 64)
          tracker.Error(err)
          debit.PartOuvriere = debit.PartOuvriere / 100
          debit.PartPatronale, err = strconv.ParseFloat(row[partPatronaleIndex], 64)
          tracker.Error(err)
          debit.PartPatronale = debit.PartPatronale / 100
          debit.NumeroHistoriqueEcartNegatif, err = strconv.Atoi(row[numeroHistoriqueEcartNegatifIndex])
          tracker.Error(err)
          debit.EtatCompte, err = strconv.Atoi(row[etatCompteIndex])
          tracker.Error(err)
          debit.Periode, err = urssafToPeriod(row[periodeIndex])
          tracker.Error(err)
          // debit.MontantMajorations, err = strconv.ParseFloat(row[montantMajorationsIndex], 64)
          // tracker.Error(err)
          // debit.MontantMajorations = debit.MontantMajorations / 100

          if !tracker.ErrorInCycle() {
            outputChannel <- debit
          } else {
            //event.Debug(tracker.Report("errors"))
          }
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
