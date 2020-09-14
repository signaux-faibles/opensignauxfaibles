import { f } from "./functions"
import { EntréeEffectif, ParHash, Timestamp, ParPériode } from "../RawDataTypes"

// Paramètres globaux utilisés par "reduce.algo2"
declare const offset_effectif: number

type CléSortieEffectif = "effectif_ent" | "effectif" // effectif entreprise ou établissement
type CléSortieEffectifReporté = "effectif_ent_reporte" | "effectif_reporte"
type CléSortieEffectifPassé =
  | "effectif_past_6"
  | "effectif_past_12"
  | "effectif_past_18"
  | "effectif_past_24"
  | "effectif_ent_past_6"
  | "effectif_ent_past_12"
  | "effectif_ent_past_18"
  | "effectif_ent_past_24"

type ValeurEffectif = number

type SortieEffectifs = Record<CléSortieEffectif, ValeurEffectif | null> &
  Record<CléSortieEffectifReporté, 1 | 0> &
  Record<CléSortieEffectifPassé, ValeurEffectif>

type EffectifEntreprise = ParHash<EntréeEffectif>

export function effectifs(
  entréeEffectif: EffectifEntreprise,
  periodes: Timestamp[],
  clé: CléSortieEffectif
): ParPériode<SortieEffectifs> {
  "use strict"

  const sortieEffectif: ParPériode<SortieEffectifs> = {}

  // Construction d'une map[time] = effectif à cette periode
  const mapEffectif: ParPériode<ValeurEffectif> = {}

  Object.keys(entréeEffectif).forEach((hash) => {
    const effectif = entréeEffectif[hash]
    if (effectif !== null) {
      mapEffectif[effectif.periode.getTime()] = effectif.effectif
    }
  })

  // On reporte dans les dernières périodes le dernier effectif connu
  // Ne reporter que si le dernier effectif est disponible
  const dernièrePériodeAvecEffectifConnu = f.dateAddMonth(
    new Date(periodes[periodes.length - 1]),
    offset_effectif + 1
  )
  const effectifÀReporter =
    mapEffectif[dernièrePériodeAvecEffectifConnu.getTime()] ?? null

  periodes.forEach((time) => {
    sortieEffectif[time] = {
      ...sortieEffectif[time],
      [clé]: mapEffectif[time] || effectifÀReporter,
      [clé + "_reporte"]: mapEffectif[time] ? 0 : 1,
    }
  })

  Object.keys(mapEffectif).forEach((time) => {
    const futureTimestamps = [6, 12, 18, 24] // Penser à mettre à jour le type PastPropertyName pour tout changement
      .map((offset) => ({
        offset,
        timestamp: f
          .dateAddMonth(new Date(parseInt(time)), offset - offset_effectif - 1)
          .getTime(),
      }))
      .filter(({ timestamp }) => periodes.includes(timestamp))

    futureTimestamps.forEach(({ offset, timestamp }) => {
      sortieEffectif[timestamp] = {
        ...sortieEffectif[timestamp],
        [clé + "_past_" + offset]: mapEffectif[time],
      }
    })
  })
  return sortieEffectif
}
