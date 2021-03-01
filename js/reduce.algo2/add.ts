import { ParPériode } from "../common/makePeriodeMap"

export function add<T>(
  obj: ParPériode<T>,
  output: ParPériode<Partial<T>>
): void {
  "use strict"
  for (const période of output.keys()) {
    output.assign(période, obj.get(période))
  }
}
