import { ParPériode } from "../../common/ParPériode"

/**
 * parPériode convertit un objet indexé par des dates exprimées sous la
 * forme `YYYY-MM-DD` en objet indexé par le timestamp de ces Dates,
 * par souci de compatibilité avec le type ParPériode.
 */
export const parPériode = <U>(indexed: Record<string, U>): ParPériode<U> => {
  const res = new ParPériode<U>()
  Object.entries(indexed).forEach(([k, v]) => {
    res.set(new Date(k), v)
  })
  return res
}
