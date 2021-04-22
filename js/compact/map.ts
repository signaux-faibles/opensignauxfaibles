import { CompanyDataValuesWithFlags, SiretOrSiren } from "../RawDataTypes"

declare function emit(
  key: SiretOrSiren,
  value: CompanyDataValuesWithFlags
): void

export function map(this: {
  _id: unknown
  value: CompanyDataValuesWithFlags
}): void {
  "use strict"
  emit(this.value.key, this.value)
}
