import * as f from "./iterable"

type Bdf = { annee_bdf: number }

export function bdf(hs?: Record<string | number, Bdf>): Bdf[] {
  "use strict"
  return f
    .iterable<Bdf>(hs)
    .sort((a, b) => (a.annee_bdf < b.annee_bdf ? 1 : -1))
}
