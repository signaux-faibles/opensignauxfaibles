export function cibleApprentissage(
  output_indexed: { [k: string]: { tag_failure; tag_default } },
  n_months: number
): { [k: string]: { time_til_default; time_til_failure } }
