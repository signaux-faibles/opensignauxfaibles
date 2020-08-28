import * as f from "../common/dateAddMonth"

// Paramètres globaux utilisés par "reduce.algo2"
declare const offset_effectif: number

type PropertyName = "effectif_ent" | "effectif" // effectif entreprise ou établissement
type PropertyNameReporté = "effectif_ent_reporte" | "effectif_reporte"
type PastPropertyName =
  | "effectif_past_6"
  | "effectif_past_12"
  | "effectif_past_18"
  | "effectif_past_24"
  | "effectif_ent_past_6"
  | "effectif_ent_past_12"
  | "effectif_ent_past_18"
  | "effectif_ent_past_24"

type ValeurEffectif = number

type SortieEffectifs = {
  [propName in PropertyName]: ValeurEffectif | null
} &
  {
    [propName in PropertyNameReporté]: 1 | 0
  } &
  {
    [propName in PastPropertyName]: ValeurEffectif
  }

type EffectifEntreprise = Record<DataHash, EntréeEffectif>

export function effectifs(
  effobj: EffectifEntreprise,
  periodes: Timestamp[],
  propertyName: PropertyName
): ParPériode<SortieEffectifs> {
  "use strict"

  const output_effectif: ParPériode<SortieEffectifs> = {}

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

  //ne reporter que si le dernier est disponible
  // 1- quelle periode doit être disponible
  const last_period = new Date(periodes[periodes.length - 1])
  const last_period_offset = f.dateAddMonth(last_period, offset_effectif + 1)
  // 2- Cette période est-elle disponible ?

  const available = last_period_offset.getTime() in map_effectif

  //pour chaque periode (elles sont triees dans l'ordre croissant)
  periodes.reduce((accu, time) => {
    // si disponible on reporte l'effectif tel quel, sinon, on recupère l'accu
    output_effectif[time] = output_effectif[time] || {}
    output_effectif[time][propertyName] =
      map_effectif[time] || (available ? accu : null)

    // le cas échéant, on met à jour l'accu avec le dernier effectif disponible
    accu = map_effectif[time] || accu

    Object.assign(output_effectif[time], {
      [propertyName + "_reporte"]: map_effectif[time] ? 0 : 1,
    })

    return accu
  }, null as ValeurEffectif | null)

  Object.keys(map_effectif).forEach((time) => {
    const periode = new Date(parseInt(time))
    const past_month_offsets = [6, 12, 18, 24] // Note: à garder en synchro avec la définition du type PastPropertyName
    past_month_offsets.forEach((lookback) => {
      // On ajoute un offset pour partir de la dernière période où l'effectif est connu
      const time_past_lookback = f.dateAddMonth(
        periode,
        lookback - offset_effectif - 1
      )

      output_effectif[time_past_lookback.getTime()] =
        output_effectif[time_past_lookback.getTime()] || {}
      Object.assign(output_effectif[time_past_lookback.getTime()], {
        [propertyName + "_past_" + lookback]: map_effectif[time],
      })
    })
  })

  // On supprime les effectifs 'null'
  Object.keys(output_effectif).forEach((k) => {
    if (
      output_effectif[k].effectif === null &&
      output_effectif[k].effectif_ent === null
    ) {
      delete output_effectif[k]
    }
  })
  return output_effectif
}

/* TODO: appliquer même logique d'itération sur futureTimestamps que dans cotisationsdettes.ts */
