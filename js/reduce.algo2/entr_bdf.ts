import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { EntréeBdf } from "../GeneratedTypes"
import { ParHash, Timestamp } from "../RawDataTypes"

type EntréeBdfRatios = Pick<
  EntréeBdf,
  | "poids_frng"
  | "taux_marge"
  | "delai_fournisseur"
  | "dette_fiscale"
  | "financier_court_terme"
  | "frais_financier"
>

type DonnéesBdfTransmises = Omit<
  EntréeBdf,
  "raison_sociale" | "secteur" | "siren"
>

type YearOffset = 1 | 2
type CléRatiosBdfPassés = `${keyof EntréeBdfRatios}_past_${YearOffset}`
type RatiosBdfPassés = Record<CléRatiosBdfPassés, number>
type DonnéesBdfCalculées = {
  /** Année de l'exercice - 1. */
  exercice_bdf: number
} & EntréeBdfRatios &
  RatiosBdfPassés

export type SortieBdf = DonnéesBdfTransmises & DonnéesBdfCalculées

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "entr_bdf"
  computed: DonnéesBdfCalculées
  transmitted: DonnéesBdfTransmises
}

export function entr_bdf(
  donnéesBdf: ParHash<EntréeBdf>,
  periodes: Timestamp[]
): ParPériode<Partial<SortieBdf>> {
  "use strict"

  const outputBdf = f.newParPériode<Partial<SortieBdf>>()
  for (const p of periodes) {
    outputBdf.set(p, {})
  }

  for (const entréeBdf of Object.values(donnéesBdf)) {
    const periode_arrete_bilan = new Date(
      Date.UTC(
        entréeBdf.arrete_bilan_bdf.getUTCFullYear(),
        entréeBdf.arrete_bilan_bdf.getUTCMonth() + 1,
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
      const outputInPeriod = outputBdf.get(periode) || {}

      const periodData: DonnéesBdfTransmises = f.omit(
        entréeBdf,
        "raison_sociale",
        "secteur",
        "siren"
      )

      // TODO: Éviter d'ajouter des données en dehors de `periodes`, sans fausser le calcul des données passées (plus bas)
      Object.assign(outputInPeriod, periodData)
      if (outputInPeriod.annee_bdf) {
        outputInPeriod.exercice_bdf = outputInPeriod.annee_bdf - 1
      }

      const pastData = f.omit(periodData, "arrete_bilan_bdf", "annee_bdf")

      const makePastProp = (prop: keyof EntréeBdfRatios, offset: YearOffset) =>
        `${prop}_past_${offset}` as CléRatiosBdfPassés
      for (const prop of Object.keys(pastData) as (keyof typeof pastData)[]) {
        const past_year_offset: YearOffset[] = [1, 2]
        for (const offset of past_year_offset) {
          const periode_offset = f.dateAddMonth(periode, 12 * offset)
          const outputInPast = outputBdf.get(periode_offset)
          if (outputInPast) {
            outputInPast[makePastProp(prop, offset)] = entréeBdf[prop]
          }
        }
      }
      outputBdf.set(periode, outputInPeriod)
    }
  }

  return outputBdf
}
