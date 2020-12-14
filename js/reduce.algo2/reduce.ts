import { SortieMap } from "./map"

export function reduce(_key: unknown, values: SortieMap[]): SortieMap {
  "use strict"
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {} as SortieMap)
}
