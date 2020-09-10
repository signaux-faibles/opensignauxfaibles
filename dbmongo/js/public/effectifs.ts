import { iterable } from "./iterable"
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
  const f = { iterable } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const mapEffectif: ParPériode<number> = {}
  f.iterable(effectif).forEach((e) => {
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
