import "../globals"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

export type SortieBdf = {
  annee_bdf: number
  exercice_bdf: number // année
} & RatiosBdf &
  RatiosBdfPassés &
  Record<string, unknown> // for *_past_* props of bdf. // TODO: try to be more specific

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
  v: DonnéesBdf, // TODO: prendre ParPériode<EntréeBdf> au lieu de DonnéesBdf
  output_indexed: ParPériode<Partial<SortieBdf>>,
  periodes: Timestamp[]
): ParPériode<Partial<SortieBdf>> {
  "use strict"
  periodes
  const outputBdf: ParPériode<Partial<SortieBdf>> = { ...output_indexed }

  const f = { generatePeriodSerie, dateAddMonth } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  // Fonction pour omettre des props, tout en retournant le bon type
  function omit<Source, Exclusions extends Array<keyof Source>>(
    object: Source,
    ...propNames: Exclusions
  ): Omit<Source, Exclusions[number]> {
    const result: Omit<Source, Exclusions[number]> = Object.assign({}, object)
    for (const prop of propNames) {
      delete (result as any)[prop]
    }
    return result
  }
  // Fonction pour omettre des props, tout en retournant le bon type
  function omitProps<Source, Exclusions extends Array<keyof Source>>(
    object: Source,
    ...propNames: Exclusions
  ): Array<keyof Omit<Source, Exclusions[number]>> {
    const result = omit(object, ...propNames)
    return Object.keys(result) as (keyof typeof result)[]
  }
  // TODO: [refacto] extraire dans common/ ou reduce.algo2/

  for (const hash in v.bdf) {
    const bdfHashData = v.bdf[hash]
    const periode_arrete_bilan = new Date(
      Date.UTC(
        bdfHashData.arrete_bilan_bdf.getUTCFullYear(),
        bdfHashData.arrete_bilan_bdf.getUTCMonth() + 1,
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

      //if (outputInPeriod || periode.getTime() in periodes) {
      Object.assign(
        outputInPeriod,
        omit(bdfHashData, "raison_sociale", "secteur", "siren")
      )
      if (outputInPeriod.annee_bdf) {
        outputInPeriod.exercice_bdf = outputInPeriod.annee_bdf - 1
      }
      //}

      const excludedProps = [
        "raison_sociale",
        "secteur",
        "siren",
        "arrete_bilan_bdf",
        "exercice_bdf",
      ] as const

      for (const prop of omitProps(bdfHashData, ...excludedProps)) {
        const past_year_offset = [1, 2]
        for (const offset of past_year_offset) {
          const periode_offset = f.dateAddMonth(periode, 12 * offset)
          const outputInPast = outputBdf[periode_offset.getTime()]
          if (outputInPast) {
            outputInPast[prop + "_past_" + offset] = v.bdf[hash][prop]
          }
        }
      }
    }
  }

  return outputBdf
}
