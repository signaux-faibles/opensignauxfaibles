import { V } from "./finalize"

export function reduce(_key: unknown, values: V[]): V {
  "use strict"
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {} as V)
}
