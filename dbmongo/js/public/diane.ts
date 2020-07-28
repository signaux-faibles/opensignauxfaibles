import * as f from "./iterable"

type Diane = { exercice_diane: number }

export function diane(hs: Record<string | number, Diane>): Diane[] {
  "use strict"
  return f
    .iterable<Diane>(hs)
    .sort((a, b) => (a.exercice_diane < b.exercice_diane ? 1 : -1))
}
