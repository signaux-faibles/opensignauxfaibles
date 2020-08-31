type Doc<T> = { key: unknown; value: T }

export function reducer<T>(
  array: Doc<T>[],
  reduce: (key: unknown, values: T[]) => T
): T

export function invertedReducer<T>(
  array: Doc<T>[],
  reduce: (key: unknown, values: T[]) => T
): T
