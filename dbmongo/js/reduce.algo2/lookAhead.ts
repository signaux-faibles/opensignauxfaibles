type Outcome = {
  time_til_outcome: number
  outcome: boolean
}

export function lookAhead(
  data: { [period: string]: Record<string, unknown> },
  attr_name: string,
  n_months: number,
  past: boolean
): { [period: string]: Outcome } {
  "use strict"
  // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false on
  // pourra se demander que s'est-il passé

  /* eslint-disable */
  var sorting_fun = function (a: any, b: any): any {
    return a >= b ? 1 : -1 // TODO: normally, a sorting comparator should return a number, possibly including zero. => the TS version of the test has failed until we added `? 1 : -1` here
  }
  if (past) {
    sorting_fun = function (a: any, b: any): any {
      return a <= b ? 1 : -1 // TODO: normally, a sorting comparator should return a number, possibly including zero. => the TS version of the test has failed until we added `? 1 : -1` here
    }
  }
  /* eslint-enable */

  let counter = -1
  const output = Object.keys(data)
    .sort(sorting_fun)
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
    }, {} as Record<string, Outcome>)

  return output
}
