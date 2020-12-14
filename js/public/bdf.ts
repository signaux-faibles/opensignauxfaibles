export type Bdf = { annee_bdf: number; arrete_bilan_bdf: Date }

export function bdf(hs?: Record<string | number, Bdf>): Bdf[] {
  "use strict"

  const bdf: Record<string, Bdf> = {}

  // Déduplication par arrete_bilan_bdf
  Object.values(hs ?? {})
    .filter((b) => b.arrete_bilan_bdf)
    .forEach((b) => {
      bdf[b.arrete_bilan_bdf.toISOString()] = b
    })

  return Object.values(bdf ?? {}).sort((a, b) =>
    a.annee_bdf < b.annee_bdf ? 1 : -1
  )
}
