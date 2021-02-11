import { f } from "./functions"
import { EntréeDiane } from "../GeneratedTypes"
import { EntréeBdf, ParHash, ParPériode, Timestamp } from "../RawDataTypes"

type YearOffset = 1 | 2

type DonnéesDianeTransmises = Omit<
  EntréeDiane,
  | "marquee"
  | "nom_entreprise"
  | "numero_siren"
  | "statut_juridique"
  | "procedure_collective"
>
type CléRatioDianePassé = `${keyof DonnéesDianeTransmises}_past_${YearOffset}`

type RatioBdfCalculable = Pick<
  EntréeBdf,
  "poids_frng" | "dette_fiscale" | "frais_financier"
>

type RatioBdfInclus = RatioBdfCalculable &
  Pick<EntréeBdf, "taux_marge" | "financier_court_terme">
type CléRatioBdfPassé = `${keyof RatioBdfInclus}_past_${YearOffset}`

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "diane"
  computed: Record<CléRatioDianePassé, number | null | undefined> &
    RatioBdfInclus &
    Record<CléRatioBdfPassé, number>
  transmitted: DonnéesDianeTransmises
}

export type SortieDiane = Variables["computed"] & Variables["transmitted"]

export function entr_diane(
  donnéesDiane: ParHash<EntréeDiane>,
  output_indexed: ParPériode<Partial<SortieDiane>>,
  periodes: Timestamp[]
): ParPériode<Partial<SortieDiane>> {
  for (const entréeDiane of Object.values(donnéesDiane)) {
    if (!entréeDiane.arrete_bilan_diane) continue
    //entréeDiane.arrete_bilan_diane = new Date(Date.UTC(entréeDiane.exercice_diane, 11, 31, 0, 0, 0, 0))
    const periode_arrete_bilan = new Date(
      Date.UTC(
        entréeDiane.arrete_bilan_diane.getUTCFullYear(),
        entréeDiane.arrete_bilan_diane.getUTCMonth() + 1,
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
      const rest: DonnéesDianeTransmises = f.omit(
        entréeDiane,
        // "marquee",
        "nom_entreprise",
        "numero_siren",
        "statut_juridique",
        "procedure_collective"
      )

      const makePastProp = (
        prop: keyof DonnéesDianeTransmises,
        offset: YearOffset
      ) => `${prop}_past_${offset}` as CléRatioDianePassé

      if (periodes.includes(periode.getTime())) {
        Object.assign(output_indexed[periode.getTime()], rest)
      }

      for (const ratio of Object.keys(rest) as (keyof typeof rest)[]) {
        if (entréeDiane[ratio] === null) {
          const outputAtTime = output_indexed[periode.getTime()]
          if (
            outputAtTime !== undefined &&
            periodes.includes(periode.getTime())
          ) {
            delete outputAtTime[ratio]
          }
          continue
        }

        // Passé

        const past_year_offset: YearOffset[] = [1, 2]
        for (const offset of past_year_offset) {
          const periode_offset = f.dateAddMonth(periode, 12 * offset)
          const variable_name: CléRatioDianePassé = makePastProp(ratio, offset)

          const outputAtOffset = output_indexed[periode_offset.getTime()]
          if (
            outputAtOffset !== undefined &&
            ratio !== "arrete_bilan_diane" &&
            ratio !== "exercice_diane"
          ) {
            outputAtOffset[variable_name] = entréeDiane[ratio]
          }
        }
      }
    }

    for (const periode of series) {
      const inputInPeriod = output_indexed[periode.getTime()]
      const outputInPeriod = output_indexed[periode.getTime()]
      if (
        periodes.includes(periode.getTime()) &&
        inputInPeriod &&
        outputInPeriod
      ) {
        // Recalcul BdF si ratios bdf sont absents
        if (!("poids_frng" in inputInPeriod)) {
          const poids = f.poidsFrng(entréeDiane)
          if (poids !== null) outputInPeriod.poids_frng = poids
        }
        if (!("dette_fiscale" in inputInPeriod)) {
          const dette = f.detteFiscale(entréeDiane)
          if (dette !== null) outputInPeriod.dette_fiscale = dette
        }
        if (!("frais_financier" in inputInPeriod)) {
          const frais = f.fraisFinancier(entréeDiane)
          if (frais !== null) outputInPeriod.frais_financier = frais
        }

        // TODO: mettre en commun population des champs _past_ avec bdf ?
        const bdf_vars: (keyof RatioBdfInclus)[] = [
          "taux_marge",
          "poids_frng",
          "dette_fiscale",
          "financier_court_terme",
          "frais_financier",
        ]
        const past_year_offset: YearOffset[] = [1, 2]

        const makePastProp = (clé: keyof RatioBdfInclus, offset: YearOffset) =>
          `${clé}_past_${offset}` as CléRatioBdfPassé

        bdf_vars.forEach((k) => {
          if (k in outputInPeriod) {
            past_year_offset.forEach((offset) => {
              const periode_offset = f.dateAddMonth(periode, 12 * offset)
              const variable_name: CléRatioBdfPassé = makePastProp(k, offset)

              const outputAtOffset = output_indexed[periode_offset.getTime()]
              if (
                outputAtOffset &&
                periodes.includes(periode_offset.getTime())
              ) {
                outputAtOffset[variable_name] = outputInPeriod[k]
              }
            })
          }
        })
      }
    }
  }
  return output_indexed
}
