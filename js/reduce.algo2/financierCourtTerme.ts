import { EntréeDiane } from "../RawDataTypes"

export function financierCourtTerme(diane: EntréeDiane): number | null {
  "use strict"

  const ratio =
    (diane["concours_bancaire_courant"] ?? NaN) / (diane["ca"] ?? NaN)
  return isNaN(ratio) ? null : ratio * 100
}
