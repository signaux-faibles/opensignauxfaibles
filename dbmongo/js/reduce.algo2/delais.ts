import * as f from "../common/generatePeriodSerie.js"

// Définition dérivée de dbmongo/lib/urssaf/delai.go (seulement propriétés nécéssaires)
export type Delai = {
  date_creation: Date
  date_echeance: Date
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  montant_echeancier: number // exprimé en euros
}

export type DelaiMap = { [key: string]: Delai }

// Valeurs attendues dans le paramètre indexed_output passé à delais()
export type IndexedOutputExpectedValues = {
  montant_part_patronale?: number
  montant_part_ouvriere?: number
}

// Valeurs ajoutées dans la paramètre indexed_output passé à delais()
export type DelaiComputedValues = {
  delai: number
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  ratio_dette_delai: number
  montant_echeancier: number // exprimé en euros
}

// Type du paramètre indexed_output passé à delais()
export type IndexedOutputPartial = {
  [time: string]: IndexedOutputExpectedValues & Partial<DelaiComputedValues>
}

export function delais(
  v: { delai: DelaiMap },
  output_indexed: IndexedOutputPartial
): void {
  "use strict"
  Object.keys(v.delai).map(function (hash) {
    const delai = v.delai[hash]
    // On arrondit les dates au premier jour du mois.
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
    // Création d'un tableau de timestamps à raison de 1 par mois.
    const pastYearTimes = f
      .generatePeriodSerie(date_creation, date_echeance)
      .map(function (date: Date) {
        return date.getTime()
      })
    pastYearTimes.map(function (time: number) {
      if (time in output_indexed) {
        const remaining_months =
          date_echeance.getUTCMonth() -
          new Date(time).getUTCMonth() +
          12 *
            (date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear())
        output_indexed[time].delai = remaining_months
        output_indexed[time].duree_delai = delai.duree_delai
        output_indexed[time].montant_echeancier = delai.montant_echeancier

        const outputAtTime = output_indexed[time]
        if (
          delai.duree_delai > 0 &&
          outputAtTime.montant_part_patronale !== undefined &&
          outputAtTime.montant_part_ouvriere !== undefined
        ) {
          output_indexed[time].ratio_dette_delai =
            (outputAtTime.montant_part_patronale +
              outputAtTime.montant_part_ouvriere -
              (delai.montant_echeancier * remaining_months * 30) /
                delai.duree_delai) /
            delai.montant_echeancier
        }
      }
    })
  })
}
