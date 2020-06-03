import * as f from "../common/generatePeriodSerie.js"

// Définition dérivée de dbmongo/lib/urssaf/delai.go (seulement propriétés nécéssaires)
export type Delai = {
  date_creation: Date
  date_echeance: Date
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  montant_echeancier: number // exprimé en euros
}

// Valeurs ajoutées dans la paramètre indexed_output passé à delais()
export type DelaiComputedValues = {
  delai: number
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  ratio_dette_delai?: number
  montant_echeancier: number // exprimé en euros
}

export type DelaiMap = { [key: string]: Delai }

// Valeurs attendues dans le paramètre indexed_output passé à delais()
export type DebitComputedValues = {
  montant_part_patronale?: number
  montant_part_ouvriere?: number
}

// Type du paramètre donnéesActuellesParPériode passé à delais()
export type DebitComputedValuesPerPeriod = {
  [time: string]: DebitComputedValues
}

type DeepReadonly<T> = Readonly<T> // pas vraiment, mais espoire que TS le supporte prochainement

export function delais(
  v: { delai: DelaiMap },
  donnéesActuellesParPériode: DeepReadonly<DebitComputedValuesPerPeriod>
): { [time: string]: DelaiComputedValues } {
  "use strict"
  const donnéesSupplémentairesParPériode: {
    [time: string]: DelaiComputedValues
  } = {}
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
      if (time in donnéesActuellesParPériode) {
        const remaining_months =
          date_echeance.getUTCMonth() -
          new Date(time).getUTCMonth() +
          12 *
            (date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear())
        const inputAtTime = donnéesActuellesParPériode[time]
        const outputAtTime: DelaiComputedValues = {
          delai: remaining_months,
          duree_delai: delai.duree_delai,
          montant_echeancier: delai.montant_echeancier,
        }
        if (
          delai.duree_delai > 0 &&
          inputAtTime.montant_part_patronale !== undefined &&
          inputAtTime.montant_part_ouvriere !== undefined
        ) {
          outputAtTime.ratio_dette_delai =
            (inputAtTime.montant_part_patronale +
              inputAtTime.montant_part_ouvriere -
              (delai.montant_echeancier * remaining_months * 30) /
                delai.duree_delai) /
            delai.montant_echeancier
        }
        donnéesSupplémentairesParPériode[time] = outputAtTime
      }
    })
  })
  return donnéesSupplémentairesParPériode
}
