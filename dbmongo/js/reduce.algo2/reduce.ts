export function reduce<T>(_key: unknown, values: T[]): T {
  "use strict"
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {} as T)
}
