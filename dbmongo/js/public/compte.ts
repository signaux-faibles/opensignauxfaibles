import { iterable } from "./iterable"

export function compte<T>(compte?: { [key: string]: T }): T | undefined {
  const f = { iterable } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  const c = f.iterable(compte)
  return c.length > 0 ? c[c.length - 1] : undefined
}
