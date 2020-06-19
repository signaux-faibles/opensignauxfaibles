import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

export type Output = {
  periode: Date
  cotisation: number
  cotisation_moy12m: number
  cotisation_array: number[]
  montant_part_patronale: number
  montant_pp_array: number[]
  montant_part_ouvriere: number
  montant_po_array: number[]
  ratio_dette: number
  ratio_dette_moy12m: number
  tag_debit: boolean
  tag_default: boolean
}

export function cotisation(
  output_indexed: { [k: string]: Output },
  output_array: Output[]
): void {
  "use strict"
  const f = { generatePeriodSerie, dateAddMonth } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  // calcul de cotisation_moyenne sur 12 mois
  Object.keys(output_indexed).forEach((k) => {
    const periode_courante = output_indexed[k].periode
    const periode_12_mois = f.dateAddMonth(periode_courante, 12)
    const series = f.generatePeriodSerie(periode_courante, periode_12_mois)
    series.forEach((periode) => {
      if (periode.getTime() in output_indexed) {
        if ("cotisation" in output_indexed[periode_courante.getTime()])
          output_indexed[periode.getTime()].cotisation_array = (
            output_indexed[periode.getTime()].cotisation_array || []
          ).concat(output_indexed[periode_courante.getTime()].cotisation)

        output_indexed[periode.getTime()].montant_pp_array = (
          output_indexed[periode.getTime()].montant_pp_array || []
        ).concat(
          output_indexed[periode_courante.getTime()].montant_part_patronale
        )
        output_indexed[periode.getTime()].montant_po_array = (
          output_indexed[periode.getTime()].montant_po_array || []
        ).concat(
          output_indexed[periode_courante.getTime()].montant_part_ouvriere
        )
      }
    })
  })

  output_array.forEach((val) => {
    val.cotisation_array = val.cotisation_array || []
    val.cotisation_moy12m =
      val.cotisation_array.reduce((p, c) => p + c, 0) /
      (val.cotisation_array.length || 1)
    if (val.cotisation_moy12m > 0) {
      val.ratio_dette =
        (val.montant_part_ouvriere + val.montant_part_patronale) /
        val.cotisation_moy12m
      const pp_average =
        (val.montant_pp_array || []).reduce((p, c) => p + c, 0) /
        (val.montant_pp_array.length || 1)
      const po_average =
        (val.montant_po_array || []).reduce((p, c) => p + c, 0) /
        (val.montant_po_array.length || 1)
      val.ratio_dette_moy12m = (po_average + pp_average) / val.cotisation_moy12m
    }
    // Remplace dans cibleApprentissage
    //val.dette_any_12m = (val.montant_pp_array || []).reduce((p,c) => (c >=
    //100) || p, false) || (val.montant_po_array || []).reduce((p, c) => (c >=
    //100) || p, false)
    delete val.cotisation_array
    delete val.montant_pp_array
    delete val.montant_po_array
  })

  // Calcul des défauts URSSAF prolongés
  let counter = 0
  Object.keys(output_indexed)
    .sort()
    .forEach((k) => {
      if (output_indexed[k].ratio_dette > 0.01) {
        output_indexed[k].tag_debit = true // Survenance d'un débit d'au moins 1% des cotisations
      }
      if (output_indexed[k].ratio_dette > 1) {
        counter = counter + 1
        if (counter >= 3) output_indexed[k].tag_default = true
      } else counter = 0
    })
}
