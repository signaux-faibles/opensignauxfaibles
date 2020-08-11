import * as f from "./iterable"

// Paramètres globaux utilisés par "public"
declare const serie_periode: Date[]

type Result = {
  periode: Date
  effectif: number | null
}

export function effectifs(v: {
  effectif: Record<DataHash, EntréeEffectif>
}): Result[] {
  const mapEffectif: ParPériode<number> = {}
  f.iterable(v.effectif).forEach((e) => {
    mapEffectif[e.periode.getTime()] =
      (mapEffectif[e.periode.getTime()] || 0) + e.effectif
  })
  return serie_periode
    .map((p) => {
      return {
        periode: p,
        effectif: mapEffectif[p.getTime()] || null,
      }
    })
    .filter((p) => p.effectif)
}
