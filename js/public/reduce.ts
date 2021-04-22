import { SortieMap } from "./map"

export type V = Partial<SortieMap>

export function reduce(_key: unknown, values: V[]): V {
  return Object.assign({}, ...values)
}
