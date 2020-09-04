import * as f from "./iterable"

export function compte<T>(compte?: Record<string, T>): T | undefined {
  const c = f.iterable(compte)
  return c.length > 0 ? c[c.length - 1] : undefined
}
