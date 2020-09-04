import "../globals"
import { lookAhead } from "./lookAhead"

type Times = {
  time_til_default?: number
  time_til_failure?: number
}

export function cibleApprentissage(
  output_indexed: {
    [k: string]: { tag_failure?: boolean; tag_default?: boolean }
  },
  n_months: number
): { [k: string]: Partial<Times> } {
  "use strict"
  const f = { lookAhead } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  // Mock two input instead of one for future modification
  const output_cotisation = output_indexed
  const output_procol = output_indexed
  // replace with const
  const all_keys = Object.keys(output_indexed)

  const merged_info = all_keys.reduce(function (m, k) {
    m[k] = {
      outcome: Boolean(
        output_procol[k].tag_failure || output_cotisation[k].tag_default
      ),
    }
    return m
  }, {} as Record<Periode, { outcome: boolean }>)

  const output_outcome = f.lookAhead(merged_info, "outcome", n_months, true)
  const output_default = f.lookAhead(
    output_cotisation,
    "tag_default",
    n_months,
    true
  )
  const output_failure = f.lookAhead(
    output_procol,
    "tag_failure",
    n_months,
    true
  )

  const output_cible = all_keys.reduce(function (m, k) {
    const outputTimes: Times = {}
    if (output_default[k])
      outputTimes.time_til_default = output_default[k].time_til_outcome
    if (output_failure[k])
      outputTimes.time_til_failure = output_failure[k].time_til_outcome
    return {
      ...m,
      [k]: {
        ...output_outcome[k],
        ...outputTimes,
      },
    }
  }, {} as Record<string, Times>)

  return output_cible
}
