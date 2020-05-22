declare function emit(key: any, value: any): void
declare function print(...any): void

export function map() {
  "use strict"
  try {
    if (this.value != null) {
      emit(this.value.key, this.value)
    }
  } catch (error) {
    print(this.value.key)
  }
}
