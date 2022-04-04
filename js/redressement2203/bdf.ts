import { f } from "./functions"
import { DataHash } from "../RawDataTypes"

export type Bdf = { annee_bdf: number; arrete_bilan_bdf: Date }

export function bdf(hs?: Record<DataHash | number, Bdf>): Bdf[] {
  "use strict"

  const bdf = f.makePeriodeMap<Bdf>()

  // Déduplication par arrete_bilan_bdf
  for (const b of Object.values(hs ?? {})) {
    if (b.arrete_bilan_bdf !== undefined) {
      bdf.set(b.arrete_bilan_bdf, b)
    }
  }

  return [...bdf.values()].sort((a, b) => (a.annee_bdf < b.annee_bdf ? 1 : -1))
}
