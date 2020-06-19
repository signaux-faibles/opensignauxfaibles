import * as f from "./dateAddMonth"

declare const offset_effectif: number

type Time = string

type EffectifName = string

type ValeurEffectif = number

type Output = Record<Time, Record<EffectifName, ValeurEffectif | null>>

export function effectifs(
  effobj: EffectifEntreprise,
  periodes: string[],
  effectif_name: EffectifName
): Output {
  "use strict"

  const output_effectif: Output = {}

  // Construction d'une map[time] = effectif à cette periode
  const map_effectif = Object.keys(effobj).reduce((m, hash) => {
    const effectif = effobj[hash]
    if (effectif == null) {
      return m
    }
    const effectifTime = effectif.periode.getTime()
    m[effectifTime] = (m[effectifTime] || 0) + effectif.effectif
    return m
  }, {} as Record<Time, ValeurEffectif>)

  //ne reporter que si le dernier est disponible
  // 1- quelle periode doit être disponible
  const last_period = new Date(parseInt(periodes[periodes.length - 1]))
  const last_period_offset = f.dateAddMonth(last_period, offset_effectif + 1)
  // 2- Cette période est-elle disponible ?

  const available = map_effectif[last_period_offset.getTime()] ? 1 : 0

  //pour chaque periode (elles sont triees dans l'ordre croissant)
  periodes.reduce((accu, time) => {
    // si disponible on reporte l'effectif tel quel, sinon, on recupère l'accu
    output_effectif[time] = output_effectif[time] || {}
    output_effectif[time][effectif_name] =
      map_effectif[time] || (available ? accu : null)

    // le cas échéant, on met à jour l'accu avec le dernier effectif disponible
    accu = map_effectif[time] || accu

    output_effectif[time][effectif_name + "_reporte"] = map_effectif[time]
      ? 0
      : 1
    return accu
  }, null as ValeurEffectif | null)

  Object.keys(map_effectif).forEach((time) => {
    const periode = new Date(parseInt(time))
    const past_month_offsets = [6, 12, 18, 24]
    past_month_offsets.forEach((lookback) => {
      // On ajoute un offset pour partir de la dernière période où l'effectif est connu
      const time_past_lookback = f.dateAddMonth(
        periode,
        lookback - offset_effectif - 1
      )

      const variable_name_effectif = effectif_name + "_past_" + lookback
      output_effectif[time_past_lookback.getTime()] =
        output_effectif[time_past_lookback.getTime()] || {}
      output_effectif[time_past_lookback.getTime()][variable_name_effectif] =
        map_effectif[time]
    })
  })

  // On supprime les effectifs 'null'
  Object.keys(output_effectif).forEach((k) => {
    if (
      output_effectif[k].effectif == null &&
      output_effectif[k].effectif_ent == null
    ) {
      delete output_effectif[k]
    }
  })
  return output_effectif
}
