function repeatable(rep) {
  "use strict"
  const output_repeatable = {}
  Object.keys(rep).forEach((hash) => {
    const one_rep = rep[hash]
    const periode = one_rep.periode.getTime()
    output_repeatable[periode] = output_repeatable[periode] || {}
    output_repeatable[periode].random_order = one_rep.random_order
  })

  return output_repeatable
}

exports.repeatable = repeatable
