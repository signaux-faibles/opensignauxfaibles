import "../globals.ts"

export function map(): void {
  "use strict"
  try {
    // TODO: this.value is RawDataValues ?
    if (this.value != null) {
      emit(this.value.key, this.value)
    }
  } catch (error) {
    print(this.value.key)
  }
}
