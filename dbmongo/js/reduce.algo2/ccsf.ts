type V = DonnéesCcsf

export type Output = {
  periode: Date
  date_ccsf: unknown
}

export function ccsf(v: V, output_array: Output[]): void {
  "use strict"

  const ccsfHashes = Object.keys(v.ccsf || {})

  output_array.forEach((val) => {
    const optccsf = ccsfHashes.reduce(
      function (accu, hash) {
        const ccsf = v.ccsf[hash]
        if (
          ccsf.date_traitement.getTime() < val.periode.getTime() &&
          ccsf.date_traitement.getTime() > accu.date_traitement.getTime()
        ) {
          accu = ccsf
        }
        return accu
      },
      {
        date_traitement: new Date(0),
      }
    )

    if (optccsf.date_traitement.getTime() != 0) {
      val.date_ccsf = optccsf.date_traitement
    }
  })
}
