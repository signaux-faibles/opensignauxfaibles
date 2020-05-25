import "../globals.ts"

type Siret = string

type Scope = "etablissement" | "entreprise"

type ReduceIndexFlags = {
  algo1: boolean
  algo2: boolean
}

type Periode = Date

type RepOrder = {
  random_order: number
  periode: Periode
  siret: Siret
}

type BatchValue = {
  reporder: { [periode: string]: RepOrder }
}

type BatchValues = { [batchKey: string]: BatchValue }

type RawDataValues = {
  scope: Scope
  index: ReduceIndexFlags
  batch: BatchValues
}

// complete_reporder ajoute une propriété `reporder` pour chaque couple
// SIRET+période, afin d'assurer la reproductibilité de l'échantillonage.
export function complete_reporder(
  siret: Siret,
  object: RawDataValues
): RawDataValues {
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
        siret: siret,
      }
      object.batch[lastBatch].reporder = reporder_obj
    })
  return object
}
