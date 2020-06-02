export function flatten(
  v: { key: any; scope: Scope; batch: BatchValues },
  actual_batch: number
): { key: any; scope: Scope }
