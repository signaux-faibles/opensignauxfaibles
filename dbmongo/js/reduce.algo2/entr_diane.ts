import "../globals"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"
import { omit } from "../common/omit"
import { poidsFrng } from "./poidsFrng"
import { detteFiscale } from "./detteFiscale"
import { fraisFinancier } from "./fraisFinancier"

export type SortieDiane = Record<string, unknown> // for *_past_* props of diane. // TODO: try to be more specific

export function entr_diane(
  donnéesDiane: Record<DataHash, EntréeDiane>,
  output_indexed: ParPériode<SortieDiane>,
  periodes: Timestamp[]
): ParPériode<SortieDiane> {
  /* DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO */ const f = {
    ...{ generatePeriodSerie, dateAddMonth, omit, poidsFrng, detteFiscale }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ fraisFinancier }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  for (const hash of Object.keys(donnéesDiane)) {
    if (!donnéesDiane[hash].arrete_bilan_diane) continue
    //donnéesDiane[hash].arrete_bilan_diane = new Date(Date.UTC(donnéesDiane[hash].exercice_diane, 11, 31, 0, 0, 0, 0))
    const periode_arrete_bilan = new Date(
      Date.UTC(
        donnéesDiane[hash].arrete_bilan_diane.getUTCFullYear(),
        donnéesDiane[hash].arrete_bilan_diane.getUTCMonth() + 1,
        1,
        0,
        0,
        0,
        0
      )
    )
    const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7) // 01/08 pour un bilan le 31/12, donc algo qui tourne en 01/09
    const series = f.generatePeriodSerie(
      periode_dispo,
      f.dateAddMonth(periode_dispo, 14) // periode de validité d'un bilan auprès de la Banque de France: 21 mois (14+7)
    )

    for (const periode of series) {
      const rest = f.omit(
        donnéesDiane[hash] as EntréeDiane & {
          marquee: unknown
          nom_entreprise: unknown
          numero_siren: unknown
          statut_juridique: unknown
          procedure_collective: unknown
        },
        "marquee",
        "nom_entreprise",
        "numero_siren",
        "statut_juridique",
        "procedure_collective"
      )

      if (periodes.includes(periode.getTime())) {
        Object.assign(output_indexed[periode.getTime()], rest)
      }

      for (const ratio of Object.keys(rest) as (keyof typeof rest)[]) {
        if (donnéesDiane[hash][ratio] === null) {
          if (periodes.includes(periode.getTime())) {
            delete output_indexed[periode.getTime()][ratio]
          }
          continue
        }

        // Passé

        const past_year_offset = [1, 2]
        for (const offset of past_year_offset) {
          const periode_offset = f.dateAddMonth(periode, 12 * offset)
          const variable_name = ratio + "_past_" + offset

          if (
            periode_offset.getTime() in output_indexed &&
            ratio !== "arrete_bilan_diane" &&
            ratio !== "exercice_diane"
          ) {
            output_indexed[periode_offset.getTime()][variable_name] =
              donnéesDiane[hash][ratio]
          }
        }
      }
    }
  }
  return output_indexed
}
