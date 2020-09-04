import { ParPériode } from "../RawDataTypes"

type Outcome = {
  time_til_outcome: number
  outcome: boolean
}

export function lookAhead(
  data: ParPériode<Record<string, unknown>>,
  attr_name: string,
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

      if (data[period][attr_name]) {
        // si l'évènement se produit on retombe à 0
        counter = 0
      }

      if (counter >= 0) {
        // l'évènement s'est produit
        m[period] = m[period] || {}
        m[period].time_til_outcome = counter
        if (m[period].time_til_outcome <= n_months) {
          m[period].outcome = true
        } else {
          m[period].outcome = false
        }
      }
      return m
    }, {} as ParPériode<Outcome>)

  return output
}
