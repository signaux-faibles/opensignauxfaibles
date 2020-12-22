import { ParPériode } from "../RawDataTypes"

export function add<T>(obj: ParPériode<Partial<T>>, output: ParPériode<T>) {
  "use strict"
  Object.keys(output).forEach(function (periode) {
    if (periode in obj) {
      Object.assign(output[periode], obj[periode])
    }
  })
}
