import * as f from "./dealWithProcols"
import { Output as ProcolOutput } from "./dealWithProcols"

type V = DonnéesDefaillances

export type Output = ProcolOutput

export function defaillances(
  v: V,
  output_indexed: Record<Periode, Partial<Output>>
): void {
  "use strict"
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
