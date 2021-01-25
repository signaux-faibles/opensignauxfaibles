import { f } from "./functions"
import { ParPériode } from "../RawDataTypes"

export type Input = {
  periode: Date
  cotisation?: number
  montant_part_patronale?: number
  montant_part_ouvriere?: number
}

export type SortieCotisation = {
  /** Montant moyen de cotisations calculé sur 12 mois consécutifs. */
  cotisation_moy12m?: number
  /** ratio_dette = (montant_part_ouvriere + montant_part_patronale) / cotisation_moy12m */
  ratio_dette: number
  /** Moyenne de ratio_dette sur 12 mois. */
  ratio_dette_moy12m?: number
  /** Survenance d'un débit d'au moins 1% des cotisations */
  tag_debit: boolean
  /** Survenance de trois débits de 100% (ou plus) des cotisations */
  tag_default: boolean
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "cotisation"
  computed: SortieCotisation
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function cotisation(
  output_indexed: ParPériode<Input & Partial<SortieCotisation>>
): ParPériode<SortieCotisation> {
  "use strict"

  const sortieCotisation: ParPériode<SortieCotisation> = {}

  const moyenne = (valeurs: (number | undefined)[] = []): number | undefined =>
    valeurs.some((val) => typeof val === "undefined")
      ? undefined
      : (valeurs as number[]).reduce((p, c) => p + c, 0) / (valeurs.length || 1)

  // calcul de cotisation_moyenne sur 12 mois
  const futureArrays: ParPériode<{
    cotisations: (number | undefined)[]
    montantsPP: number[]
    montantsPO: number[]
  }> = {}

  for (const [periode, input] of Object.entries(output_indexed)) {
    const périodeCourante = output_indexed[periode]?.periode
    if (périodeCourante === undefined) continue

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
      future.montantsPP.push(input.montant_part_patronale || 0)
      future.montantsPO.push(input.montant_part_ouvriere || 0)
    })

    // Calcul des cotisations moyennes à partir des valeurs accumulées ci-dessus
    const { cotisations, montantsPO, montantsPP } = futureArrays[periode] ?? {}
    const out = sortieCotisation[periode] ?? ({} as SortieCotisation)
    if (cotisations && cotisations.length >= 12) {
      out.cotisation_moy12m = moyenne(cotisations)
    }
    if (typeof out.cotisation_moy12m === "undefined") {
      delete out.cotisation_moy12m
    } else if (out.cotisation_moy12m > 0) {
      out.ratio_dette =
        ((input.montant_part_ouvriere || 0) +
          (input.montant_part_patronale || 0)) /
        out.cotisation_moy12m
      if (
        montantsPO &&
        montantsPP &&
        cotisations &&
        !cotisations.includes(undefined) &&
        !cotisations.includes(0)
      ) {
        const detteVals = []
        for (const [i, cotisation] of cotisations.entries()) {
          const montPO = montantsPO[i]
          const montPP = montantsPP[i]
          if (
            cotisation !== undefined &&
            montPO !== undefined &&
            montPP !== undefined
          ) {
            detteVals.push((montPO + montPP) / cotisation)
          }
        }
        out.ratio_dette_moy12m = moyenne(detteVals)
      }
    }
    sortieCotisation[periode] = out
    // Remplace dans cibleApprentissage
    //val.dette_any_12m = (val.montantsPA || []).reduce((p,c) => (c >=
    //100) || p, false) || (val.montantsPO || []).reduce((p, c) => (c >=
    //100) || p, false)
  }

  // Calcul des défauts URSSAF prolongés
  let counter = 0
  for (const k of Object.keys(sortieCotisation).sort()) {
    const cotis = sortieCotisation[k] as SortieCotisation
    const { ratio_dette } = sortieCotisation[k] ?? {}
    if (!ratio_dette) continue
    if (ratio_dette > 0.01) {
      cotis.tag_debit = true // Survenance d'un débit d'au moins 1% des cotisations
    }
    if (ratio_dette > 1) {
      counter = counter + 1
      if (counter >= 3) cotis.tag_default = true
    } else counter = 0
  }

  return sortieCotisation
}
