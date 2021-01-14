import { ParPériode } from "../RawDataTypes"

export function add<T>(
  obj: ParPériode<T>,
  output: ParPériode<Partial<T>>
): void {
  "use strict"
  Object.keys(output).forEach(function (periode) {
    if (periode in obj) {
      Object.assign(output[periode], obj[periode])
    }
  })
}
