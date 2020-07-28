import "../globals.ts"

declare function emit(
  key: SiretOrSiren,
  value: CompanyDataValuesWithFlags
): void

export function map(this: { value: CompanyDataValuesWithFlags }): void {
  "use strict"
  if (typeof this.value !== "object") {
    throw new Error("this.value should be a valid object, in compact::map()")
  }
  emit(this.value.key, this.value)
}
