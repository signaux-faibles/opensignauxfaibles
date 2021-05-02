import { CompanyDataValuesWithFlags, Siret, Siren } from "../RawDataTypes"

declare function emit(
  key: Siret | Siren,
  value: CompanyDataValuesWithFlags
): void

export function map(this: {
  _id: unknown
  value: CompanyDataValuesWithFlags
}): void {
  "use strict"
  emit(this.value.key, this.value)
}
