export function finalize<T extends { [siret]: any }>(
  k: any,
  v: T
): T | { incomplete: true }
