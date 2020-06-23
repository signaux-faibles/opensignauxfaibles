function interim(interim, output_indexed) {
  "use strict"
  const output_effectif = output_indexed
  // let periodes = Object.keys(output_indexed)
  // output_indexed devra être remplacé par output_effectif, et ne contenir que les données d'effectif.
  // periodes sera passé en argument.

  const output_interim = {}

  //  var offset_interim = 3

  Object.keys(interim).forEach((hash) => {
    const one_interim = interim[hash]
    const periode = one_interim.periode.getTime()
    // var periode_d = new Date(parseInt(interimTime))
    // var time_offset = f.dateAddMonth(time_d, -offset_interim)
    if (periode in output_effectif) {
      output_interim[periode] = output_interim[periode] || {}
      output_interim[periode].interim_proportion =
        one_interim.etp / output_effectif[periode].effectif
    }

    const past_month_offsets = [6, 12, 18, 24]
    past_month_offsets.forEach((offset) => {
      const time_past_offset = f.dateAddMonth(one_interim.periode, offset)
      const variable_name_interim = "interim_ratio_past_" + offset
      if (
        periode in output_effectif &&
        time_past_offset.getTime() in output_effectif
      ) {
        output_interim[time_past_offset.getTime()] =
          output_interim[time_past_offset.getTime()] || {}
        const val_offset = output_interim[time_past_offset.getTime()]
        val_offset[variable_name_interim] =
          one_interim.etp / output_effectif[periode].effectif
      }
    })
  })

  return output_interim
}
