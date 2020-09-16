import { f } from "./functions"

export type Bdf = { annee_bdf: number; arrete_bilan_bdf: Date }

export function bdf(hs?: Record<string | number, Bdf>): Bdf[] {
  "use strict"

  const bdf: Record<string, Bdf> = {}

  // Déduplication par arrete_bilan_bdf
  f.iterable<Bdf>(hs)
    .filter((b) => b.arrete_bilan_bdf)
    .forEach((b) => {
      bdf[b.arrete_bilan_bdf.toISOString()] = b
    })

  return f
    .iterable<Bdf>(bdf)
    .sort((a, b) => (a.annee_bdf < b.annee_bdf ? 1 : -1))
}
