import { f } from "./functions"
import { ParPériode } from "../RawDataTypes"

type Times = {
  time_til_default?: number
  time_til_failure?: number
}

export function cibleApprentissage(
  output_indexed: ParPériode<{ tag_failure?: boolean; tag_default?: boolean }>,
  n_months: number
): ParPériode<Partial<Times>> {
  "use strict"

  // Mock two input instead of one for future modification
  const output_cotisation = output_indexed
  const output_procol = output_indexed
  // replace with const
  const all_keys = Object.keys(output_indexed)

  const merged_info: ParPériode<{ outcome: boolean }> = {}
  for (const k of all_keys) {
    merged_info[k] = {
      outcome: Boolean(
        output_procol[k]?.tag_failure || output_cotisation[k]?.tag_default
      ),
    }
  }

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
    if (output_default[k] !== undefined)
      outputTimes.time_til_default = output_default[k]?.time_til_outcome
    if (output_failure[k] !== undefined)
      outputTimes.time_til_failure = output_failure[k]?.time_til_outcome
    return {
      ...m,
      [k]: {
        ...output_outcome[k],
        ...outputTimes,
      },
    }
  }, {} as ParPériode<Times>)

  return output_cible
}
