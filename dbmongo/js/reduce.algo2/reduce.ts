import { EntréeFinalize } from "./finalize"

export function reduce(
  _key: unknown,
  values: EntréeFinalize[]
): EntréeFinalize {
  "use strict"
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {} as EntréeFinalize)
}
