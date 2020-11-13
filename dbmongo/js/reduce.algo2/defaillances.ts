import { f } from "./functions"
import { SortieProcols } from "./dealWithProcols"
import { EntréeDéfaillances, ParPériode, ParHash } from "../RawDataTypes"

export type SortieDefaillances = SortieProcols

export function defaillances(
  procol: ParHash<EntréeDéfaillances>,
  output_indexed: ParPériode<Partial<SortieDefaillances>>
): void {
  "use strict"
  f.dealWithProcols(procol, output_indexed)
}
