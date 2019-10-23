function complete_reporder(key, object){
  var batches = Object.keys(object.batch)
  batches.sort()
  var dates = serie_periode
  var missing = {}
  serie_periode.forEach(p => {
    missing[p] = true
  })
  print(JSON.stringify(object, null, 2))


  batches.forEach(batch => {
    let reporder = object.batch[batch].reporder || {}

    Object.keys(reporder).forEach(ro => {
      missing[reporder[ro].periode] = false
    })
  })

  var lastBatch = batches[batches.length - 1]
  serie_periode.filter(p => missing[p]).forEach(p => {
    let reporder_obj = object.batch[lastBatch].reporder || {}
    reporder_obj[p] = { random_order: Math.random(), periode: p, siret: key }
    object.batch[lastBatch].reporder = reporder_obj
  })
  return(object)
}
