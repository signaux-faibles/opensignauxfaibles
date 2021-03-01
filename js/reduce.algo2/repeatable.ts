import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { EntréeRepOrder, ParHash } from "../RawDataTypes"

export type SortieRepeatable = {
  /** Numéro permettant de réaliser un échantillon reproductible des données. */
  random_order: number
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "repeatable"
  computed: SortieRepeatable
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function repeatable(
  rep: ParHash<EntréeRepOrder>
): ParPériode<SortieRepeatable> {
  "use strict"
  const output_repeatable = f.makePeriodeMap<SortieRepeatable>()
  for (const one_rep of Object.values(rep)) {
    const periode = one_rep.periode.getTime()
    const out = output_repeatable.get(periode) ?? ({} as SortieRepeatable)
    out.random_order = one_rep.random_order
    output_repeatable.set(periode, out) // TODO: utiliser append() ou upsert()
  }

  return output_repeatable
}
