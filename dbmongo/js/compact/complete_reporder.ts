import "../globals.ts"

// Pour rendre les dates au format 'Fri Apr 01 2016 00:00:00 GMT+0000 (UTC)'
const DATE_FORMAT = new Intl.DateTimeFormat("en-US", {
  weekday: "short",
  day: "2-digit",
  month: "short",
  year: "numeric",
  hour: "2-digit",
  hour12: false,
  minute: "2-digit",
  second: "2-digit",
  timeZone: "UTC",
})
const renderDate = (date: Date): string =>
  DATE_FORMAT.format(date).replace(/, /g, " ") + " GMT+0000 (UTC)"

// complete_reporder ajoute une propriété "reporder" pour chaque couple
// SIRET+période, afin d'assurer la reproductibilité de l'échantillonage.
export function complete_reporder(
  siret: SiretOrSiren,
  object: CompanyDataValuesWithFlags
): CompanyDataValuesWithFlags {
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
      reporder_obj[renderDate(p)] = {
        random_order: Math.random(),
        periode: p,
        siret: siret,
      }
      object.batch[lastBatch].reporder = reporder_obj
    })
  return object
}
