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
  output_indexed: Record<Periode, Partial<SortieBdf>>
  // periodes: Timestamp[]
): void /*ParPériode<SortieBdf>*/ {
  "use strict"
  // const outputBdf: ParPériode<SortieBdf> = {}
  // const outputBdf = entréeBdf.bdf

  const f = { generatePeriodSerie, dateAddMonth } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  /*
  // Retourne les clés de obj, en respectant le type défini dans le type de obj.
  // Contrat: obj ne doit contenir que les clés définies dans son type.
  const typedObjectKeys = <T>(obj: T): Array<keyof T> =>
    Object.keys(obj) as Array<keyof T>
  */
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
  // TODO: [refacto] extraire dans common/ ou reduce.algo2/

  for (const hash in /*of typedObjectKeys*/ v.bdf) {
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
      const outputInPeriod = output_indexed[periode.getTime()]
      const rest = omit(bdfHashData, "raison_sociale", "secteur", "siren")

      if (outputInPeriod /*periode.getTime() in periodes*/) {
        Object.assign(outputInPeriod, rest)
        if (outputInPeriod.annee_bdf) {
          outputInPeriod.exercice_bdf = outputInPeriod.annee_bdf - 1
        }
      }

      for (const k of Object.keys(rest) as (keyof typeof rest)[]) {
        const past_year_offset = [1, 2]
        for (const offset of past_year_offset) {
          const periode_offset = f.dateAddMonth(periode, 12 * offset)
          const variable_name = k + "_past_" + offset
          if (
            periode_offset.getTime() in output_indexed &&
            // TODO: `in periodes` en récupérant un paramètre périodes.
            k !== "arrete_bilan_bdf" &&
            k !== "exercice_bdf"
            // TODO: props à inclure dans le omit ci-dessus ?
          ) {
            output_indexed[periode_offset.getTime()][variable_name] =
              v.bdf[hash][k]
            /*
            output_indexed[periode_offset.getTime()] = {
              ...output_indexed[periode_offset.getTime()],
              [variable_name]: v.bdf[hash][k],
            }
            */
          }
        }
      }
    }
  }

  // return outputBdf
}
