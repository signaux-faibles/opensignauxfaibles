import * as f from "./iterable"

export function delai<T>(delai: { [key: string]: T }): T[] {
  return f.iterable(delai)
}
