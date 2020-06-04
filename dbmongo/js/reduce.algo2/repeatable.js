function repeatable(rep){
  "use strict";
  let output_repeatable = {}
  Object.keys(rep).forEach(hash => {
    var one_rep = rep[hash]
    var periode = one_rep.periode.getTime()
    output_repeatable[periode] = output_repeatable[periode] || {}
    output_repeatable[periode].random_order = one_rep.random_order
  })

  return(output_repeatable)

}

exports.repeatable = repeatable
