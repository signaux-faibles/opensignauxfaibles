// Object golang défini dans dbmongo/lib/urssaf/delai.go
// NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
// NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
// DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
// DateEcheance      time.Time `json:"date_echeance" bson:"date_echeance"`
// DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
// Denomination      string    `json:"denomination" bson:"denomination"`
// Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
// AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
// MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
// Stade             string    `json:"stade" bson:"stade"`
// Action            string    `json:"action" bson:"action"`

type Delai = {
  numero_compte: string
  numero_contentieux: string
  date_creation: Date
  date_echeance: Date
  duree_delai: number
  denomination: string
  indic_6m: string
  annee_creation: number
  montant_echeancier: number
  stade: string
  action: string
}

type DelaiMap = { [key: string]: Delai }

declare const f: {
  generatePeriodSerie(date_creation: Date, date_echeance: Date): Date[]
}
export function delais(v: { delai: DelaiMap }, output_indexed: object): void {
  "use strict"
  Object.keys(v.delai).map(function (hash) {
    const delai = v.delai[hash]
    const date_creation = new Date(
      Date.UTC(
        delai.date_creation.getUTCFullYear(),
        delai.date_creation.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    const date_echeance = new Date(
      Date.UTC(
        delai.date_echeance.getUTCFullYear(),
        delai.date_echeance.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    const pastYearTimes = f
      .generatePeriodSerie(date_creation, date_echeance)
      .map(function (date) {
        return date.getTime()
      })
    pastYearTimes.map(function (time) {
      if (time in output_indexed) {
        const remaining_months =
          date_echeance.getUTCMonth() -
          new Date(time).getUTCMonth() +
          12 *
            (date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear())
        output_indexed[time].delai = remaining_months
        output_indexed[time].duree_delai = delai.duree_delai
        output_indexed[time].montant_echeancier = delai.montant_echeancier

        if (delai.duree_delai > 0) {
          output_indexed[time].ratio_dette_delai =
            (output_indexed[time].montant_part_patronale +
              output_indexed[time].montant_part_ouvriere -
              (delai.montant_echeancier * remaining_months * 30) /
                delai.duree_delai) /
            delai.montant_echeancier
        }
      }
    })
  })
}
