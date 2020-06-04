export function finalize<T>(
  k: any,
  v: { [key: string]: T }
): T[] | { incomplete: true }
