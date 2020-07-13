import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

type Input = {
  periode: Date
  cotisation?: number
  montant_part_patronale?: number
  montant_part_ouvriere?: number
}

export type SortieCotisation = {
  montant_pp_array: number[]
  montant_po_array: number[]
  cotisation_moy12m: number
  cotisation_array: number[]
  ratio_dette: number
  ratio_dette_moy12m: number
  tag_debit: boolean
  tag_default: boolean
}

export function cotisation(
  output_indexed: { [k: string]: Input & Partial<SortieCotisation> },
  output_array: (Input & Partial<SortieCotisation>)[]
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
        const outputInPeriod = output_indexed[periode.getTime()]
        const outputCourante = output_indexed[periode_courante.getTime()]
        if (outputCourante.cotisation !== undefined)
          outputInPeriod.cotisation_array = (
            outputInPeriod.cotisation_array || []
          ).concat(outputCourante.cotisation)
        if (outputCourante.montant_part_patronale !== undefined)
          outputInPeriod.montant_pp_array = (
            outputInPeriod.montant_pp_array || []
          ).concat(outputCourante.montant_part_patronale)
        if (outputCourante.montant_part_ouvriere !== undefined)
          outputInPeriod.montant_po_array = (
            outputInPeriod.montant_po_array || []
          ).concat(outputCourante.montant_part_ouvriere)
      }
    })
  })

  for (const val of output_array) {
    val.cotisation_array = val.cotisation_array || []
    val.cotisation_moy12m =
      val.cotisation_array.reduce((p, c) => p + c, 0) /
      (val.cotisation_array.length || 1)
    if (
      val.cotisation_moy12m > 0 &&
      val.montant_part_ouvriere !== undefined &&
      val.montant_part_patronale !== undefined
    ) {
      val.ratio_dette =
        (val.montant_part_ouvriere + val.montant_part_patronale) /
        val.cotisation_moy12m
      const pp_average =
        (val.montant_pp_array || []).reduce((p, c) => p + c, 0) /
        (val.montant_pp_array?.length || 1)
      const po_average =
        (val.montant_po_array || []).reduce((p, c) => p + c, 0) /
        (val.montant_po_array?.length || 1)
      val.ratio_dette_moy12m = (po_average + pp_average) / val.cotisation_moy12m
    }
    // Remplace dans cibleApprentissage
    //val.dette_any_12m = (val.montant_pp_array || []).reduce((p,c) => (c >=
    //100) || p, false) || (val.montant_po_array || []).reduce((p, c) => (c >=
    //100) || p, false)
    delete val.cotisation_array
    delete val.montant_pp_array
    delete val.montant_po_array
  }

  // Calcul des défauts URSSAF prolongés
  let counter = 0
  Object.keys(output_indexed)
    .sort()
    .forEach((k) => {
      const { ratio_dette } = output_indexed[k]
      if (!ratio_dette) return
      if (ratio_dette > 0.01) {
        output_indexed[k].tag_debit = true // Survenance d'un débit d'au moins 1% des cotisations
      }
      if (ratio_dette > 1) {
        counter = counter + 1
        if (counter >= 3) output_indexed[k].tag_default = true
      } else counter = 0
    })
}

/* TODO: appliquer même logique d'itération sur futureTimestamps que dans cotisationsdettes.ts */
