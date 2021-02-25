/**
 * parPériode convertit un objet indexé par des dates exprimées sous la
 * forme `YYYY-MM-DD` en objet indexé par le timestamp de ces Dates,
 * par souci de compatibilité avec le type ParPériode.
 */
export const parPériode = <T extends Record<number, unknown>>(
  indexed: Record<string, T[keyof T]>
): T => {
  const res = {} as T
  Object.entries(indexed).forEach(([k, v]) => {
    res[new Date(k).getTime()] = v
  })
  return res
}
