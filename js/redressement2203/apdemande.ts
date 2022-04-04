import { EntréeApDemande } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"

export function apdemande(
  apdemande?: ParHash<EntréeApDemande>
): EntréeApDemande[] {
  return Object.values(apdemande ?? {}).sort((p1, p2) =>
    p1.periode.start.getTime() < p2.periode.start.getTime() ? 1 : -1
  )
}
