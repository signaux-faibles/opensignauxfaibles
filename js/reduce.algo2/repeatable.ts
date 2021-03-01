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
  for (const { periode, random_order } of Object.values(rep)) {
    output_repeatable.assign(periode, { random_order })
  }
  return output_repeatable
}
