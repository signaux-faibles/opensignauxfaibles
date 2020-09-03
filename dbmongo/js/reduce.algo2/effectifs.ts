import * as f from "../common/dateAddMonth"

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

type EffectifEntreprise = Record<DataHash, EntréeEffectif>

export function effectifs(
  effobj: EffectifEntreprise,
  periodes: Timestamp[],
  propertyName: CléSortieEffectif
): ParPériode<SortieEffectifs> {
  "use strict"

  const sortieEffectif: ParPériode<SortieEffectifs> = {}

  // Construction d'une map[time] = effectif à cette periode
  const map_effectif = Object.keys(effobj).reduce((m, hash) => {
    const effectif = effobj[hash]
    if (effectif === null) {
      return m
    }
    const effectifTime = effectif.periode.getTime()
    m[effectifTime] = (m[effectifTime] || 0) + effectif.effectif
    return m
  }, {} as Record<Periode, number>)

  // Ne reporter que si le dernier effectif est disponible
  // On reporte dans les dernières périodes le dernier effectif connu
  const dernièrePériodeAvecEffectifConnu = f.dateAddMonth(
    new Date(periodes[periodes.length - 1]),
    offset_effectif + 1
  )
  const dernièrePériodeDisponible =
    dernièrePériodeAvecEffectifConnu.getTime() in map_effectif

  //pour chaque periode (elles sont triees dans l'ordre croissant)
  periodes.reduce((accu, time) => {
    // si disponible on reporte l'effectif tel quel, sinon, on recupère l'accu
    sortieEffectif[time] = sortieEffectif[time] || {}
    sortieEffectif[time][propertyName] =
      map_effectif[time] || (dernièrePériodeDisponible ? accu : null)

    // le cas échéant, on met à jour l'accu avec le dernier effectif disponible
    accu = map_effectif[time] || accu

    Object.assign(sortieEffectif[time], {
      [propertyName + "_reporte"]: map_effectif[time] ? 0 : 1,
    })

    return accu
  }, null as ValeurEffectif | null)

  Object.keys(map_effectif).forEach((time) => {
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
        [propertyName + "_past_" + offset]: map_effectif[time],
      }
    })
  })

  // On supprime les effectifs 'null'
  Object.keys(sortieEffectif).forEach((k) => {
    if (
      sortieEffectif[k].effectif === null &&
      sortieEffectif[k].effectif_ent === null
    ) {
      delete sortieEffectif[k]
    }
  })
  return sortieEffectif
}
