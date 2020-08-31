type Input = {
  periode: Date
}

export type SortieCcsf = {
  date_ccsf: unknown
}

export function ccsf(
  vCcsf: Record<string, { date_traitement: Date }>,
  output_array: (Input & Partial<SortieCcsf>)[]
): void {
  "use strict"

  const ccsfHashes = Object.keys(vCcsf || {})

  output_array.forEach((val) => {
    const optccsf = ccsfHashes.reduce(
      function (accu, hash) {
        const ccsf = vCcsf[hash]
        if (
          ccsf.date_traitement.getTime() < val.periode.getTime() &&
          ccsf.date_traitement.getTime() > accu.date_traitement.getTime()
        ) {
          return ccsf
        }
        return accu
      },
      {
        date_traitement: new Date(0),
      }
    )

    if (optccsf.date_traitement.getTime() !== 0) {
      val.date_ccsf = optccsf.date_traitement
    }
  })
}
