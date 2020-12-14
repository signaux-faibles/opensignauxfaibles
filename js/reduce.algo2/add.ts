import { ParPériode } from "../RawDataTypes"

export function add(
  obj: ParPériode<unknown>,
  output: ParPériode<unknown>
): void {
  "use strict"
  Object.keys(output).forEach(function (periode) {
    if (periode in obj) {
      Object.assign(output[periode], obj[periode])
    }
  })
}
