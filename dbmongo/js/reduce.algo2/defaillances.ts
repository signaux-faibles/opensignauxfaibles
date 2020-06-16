import * as f from "./dealWithProcols"
import { Event } from "./dealWithProcols"

export function defaillances(
  v: {
    altares: { [hash: string]: Event }
    procol: { [hash: string]: Event }
  },
  output_indexed: {}
): void {
  "use strict"
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
