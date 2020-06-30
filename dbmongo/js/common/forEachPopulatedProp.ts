// Appelle fct() pour chaque propriété définie (non undefined) de obj.
// Contrat: obj ne doit contenir que les clés définies dans son type.
export function forEachPopulatedProp<T>(
  obj: T,
  fct: (key: keyof T, val: Required<T>[keyof T]) => unknown
): void {
  ;(Object.keys(obj) as Array<keyof T>).forEach((key) => {
    if (typeof obj[key] !== "undefined") fct(key, obj[key])
  })
}
