import { f } from "./functions"
import { ParPériode } from "../RawDataTypes"
import { Outcome } from "./lookAhead"

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
  n_months: number
): ParPériode<SortieCibleApprentissage> {
  "use strict"

  // Mock two input instead of one for future modification
  const output_cotisation = output_indexed
  const output_procol = output_indexed
  // replace with const
  const strPériodes = Object.keys(output_indexed)

  const merged_info: ParPériode<{ outcome: boolean }> = {}
  for (const strPériode of strPériodes) {
    const période = parseInt(strPériode)
    merged_info[période] = {
      outcome: Boolean(
        output_procol[période]?.tag_failure ||
          output_cotisation[période]?.tag_default
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

  const output_cible = strPériodes.reduce(function (m, strPériode) {
    const k = parseInt(strPériode)
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
