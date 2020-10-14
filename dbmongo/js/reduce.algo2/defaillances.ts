import { f } from "./functions"
import { SortieProcols } from "./dealWithProcols"
import { EntréeDefaillances, ParPériode, ParHash } from "../RawDataTypes"

export type SortieDefaillances = SortieProcols

export function defaillances(
  procol: ParHash<EntréeDefaillances>,
  output_indexed: ParPériode<Partial<SortieDefaillances>>
): void {
  "use strict"
  f.dealWithProcols(procol, output_indexed)
}
