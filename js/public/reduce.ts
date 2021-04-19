import { SortieMap } from "./map"

export type V = Partial<SortieMap>

export function reduce(_key: unknown, values: V[]): V {
  return values.reduce((m, v) => {
    Object.assign(m, v)
    return m
  }, {} as V)
}
