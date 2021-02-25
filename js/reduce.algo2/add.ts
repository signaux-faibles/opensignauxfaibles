import { ParPériode } from "../RawDataTypes"

export function add<T>(
  obj: ParPériode<T>,
  output: ParPériode<Partial<T>>
): void {
  "use strict"
  output.forEach((val, période) => {
    if (obj.has(période)) {
      Object.assign(val, obj.get(période))
    }
  })
}
