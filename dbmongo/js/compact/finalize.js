function finalize(k, o) {
  o.index = {"algo1":false,
    "algo2":false}


  if (o.scope == "entreprise") {
    o.index.algo1 = true
    o.index.algo2 = true
  } else {
    // Est-ce que l'un des batchs a un effectif ?
    var batches = Object.keys(o.batch)
    batches.some(batch => {
      let hasEffectif = Object.keys(o.batch[batch].effectif || {}).length > 0
      o.index.algo1 = hasEffectif
      o.index.algo2 = hasEffectif
      return (hasEffectif)
    })

    // Complete reporder if missing

    var dates = serie_periode
    var missing = {}
    serie_periode.forEach(p => {
      missing[p] = true
    })

    batches.forEach(batch => {
      let reporder = o.batch[batch].reporder || {}

      Object.keys(reporder).forEach(ro => {
        missing[reporder[ro].periode] = false
      })
    })

    var lastBatch = batches[batches.length - 1]
    serie_periode.filter(p => missing[p]).forEach(p => {
      let reporder_obj = o.batch[lastBatch].reporder || {}
      reporder_obj[p] = { random_order: Math.random(), periode: p, siret: k }
      o.batch[lastBatch].reporder = reporder_obj
    })
  }
  return(o)
}
