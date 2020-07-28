import { V } from "./finalize"

export function reduce<T>(_key: unknown, values: T[]): V {
  "use strict"
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {} as V)
}
