function cibleApprentissage(output_indexed, n_months) {
  "use strict"

  // Mock two input instead of one for future modification
  const output_cotisation = output_indexed
  const output_procol = output_indexed
  // replace with const
  const all_keys = Object.keys(output_indexed)
  //

  const merged_info = all_keys.reduce(function (m, k) {
    m[k] = {
      outcome: Boolean(
        output_procol[k].tag_failure || output_cotisation[k].tag_default
      ),
    }
    return m
  }, {})

  const output_outcome = f.lookAhead(merged_info, "outcome", n_months, true)
  const output_default = f.lookAhead(
    output_cotisation,
    "tag_default",
    n_months,
    true
  )
  const output_failure = f.lookAhead(
    output_procol,
    "tag_failure",
    n_months,
    true
  )

  const output_cible = all_keys.reduce(function (m, k) {
    m[k] = {}

    if (output_outcome[k]) m[k] = output_outcome[k]
    if (output_default[k])
      m[k].time_til_default = output_default[k].time_til_outcome
    if (output_failure[k])
      m[k].time_til_failure = output_failure[k].time_til_outcome
    return m
  }, {})

  return output_cible
}

exports.cibleApprentissage = cibleApprentissage
