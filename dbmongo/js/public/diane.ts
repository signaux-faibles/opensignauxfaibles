import * as f from "./iterable"

export function diane(hs?: Record<string, EntréeDiane>): EntréeDiane[] {
  "use strict"

  const diane: Record<string, EntréeDiane> = {}

  f.iterable(hs)
    .filter((d) => d.arrete_bilan_diane)
    .forEach((d) => {
      diane[d.arrete_bilan_diane.toISOString()] = d
    })

  return f
    .iterable(diane)
    .sort((a, b) => (a.exercice_diane < b.exercice_diane ? 1 : -1))
}
