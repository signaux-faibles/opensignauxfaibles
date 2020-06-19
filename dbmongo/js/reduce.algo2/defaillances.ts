import * as f from "./dealWithProcols"
import { InputEvent, Output } from "./dealWithProcols"

export function defaillances(
  v: {
    altares: { [hash: string]: InputEvent }
    procol: { [hash: string]: InputEvent }
  },
  output_indexed: { [time: string]: Output }
): void {
  "use strict"
  f.dealWithProcols(v.altares, "altares", output_indexed)
  f.dealWithProcols(v.procol, "procol", output_indexed)
}
