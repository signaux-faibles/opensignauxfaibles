import * as f from "./iterable"

export function diane(hs?: Record<string, EntréeDiane>): EntréeDiane[] {
  "use strict"
  return f
    .iterable<EntréeDiane>(hs)
    .sort((a, b) => (a.exercice_diane < b.exercice_diane ? 1 : -1))
}
