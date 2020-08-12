export function iterable<T>(dict?: Record<string | number, T>): T[] {
  return typeof dict === "object" ? Object.keys(dict).map((h) => dict[h]) : []
}
