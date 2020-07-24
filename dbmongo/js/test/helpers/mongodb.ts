const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

export const runMongoMap = (
  mapFct: () => void,
  keyVal: unknown
): Record<string, unknown> => {
  const results: Record<string, unknown> = {}
  global.emit = (key: string, value: CompanyDataValuesWithFlags): void => {
    results[key] = value
  }
  mapFct.call(keyVal)
  return results
}
