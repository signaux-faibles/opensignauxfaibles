import "../globals.ts"

export function complete_reporder(key, object) {
  "use strict"
  const batches = Object.keys(object.batch)
  batches.sort()
  const missing = {}
  serie_periode.forEach((p) => {
    missing[p.getTime()] = true
  })

  batches.forEach((batch) => {
    const reporder = object.batch[batch].reporder || {}

    Object.keys(reporder).forEach((ro) => {
      if (!missing[reporder[ro].periode.getTime()]) {
        delete object.batch[batch].reporder[ro]
      } else {
        missing[reporder[ro].periode.getTime()] = false
      }
    })
  })

  const lastBatch = batches[batches.length - 1]
  serie_periode
    .filter((p) => missing[p.getTime()])
    .forEach((p) => {
      const reporder_obj = object.batch[lastBatch].reporder || {}
      reporder_obj[p.toString()] = {
        random_order: Math.random(),
        periode: p,
        siret: key,
      }
      object.batch[lastBatch].reporder = reporder_obj
    })
  return object
}
