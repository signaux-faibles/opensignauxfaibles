import { f } from "./functions"
import { ParPériode } from "../RawDataTypes"
import { SortieDefaillances } from "./defaillances"
import { SortieCotisation } from "./cotisation"
import { Outcome } from "./lookAhead"

export type SortieCibleApprentissage = {
  outcome?: Outcome["outcome"]
  /** Distance de l'évènement, exprimé en nombre de périodes. */
  time_til_outcome?: Outcome["time_til_outcome"]
  /** Distance de l'évènement basé sur le défaut de paiement des cotisations (cf tag_default), exprimé en nombre de périodes. */
  time_til_default?: Outcome["time_til_outcome"]
  /** Distance de l'évènement basé sur une défaillance (cf tag_failure des procédures collectives), exprimé en nombre de périodes. */
  time_til_failure?: Outcome["time_til_outcome"]
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "cibleApprentissage"
  computed: SortieCibleApprentissage
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function cibleApprentissage(
  output_indexed: ParPériode<{
    tag_failure?: SortieDefaillances["tag_failure"]
    tag_default?: SortieCotisation["tag_default"]
  }>,
  n_months: number /** nombre de mois avant/après l'évènement pendant lesquels outcome sera true */
): ParPériode<SortieCibleApprentissage> {
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

  function objectMap<InputVal, OutputVal>(
    obj: Record<string, InputVal>,
    fct: (key: string, val: InputVal) => OutputVal
  ): Record<string, OutputVal> {
    const result: Record<string, OutputVal> = {}
    Object.entries(obj).forEach(([key, val]) => {
      result[key] = fct(key, val)
    })
    return result
  }

  const outputPastOutcome = objectMap(
    f.lookAhead(merged_info, "outcome", n_months, false),
    (_, val) => ({
      ...val,
      time_til_outcome: -val.time_til_outcome, // ex: -1 veut dire qu'il y a eu une défaillance il y a 1 mois
    })
  )
  const output_outcome = {
    ...outputPastOutcome,
    ...f.lookAhead(merged_info, "outcome", n_months, true),
  }
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
    const outputTimes: SortieCibleApprentissage = {}
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
  }, {} as ParPériode<SortieCibleApprentissage>)

  return output_cible
}
