import "../globals"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

type SortieBdf = {
  annee_bdf: number
  exercice_bdf: number // année
}

export function entr_bdf(entréeBdf: DonnéesBdf): ParPériode<SortieBdf> {
  const outputBdf: ParPériode<SortieBdf> = {}

  const f = { generatePeriodSerie, dateAddMonth } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  // Retourne les clés de obj, en respectant le type défini dans le type de obj.
  // Contrat: obj ne doit contenir que les clés définies dans son type.
  const typedObjectKeys = <T>(obj: T): Array<keyof T> =>
    Object.keys(obj) as Array<keyof T>

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

  for (const hash of typedObjectKeys(entréeBdf.bdf)) {
    const periode_arrete_bilan = new Date(
      Date.UTC(
        entréeBdf.bdf[hash].arrete_bilan_bdf.getUTCFullYear(),
        entréeBdf.bdf[hash].arrete_bilan_bdf.getUTCMonth() + 1,
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
      const bdfHashData = entréeBdf.bdf[hash]
      const outputInPeriod = outputBdf[periode.getTime()]
      const rest = omit(
        bdfHashData as EntréeBdf & {
          raison_sociale: unknown
          secteur: unknown
          siren: unknown
        },
        "raison_sociale",
        "secteur",
        "siren"
      )

      if (outputInPeriod) {
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
            periode_offset.getTime() in outputBdf &&
            k !== "arrete_bilan_bdf" &&
            k !== "exercice_bdf"
          ) {
            outputBdf[periode_offset.getTime()][variable_name] =
              entréeBdf.bdf[hash][k]
          }
        }
      }
    }
  }

  return outputBdf
}
