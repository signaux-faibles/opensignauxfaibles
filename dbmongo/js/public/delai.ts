import { iterable } from "./iterable"

export function delai<T>(delai?: Record<string, T>): T[] {
  const f = { iterable } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  return f.iterable(delai)
}
