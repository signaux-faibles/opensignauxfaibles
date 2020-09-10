import { iterable } from "./iterable"

export function apconso(
  apconso?: Record<DataHash, EntréeApConso>
): EntréeApConso[] {
  const f = { iterable } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  return f
    .iterable(apconso)
    .sort((p1, p2) => (p1.periode < p2.periode ? 1 : -1))
}
