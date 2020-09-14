import { f } from "./functions"
import { SortieProcols } from "./dealWithProcols"
import { EntréeDefaillances, ParPériode, ParHash } from "../RawDataTypes"

export type SortieDefaillances = SortieProcols

export function defaillances(
  altares: ParHash<EntréeDefaillances>,
  procol: ParHash<EntréeDefaillances>,
  output_indexed: ParPériode<Partial<SortieDefaillances>>
): void {
  "use strict"
  f.dealWithProcols(altares, "altares", output_indexed)
  f.dealWithProcols(procol, "procol", output_indexed)
}
