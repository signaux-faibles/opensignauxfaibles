export function setBatchValueForType(
  batchValue: BatchValue,
  typeName: keyof BatchValue,
  updatedValues: BatchValue[keyof BatchValue]
): void {
  switch (typeName) {
    case "reporder":
      batchValue[typeName] = updatedValues as BatchValue["reporder"]
      break
    case "compact":
      batchValue[typeName] = updatedValues as BatchValue["compact"]
      break
    case "effectif":
      batchValue[typeName] = updatedValues as BatchValue["effectif"]
      break
    case "apconso":
      batchValue[typeName] = updatedValues as BatchValue["apconso"]
      break
    case "apdemande":
      batchValue[typeName] = updatedValues as BatchValue["apdemande"]
      break
    default:
      // This switch should be exhaustive: cover all the keys defined in the BatchValue type.
      // => Warning TS(2345) if we miss a case, e.g. Argument of type '"new_effectif"' is not assignable to parameter of type 'never'.
      // source: https://stackoverflow.com/a/61806149/592254
      // eslint-disable-next-line no-extra-semi
      ;((caseVal: never): void => {
        throw new Error(`case "${caseVal}" should be added to switch`)
      })(typeName)
  }
}
