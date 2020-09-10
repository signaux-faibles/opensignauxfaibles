import { dealWithProcols } from "./dealWithProcols"
import { SortieProcols } from "./dealWithProcols"

export type SortieDefaillances = SortieProcols

export function defaillances(
  altares: Record<string, EntréeDefaillances>,
  procol: Record<string, EntréeDefaillances>,
  output_indexed: Record<Periode, Partial<SortieDefaillances>>
): void {
  "use strict"
  const f = { dealWithProcols } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  f.dealWithProcols(altares, "altares", output_indexed)
  f.dealWithProcols(procol, "procol", output_indexed)
}
