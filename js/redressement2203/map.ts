/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { f } from "./functions"
import {
  CompanyDataValues,
  SommesDettes,
  SortieRedressementUrssaf2203,
} from "../RawDataTypes"

export type SortieMap = SortieRedressementUrssaf2203

declare const dateStr: string
declare const dateFin: Date

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>

declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  "use strict"
  const dateDebutObservation = new Date(dateFin)
  const testDate = new Date(dateStr)
  const dateFinObservation = new Date(dateFin.getTime())
  dateDebutObservation.setFullYear(dateDebutObservation.getFullYear() - 1)

  const plafonnerDateObservation = (
    dateObs: Date,
    dateDebut: Date,
    dateFin: Date
  ) => {
    return Math.min(
      Math.max(dateObs.getTime(), dateDebut.getTime()),
      dateFin.getTime()
    )
  }

  const values = f.flatten(this.value, "2203")
  const beforeBatches = [] // TODO : renommer les variables
  const afterBatches = []

  if (values.debit) {
    for (const debit of Object.values(values.debit)) {
      debit.periode.start > testDate
        ? afterBatches.push(debit)
        : beforeBatches.push(debit)
    }
  }
  const dettesAnciennesParECN: SommesDettes = f.recupererDetteTotale(
    beforeBatches
  )

  // const dettesAnciennesDebutParECN: SommesDettes = f.recupererDetteTotale(
  //   beforeBatches.filter((b) => b.date_traitement <= testDate)
  // )

  const dettesRecentesParECN: SommesDettes = f.recupererDetteTotale(
    afterBatches
  )

  const cotisationMoyenne = f.cotisation(values.cotisation || {}, dateFin)

  // Jours de demande
  const totalMoisDemande =
    Object.values(values.apdemande || [])
      .filter((a) => a.motif_recours_se < 6)
      .reduce(
        (a, b) =>
          plafonnerDateObservation(
            b.periode.end,
            dateDebutObservation,
            dateFinObservation
          ) -
          plafonnerDateObservation(
            b.periode.start,
            dateDebutObservation,
            dateFinObservation
          ) +
          a,
        0
      ) /
    (3600 * 1000 * 24 * 30)

  if (
    dettesAnciennesParECN.partPatronale !== 0 ||
    dettesAnciennesParECN.partOuvriere !== 0 ||
    dettesRecentesParECN.partOuvriere !== 0 ||
    dettesRecentesParECN.partPatronale !== 0 ||
    // dettesAnciennesDebutParECN.partOuvriere !== 0 ||
    // dettesAnciennesDebutParECN.partPatronale !== 0 ||
    totalMoisDemande !== 0
  ) {
    emit(this.value.key, {
      montant_part_patronale_ancienne_courante:
        dettesAnciennesParECN.partPatronale,
      montant_part_ouvriere_ancienne_courante:
        dettesAnciennesParECN.partOuvriere,
      montant_part_patronale_recente_courante:
        dettesRecentesParECN.partPatronale,
      montant_part_ouvriere_recente_courante: dettesRecentesParECN.partOuvriere,
      // montant_part_ouvriere_ancienne_reference:
      //   dettesAnciennesDebutParECN.partOuvriere,
      // montant_part_patronale_ancienne_reference:
      //   dettesAnciennesDebutParECN.partPatronale,
      cotisation_moyenne_12m: cotisationMoyenne,
      total_demande_ap: totalMoisDemande,
    })
  }
}
