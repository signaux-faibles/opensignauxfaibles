import { ParPériode } from "../RawDataTypes"

export function add<T>(
  obj: ParPériode<T>,
  output: ParPériode<Partial<T>>
): void {
  "use strict"
  Object.keys(output).forEach(function (strPériode) {
    if (strPériode in obj) {
      const période = parseInt(strPériode)
      Object.assign(output[période], obj[période])
    }
  })
}
