import * as f from "./iterable"

export function apconso(
  apconso?: Record<DataHash, EntréeApConso>
): EntréeApConso[] {
  return f
    .iterable(apconso)
    .sort((p1, p2) => (p1.periode < p2.periode ? 1 : -1))
}