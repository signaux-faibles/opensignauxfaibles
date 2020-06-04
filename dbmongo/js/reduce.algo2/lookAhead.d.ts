type Outcome = { time_til_outcome: number; outcome: boolean }

export function lookAhead(
  data: { [period: string]: any },
  attr_name: string,
  n_months: number,
  past: boolean
): { [period: string]: Outcome }
