import * as f from "./dealWithProcols"
import { SortieProcols } from "./dealWithProcols"

export type SortieDefaillances = SortieProcols

export function defaillances(
  v: DonnéesDefaillances,
  output_indexed: Record<Periode, Partial<SortieDefaillances>>
): void {
  "use strict"
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
