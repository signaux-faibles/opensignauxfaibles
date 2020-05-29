// Cette fonction TypeScript permet de vérifier que seuls les types reconnus
// peuvent être intégrés dans un BatchValue de destination.
// Ex: setBatchValueForType(batchValue, "pouet", {}) cause une erreur ts(2345).
export function setBatchValueForType<T extends keyof BatchValue>(
  batchValue: BatchValue,
  typeName: T,
  updatedValues: BatchValue[T]
): void {
  batchValue[typeName] = updatedValues
}
