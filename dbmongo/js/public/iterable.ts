export function iterable<T>(dict: Record<string | number, T>): T[] {
  try {
    return Object.keys(dict).map((h) => dict[h])
  } catch (error) {
    return []
  }
}
