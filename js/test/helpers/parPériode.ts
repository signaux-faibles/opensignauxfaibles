export const parPériode = <T extends Record<number, unknown>>(
  indexed: Record<string, T[keyof T]>
): T => {
  const res = {} as T
  Object.entries(indexed).forEach(([k, v]) => {
    res[new Date(k).getTime()] = v
  })
  return res
}
