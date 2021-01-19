import { ParHash } from "../RawDataTypes"

type Input = {
  periode: Date
}

// VariablesSource est utilisé pour populer `source` dans docs/variables.json (cf generate-docs.ts)
export type VariablesSource = "ccsf"

// ComputedVariables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type ComputedVariables = {
  date_ccsf: Date
}

// TransmittedVariables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
// export type TransmittedVariables = {}

export type SortieCcsf = ComputedVariables // & TransmittedVariables

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
