import { f } from "./functions"
import { EntréeApConso, ParHash } from "../RawDataTypes"

export function apconso(apconso?: ParHash<EntréeApConso>): EntréeApConso[] {
  return f
    .iterable(apconso)
    .sort((p1, p2) => (p1.periode < p2.periode ? 1 : -1))
}
