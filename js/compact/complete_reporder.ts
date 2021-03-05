import {
  CompanyDataValuesWithFlags,
  SiretOrSiren,
  Timestamp,
} from "../RawDataTypes"

// Paramètres globaux utilisés par "compact"
declare const serie_periode: Date[]

// complete_reporder ajoute une propriété "reporder" pour chaque couple
// SIRET+période, afin d'assurer la reproductibilité de l'échantillonage.
export function complete_reporder(
  siret: SiretOrSiren,
  object: CompanyDataValuesWithFlags
): CompanyDataValuesWithFlags {
  "use strict"
  const batches = Object.keys(object.batch)
  batches.sort()
  const missing = new Set<Timestamp>(serie_periode.map((p) => p.getTime()))

  batches.forEach((batch) => {
    const reporder = object.batch[batch]?.reporder
    if (reporder === undefined) return
    Object.keys(reporder).forEach((ro) => {
      const periode = reporder[ro]?.periode
      if (periode === undefined) return
      if (!missing.has(periode.getTime())) {
        delete reporder[ro]
      } else {
        missing.delete(periode.getTime())
      }
    })
  })

  const lastBatch = batches[batches.length - 1]
  if (lastBatch === undefined)
    throw new Error("the last batch should not be undefined")
  serie_periode
    .filter((p) => missing.has(p.getTime()))
    .forEach((p) => {
      const dataInLastBatch = object.batch[lastBatch]
      if (dataInLastBatch === undefined) return
      const reporder_obj = dataInLastBatch.reporder ?? {}
      reporder_obj[p.toString()] = {
        random_order: Math.random(),
        periode: p,
        siret,
      }
      dataInLastBatch.reporder = reporder_obj
    })
  return object
}
