import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

export type Input = {
  periode: Date
  cotisation?: number
  montant_part_patronale?: number
  montant_part_ouvriere?: number
}

export type SortieCotisation = {
  cotisation_moy12m?: number
  ratio_dette: number
  ratio_dette_moy12m?: number
  tag_debit: boolean
  tag_default: boolean
}

export function cotisation(
  output_indexed: ParPériode<Input & Partial<SortieCotisation>>
): ParPériode<SortieCotisation> {
  "use strict"

  const sortieCotisation: ParPériode<SortieCotisation> = {}

  const f = { generatePeriodSerie, dateAddMonth } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const moyenne = (valeurs: (number | undefined)[] = []): number | undefined =>
    valeurs.some((val) => typeof val === "undefined")
      ? undefined
      : (valeurs as number[]).reduce((p, c) => p + c, 0) / (valeurs.length || 1)

  // calcul de cotisation_moyenne sur 12 mois
  const futureArrays: ParPériode<{
    cotisations: (number | undefined)[]
    montantsPP: (number | undefined)[]
    montantsPO: (number | undefined)[]
  }> = {}

  Object.keys(output_indexed).forEach((periode) => {
    const input = output_indexed[periode]

    const périodeCourante = output_indexed[periode].periode
    const douzeMoisÀVenir = f
      .generatePeriodSerie(périodeCourante, f.dateAddMonth(périodeCourante, 12))
      .map((periodeFuture) => ({ timestamp: periodeFuture.getTime() }))
      .filter(({ timestamp }) => timestamp in output_indexed)

    // Accumulation de cotisations sur les 12 mois à venir, pour calcul des moyennes
    douzeMoisÀVenir.forEach(({ timestamp }) => {
      const future = (futureArrays[timestamp] = futureArrays[timestamp] || {
        cotisations: [],
        montantsPP: [],
        montantsPO: [],
      })
      future.cotisations.push(input.cotisation)
      future.montantsPP.push(input.montant_part_patronale)
      future.montantsPO.push(input.montant_part_ouvriere)
    })

    // Calcul des cotisations moyennes à partir des valeurs accumulées ci-dessus
    const { cotisations, montantsPO, montantsPP } = futureArrays[periode]
    const out = (sortieCotisation[periode] = sortieCotisation[periode] || {})
    out.cotisation_moy12m = moyenne(cotisations)
    if (
      typeof out.cotisation_moy12m !== "undefined" &&
      out.cotisation_moy12m > 0
    ) {
      out.ratio_dette =
        ((input.montant_part_ouvriere || 0) +
          (input.montant_part_patronale || 0)) /
        out.cotisation_moy12m
      const [moyPO, moyPP] = [moyenne(montantsPO), moyenne(montantsPP)]
      if (typeof moyPO === "number" && typeof moyPP === "number") {
        out.ratio_dette_moy12m = (moyPO + moyPP) / out.cotisation_moy12m
      }
    }
    // Remplace dans cibleApprentissage
    //val.dette_any_12m = (val.montantsPA || []).reduce((p,c) => (c >=
    //100) || p, false) || (val.montantsPO || []).reduce((p, c) => (c >=
    //100) || p, false)
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
