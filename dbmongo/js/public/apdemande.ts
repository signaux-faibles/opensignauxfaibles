import { EntréeApDemande, ParHash } from "../RawDataTypes"
import { iterable } from "./iterable"

export function apdemande(
  apdemande?: ParHash<EntréeApDemande>
): EntréeApDemande[] {
  const f = { iterable } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  return f
    .iterable(apdemande)
    .sort((p1, p2) =>
      p1.periode.start.getTime() < p2.periode.start.getTime() ? 1 : -1
    )
}
