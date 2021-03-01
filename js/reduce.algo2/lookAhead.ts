import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { SortieDefaillances } from "./defaillances"
import { SortieCotisation } from "./cotisation"

export type Outcome = {
  /** Distance de l'évènement, exprimé en nombre de périodes. */
  time_til_outcome: number
  outcome: boolean
}

type EntréeLookAhead = {
  outcome: { outcome: boolean }
  tag_default: Partial<Pick<SortieCotisation, "tag_default">>
  tag_failure: Partial<Pick<SortieDefaillances, "tag_failure">>
}

export function lookAhead<PropName extends keyof EntréeLookAhead>(
  data: ParPériode<EntréeLookAhead[PropName]>,
  attr_name: PropName /** "outcome" | "tag_default" | "tag_failure" */,
  n_months: number,
  past: boolean
): ParPériode<Outcome> {
  "use strict"
  // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false on
  // pourra se demander que s'est-il passé

  const chronologic = (a: number, b: number) => (a > b ? 1 : -1) // TODO: a - b
  const reverse = (a: number, b: number) => (b > a ? 1 : -1) // TODO: b - a

  let counter = -1
  const output = [...data.keys()]
    .sort(past ? reverse : chronologic)
    .reduce(function (m, période) {
      // Si on a déjà détecté quelque chose, on compte le nombre de périodes
      if (counter >= 0) counter = counter + 1

      // TODO: éviter l'explicitation de type ci-dessous:
      const dataInPeriod: Record<string, unknown> | undefined = data.get(
        période
      )
      if (dataInPeriod && dataInPeriod[attr_name]) {
        // si l'évènement se produit on retombe à 0
        counter = 0
      }

      if (counter >= 0) {
        // l'évènement s'est produit
        const out = m.get(période) ?? ({} as Outcome)
        out.time_til_outcome = counter
        if (out.time_til_outcome <= n_months) {
          out.outcome = true
        } else {
          out.outcome = false
        }
        m.set(période, out)
      }
      return m
    }, f.makePeriodeMap<Outcome>())

  return output
}
