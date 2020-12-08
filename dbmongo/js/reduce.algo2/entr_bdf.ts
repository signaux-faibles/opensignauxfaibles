import { f } from "./functions"
import {
  EntréeBdfRatios,
  EntréeBdf,
  ParHash,
  Timestamp,
  ParPériode,
} from "../RawDataTypes"

export type SortieBdf = {
  annee_bdf: number
  exercice_bdf: number // année
} & EntréeBdfRatios &
  RatiosBdfPassés

// Synchroniser les propriétés avec celles de RatiosBdf
type RatiosBdfPassés = {
  poids_frng_past_1: number
  taux_marge_past_1: number
  delai_fournisseur_past_1: number
  dette_fiscale_past_1: number
  financier_court_terme_past_1: number
  frais_financier_past_1: number
  poids_frng_past_2: number
  taux_marge_past_2: number
  delai_fournisseur_past_2: number
  dette_fiscale_past_2: number
  financier_court_terme_past_2: number
  frais_financier_past_2: number
}

export function entr_bdf(
  donnéesBdf: ParHash<EntréeBdf>,
  periodes: Timestamp[]
): ParPériode<Partial<SortieBdf>> {
  "use strict"

  const outputBdf: ParPériode<Partial<SortieBdf>> = {}
  for (const p of periodes) {
    outputBdf[p] = {}
  }

  for (const entréeBdf of Object.values(donnéesBdf)) {
    const periode_arrete_bilan = new Date(
      Date.UTC(
        entréeBdf.arrete_bilan_bdf.getUTCFullYear(),
        entréeBdf.arrete_bilan_bdf.getUTCMonth() + 1,
        1,
        0,
        0,
        0,
        0
      )
    )
    const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7)
    const series = f.generatePeriodSerie(
      periode_dispo,
      f.dateAddMonth(periode_dispo, 13)
    )

    for (const periode of series) {
      const outputInPeriod = (outputBdf[periode.getTime()] =
        outputBdf[periode.getTime()] || {})

      const periodData = f.omit(entréeBdf, "raison_sociale", "secteur", "siren")

      // TODO: Éviter d'ajouter des données en dehors de `periodes`, sans fausser le calcul des données passées (plus bas)
      Object.assign(outputInPeriod, periodData)
      if (outputInPeriod.annee_bdf) {
        outputInPeriod.exercice_bdf = outputInPeriod.annee_bdf - 1
      }

      const pastData = f.omit(periodData, "arrete_bilan_bdf", "exercice_bdf")

      for (const prop of Object.keys(pastData) as (keyof typeof pastData)[]) {
        const past_year_offset = [1, 2]
        for (const offset of past_year_offset) {
          const periode_offset = f.dateAddMonth(periode, 12 * offset)
          const outputInPast = outputBdf[periode_offset.getTime()]
          if (outputInPast) {
            Object.assign(outputInPast, {
              [prop + "_past_" + offset]: entréeBdf[prop],
            })
          }
        }
      }
    }
  }

  return outputBdf
}
