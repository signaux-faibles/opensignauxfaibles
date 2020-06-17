import * as f from "./dealWithProcols"
import { InputEvent } from "./dealWithProcols"

export function defaillances(
  v: {
    altares: { [hash: string]: InputEvent }
    procol: { [hash: string]: InputEvent }
  },
  output_indexed: {}
): void {
  "use strict"
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
