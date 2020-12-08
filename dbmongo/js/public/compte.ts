export function compte<T>(compte?: Record<string, T>): T | undefined {
  const c = Object.values(compte ?? {})
  return c.length > 0 ? c[c.length - 1] : undefined
}
