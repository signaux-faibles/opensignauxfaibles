export function repeatable(rep: {
  [hash: string]: { periode: Periode; random_order: number }
}): { [key: string]: { random_order: number } }
