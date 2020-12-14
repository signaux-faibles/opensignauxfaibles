export function delai<T>(delai?: Record<string, T>): T[] {
  return Object.values(delai ?? {})
}
