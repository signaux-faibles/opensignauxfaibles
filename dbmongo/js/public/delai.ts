import { f } from "./functions"

export function delai<T>(delai?: Record<string, T>): T[] {
  return f.iterable(delai)
}
