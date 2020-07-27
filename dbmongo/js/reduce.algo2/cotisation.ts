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
  output_indexed: ParPériode<Input & Partial<SortieCotisation>>
): ParPériode<SortieCotisation> {
  "use strict"

  const f = { generatePeriodSerie, dateAddMonth } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const sortieCotisation: ParPériode<SortieCotisation> = {}

  // calcul de cotisation_moyenne sur 12 mois
  Object.keys(output_indexed).forEach((periode) => {
    const input = output_indexed[periode]
    const periode_courante = output_indexed[periode].periode

    const periode_12_mois = f.dateAddMonth(periode_courante, 12)
    const series = f.generatePeriodSerie(periode_courante, periode_12_mois)
    series.forEach((periodeFuture) => {
      if (periodeFuture.getTime() in output_indexed) {
        const outputInFuture = (sortieCotisation[
          periodeFuture.getTime()
        ] = sortieCotisation[periodeFuture.getTime()] || {
          cotisation_array: [],
          montant_pp_array: [],
          montant_po_array: [],
        })
        if (input.cotisation !== undefined)
          outputInFuture.cotisation_array.push(input.cotisation)
        if (input.montant_part_patronale !== undefined)
          outputInFuture.montant_pp_array.push(input.montant_part_patronale)
        if (input.montant_part_ouvriere !== undefined)
          outputInFuture.montant_po_array.push(input.montant_part_ouvriere)
      }
    })

    const outputInPeriod = sortieCotisation[periode]
    outputInPeriod.cotisation_array = outputInPeriod.cotisation_array || []
    outputInPeriod.cotisation_moy12m =
      outputInPeriod.cotisation_array.reduce((p, c) => p + c, 0) /
      (outputInPeriod.cotisation_array.length || 1)
    if (
      outputInPeriod.cotisation_moy12m > 0 &&
      input.montant_part_ouvriere !== undefined &&
      input.montant_part_patronale !== undefined
    ) {
      outputInPeriod.ratio_dette =
        (input.montant_part_ouvriere + input.montant_part_patronale) /
        outputInPeriod.cotisation_moy12m
      const pp_average =
        (outputInPeriod.montant_pp_array || []).reduce((p, c) => p + c, 0) /
        (outputInPeriod.montant_pp_array?.length || 1)
      const po_average =
        (outputInPeriod.montant_po_array || []).reduce((p, c) => p + c, 0) /
        (outputInPeriod.montant_po_array?.length || 1)
      outputInPeriod.ratio_dette_moy12m =
        (po_average + pp_average) / outputInPeriod.cotisation_moy12m
    }
    // Remplace dans cibleApprentissage
    //val.dette_any_12m = (val.montant_pp_array || []).reduce((p,c) => (c >=
    //100) || p, false) || (val.montant_po_array || []).reduce((p, c) => (c >=
    //100) || p, false)
    delete outputInPeriod.cotisation_array
    delete outputInPeriod.montant_pp_array
    delete outputInPeriod.montant_po_array
  })

  // Calcul des défauts URSSAF prolongés
  let counter = 0
  Object.keys(sortieCotisation)
    .sort()
    .forEach((k) => {
      const { ratio_dette } = sortieCotisation[k]
      if (!ratio_dette) return
      if (ratio_dette > 0.01) {
        sortieCotisation[k].tag_debit = true // Survenance d'un débit d'au moins 1% des cotisations
      }
      if (ratio_dette > 1) {
        counter = counter + 1
        if (counter >= 3) sortieCotisation[k].tag_default = true
      } else counter = 0
    })

  return sortieCotisation
}

/* TODO: appliquer même logique d'itération sur futureTimestamps que dans cotisationsdettes.ts */
