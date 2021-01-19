import { ParHash } from "../RawDataTypes"

type Input = {
  periode: Date
}

export type SortieCcsf = {
  /** Date de début de la procédure CCSF */
  date_ccsf: Date
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "ccsf"
  computed: SortieCcsf
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function ccsf(
  vCcsf: ParHash<{ date_traitement: Date }>,
  output_array: (Input & Partial<SortieCcsf>)[]
): void {
  "use strict"

  output_array.forEach((val) => {
    let optccsfDateTraitement = new Date(0)
    for (const ccsf of Object.values(vCcsf)) {
      if (
        ccsf.date_traitement.getTime() < val.periode.getTime() &&
        ccsf.date_traitement.getTime() > optccsfDateTraitement.getTime()
      ) {
        optccsfDateTraitement = ccsf.date_traitement
      }
    }

    if (optccsfDateTraitement.getTime() !== 0) {
      val.date_ccsf = optccsfDateTraitement
    }
  })
}
