import { EntréeEffectif } from "../GeneratedTypes"
import { ParHash, ParPériode } from "../RawDataTypes"

// Paramètres globaux utilisés par "public"
declare const serie_periode: Date[]

export type SortieEffectif = {
  periode: Date
  effectif: number
}

export function effectifs(
  effectif?: ParHash<EntréeEffectif>
): SortieEffectif[] {
  const mapEffectif = new ParPériode<number>()
  Object.values(effectif ?? {}).forEach((e) => {
    mapEffectif.set(e.periode, (mapEffectif.get(e.periode) || 0) + e.effectif)
  })
  return serie_periode
    .map((p) => {
      return {
        periode: p,
        effectif: mapEffectif.get(p) || -1,
      }
    })
    .filter((p) => p.effectif >= 0)
}
