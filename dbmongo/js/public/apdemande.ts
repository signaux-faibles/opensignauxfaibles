import * as f from "./iterable"

export function apdemande(
  apdemande?: Record<DataHash, EntréeApDemande>
): EntréeApDemande[] {
  return f
    .iterable(apdemande)
    .sort((p1, p2) => (p1.periode < p2.periode ? 1 : -1))
}
