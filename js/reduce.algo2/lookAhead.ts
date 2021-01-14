import { SortieDefaillances } from "./defaillances"
import { ParPériode } from "../RawDataTypes"
import { SortieCotisation } from "./cotisation"

type Outcome = {
  time_til_outcome: number
  outcome: boolean
}

type EntréeLookAhead = {
  outcome: { outcome: boolean }
  tag_default: Partial<Pick<SortieCotisation, "tag_default">>
  tag_failure: Partial<Pick<SortieDefaillances, "tag_failure">>
}

export function lookAhead<
  K extends keyof EntréeLookAhead,
  T extends EntréeLookAhead[K]
>(
  data: ParPériode<T>,
  attr_name: K,
  n_months: number,
  past: boolean
): ParPériode<Outcome> {
  "use strict"
  // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false on
  // pourra se demander que s'est-il passé

  const chronologic = (a: string, b: string) => (a > b ? 1 : -1)
  const reverse = (a: string, b: string) => (b > a ? 1 : -1)

  let counter = -1
  const output = Object.keys(data)
    .sort(past ? reverse : chronologic)
    .reduce(function (m, period) {
      // Si on a déjà détecté quelque chose, on compte le nombre de périodes
      if (counter >= 0) counter = counter + 1

      const dataInPeriod: Record<string, unknown> | undefined = data[period]
      if (dataInPeriod && dataInPeriod[attr_name]) {
        // si l'évènement se produit on retombe à 0
        counter = 0
      }

      if (counter >= 0) {
        // l'évènement s'est produit
        const out = m[period] ?? ({} as Outcome)
        out.time_til_outcome = counter
        if (out.time_til_outcome <= n_months) {
          out.outcome = true
        } else {
          out.outcome = false
        }
        m[period] = out
      }
      return m
    }, {} as ParPériode<Outcome>)

  return output
}
