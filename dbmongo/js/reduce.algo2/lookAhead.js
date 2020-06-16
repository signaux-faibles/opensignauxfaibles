function lookAhead(data, attr_name, n_months, past) {
  "use strict"
  // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false on
  // pourra se demander que s'est-il passé
  var sorting_fun = function (a, b) {
    return a >= b
  }
  if (past) {
    sorting_fun = function (a, b) {
      return a <= b
    }
  }

  var counter = -1
  var output = Object.keys(data)
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
    }, {})

  return output
}

exports.lookAhead = lookAhead
