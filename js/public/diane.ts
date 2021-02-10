import { EntréeDiane } from "../GeneratedTypes"

export function diane(hs?: Record<string, EntréeDiane>): EntréeDiane[] {
  "use strict"

  const diane: Record<string, EntréeDiane> = {}

  // Déduplication par arrete_bilan_diane
  Object.values(hs ?? {})
    .filter((d) => d.arrete_bilan_diane)
    .forEach((d) => {
      diane[d.arrete_bilan_diane.toISOString()] = d
    })

  return Object.values(diane ?? {}).sort((a, b) =>
    (a.exercice_diane ?? 0) < (b.exercice_diane ?? 0) ? 1 : -1
  )
}
