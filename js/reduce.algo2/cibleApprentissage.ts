import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { Outcome } from "./lookAhead"
import { Timestamp } from "../RawDataTypes"

export type SortieCibleApprentissage = {
  outcome?: Outcome["outcome"]
  /** Distance de l'évènement, exprimé en nombre de périodes. */
  time_til_outcome?: Outcome["time_til_outcome"]
  /** Distance de l'évènement basé sur le défaut de paiement des cotisations (cf tag_default), exprimé en nombre de périodes. */
  time_til_default?: number
  /** Distance de l'évènement basé sur une défaillance (cf tag_failure des procédures collectives), exprimé en nombre de périodes. */
  time_til_failure?: number
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "cibleApprentissage"
  computed: SortieCibleApprentissage
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function cibleApprentissage(
  output_indexed: ParPériode<{ tag_failure?: boolean; tag_default?: boolean }>,
  n_months: number // nombre de mois avant/après l'évènement pendant lesquels outcome sera true
): ParPériode<SortieCibleApprentissage> {
  "use strict"

  // Mock two input instead of one for future modification
  const output_cotisation = output_indexed
  const output_procol = output_indexed
  // replace with const
  const périodes = [...output_indexed.keys()]

  const merged_info = f.makePeriodeMap<{ outcome: boolean }>()
  for (const période of périodes) {
    merged_info.set(période, {
      outcome: Boolean(
        output_procol.get(période)?.tag_failure ||
          output_cotisation.get(période)?.tag_default
      ),
    })
  }

  function objectMap<InputVal, OutputVal>(
    input: ParPériode<InputVal>,
    fct: (key: Timestamp, val: InputVal) => OutputVal
  ): ParPériode<OutputVal> {
    const result = f.makePeriodeMap<OutputVal>()
    input.forEach((val, key) => {
      result.set(key, fct(key, val))
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

  const output_cible = périodes.reduce(function (m, k) {
    const oDefault = output_default.get(k)
    const oFailure = output_failure.get(k)
    return m.set(k, {
      ...outputPastOutcome.get(k),
      ...output_outcome.get(k),
      ...(oDefault && { time_til_default: oDefault.time_til_outcome }),
      ...(oFailure && { time_til_failure: oFailure.time_til_outcome }),
    })
  }, f.makePeriodeMap<SortieCibleApprentissage>())

  return output_cible
}
