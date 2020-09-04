import { EntréeRepOrder, Periode } from "../RawDataTypes"

export type SortieRepeatable = { random_order: number }

export function repeatable(
  rep: Record<Periode, EntréeRepOrder>
): Record<Periode, SortieRepeatable> {
  "use strict"
  const output_repeatable: Record<string, { random_order: number }> = {}
  Object.keys(rep).forEach((hash) => {
    const one_rep = rep[hash]
    const periode = one_rep.periode.getTime()
    output_repeatable[periode] = output_repeatable[periode] || {}
    output_repeatable[periode].random_order = one_rep.random_order
  })

  return output_repeatable
}
