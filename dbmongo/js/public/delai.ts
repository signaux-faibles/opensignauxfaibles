import * as f from "./iterable"

export function delai<T>(delai?: Record<string, T>): T[] {
  return f.iterable(delai)
}
