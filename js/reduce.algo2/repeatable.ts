import { EntréeRepOrder, ParPériode } from "../RawDataTypes"

export type SortieRepeatable = { random_order: number }

export function repeatable(
  rep: ParPériode<EntréeRepOrder>
): ParPériode<SortieRepeatable> {
  "use strict"
  const output_repeatable: ParPériode<{ random_order: number }> = {}
  for (const one_rep of Object.values(rep)) {
    const periode = one_rep.periode.getTime()
    const out = output_repeatable[periode] ?? ({} as { random_order: number })
    out.random_order = one_rep.random_order
    output_repeatable[periode] = out
  }

  return output_repeatable
}
