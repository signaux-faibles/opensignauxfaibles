import { f } from "./functions"
import { DataHash } from "../RawDataTypes"
import { EntréeDiane } from "../GeneratedTypes"

export function diane(hs?: Record<DataHash, EntréeDiane>): EntréeDiane[] {
  "use strict"

  const diane = f.makePeriodeMap<EntréeDiane>()

  // Déduplication par arrete_bilan_diane
  for (const d of Object.values(hs ?? {})) {
    if (d.arrete_bilan_diane !== undefined) {
      diane.set(d.arrete_bilan_diane, d)
    }
  }

  return [...diane.values()].sort((a, b) =>
    (a.exercice_diane ?? 0) < (b.exercice_diane ?? 0) ? 1 : -1
  )
}
