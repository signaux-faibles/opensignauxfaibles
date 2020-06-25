import * as f from "./dealWithProcols"
import { Output as ProcolOutput } from "./dealWithProcols"

export type Output = ProcolOutput

export function defaillances(
  v: DonnéesDefaillances,
  output_indexed: Record<Periode, Partial<Output>>
): void {
  "use strict"
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
