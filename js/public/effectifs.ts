import { EntréeEffectif, ParHash, ParPériode } from "../RawDataTypes"

// Paramètres globaux utilisés par "public"
declare const serie_periode: Date[]

export type SortieEffectif = {
  periode: Date
  effectif: number
}

export function effectifs(
  effectif?: ParHash<EntréeEffectif>
): SortieEffectif[] {
  const mapEffectif: ParPériode<number> = {}
  Object.values(effectif ?? {}).forEach((e) => {
    mapEffectif[e.periode.getTime()] =
      (mapEffectif[e.periode.getTime()] || 0) + e.effectif
  })
  return serie_periode
    .map((p) => {
      return {
        periode: p,
        effectif: mapEffectif[p.getTime()] || -1,
      }
    })
    .filter((p) => p.effectif >= 0)
}