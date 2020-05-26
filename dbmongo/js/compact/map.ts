import "../globals.ts"

export function map(this: { value: CompanyDataValues }): void {
  "use strict"
  if (typeof this.value !== "object") {
    throw new Error("this.value should be a valid object, in compact::map()")
  }
  emit(this.value.key, this.value)
}
