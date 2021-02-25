import { ParPériode } from "../common/ParPériode"

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
