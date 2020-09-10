import { iterable } from "./iterable"

export function diane(hs?: Record<string, EntréeDiane>): EntréeDiane[] {
  "use strict"
  const f = { iterable } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const diane: Record<string, EntréeDiane> = {}

  // Déduplication par arrete_bilan_diane
  f.iterable(hs)
    .filter((d) => d.arrete_bilan_diane)
    .forEach((d) => {
      diane[d.arrete_bilan_diane.toISOString()] = d
    })

  return f
    .iterable(diane)
    .sort((a, b) => (a.exercice_diane < b.exercice_diane ? 1 : -1))
}
