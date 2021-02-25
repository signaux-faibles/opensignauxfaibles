import { f } from "./functions"
import { ProcolToHumanRes } from "../common/procolToHuman"
import { EntréeDéfaillances } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"
import { ParPériode } from "../common/newParPériode"

export type SortieDefaillances = {
  /** État de la procédure collective. */
  etat_proc_collective: ProcolToHumanRes
  /** Date effet de la procédure collective. */
  date_proc_collective: Date
  /** État de défaillance. (c.a.d. l'entité n'est pas "in bonis") */
  tag_failure: boolean
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "defaillances"
  computed: SortieDefaillances
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function defaillances(
  défaillances: ParHash<EntréeDéfaillances>,
  output_indexed: ParPériode<Partial<SortieDefaillances>>
): void {
  "use strict"
  const codes = Object.keys(défaillances)
    .reduce((events, hash) => {
      const the_event = défaillances[hash] as EntréeDéfaillances

      let etat: ProcolToHumanRes = null
      etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol)

      if (etat !== null)
        events.push({
          etat,
          date_proc_col: new Date(the_event.date_effet),
        })

      return events
    }, [] as { etat: ProcolToHumanRes; date_proc_col: Date }[])
    .sort((a, b) => {
      return a.date_proc_col.getTime() - b.date_proc_col.getTime()
    })

  codes.forEach((event) => {
    const periode_effet = new Date(
      Date.UTC(
        event.date_proc_col.getFullYear(),
        event.date_proc_col.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    const time_til_last = [...output_indexed.keys()].filter((période) => {
      return période >= periode_effet.getTime()
    })

    time_til_last.forEach((time) => {
      const outputForTime = output_indexed.get(time) ?? {}
      if (outputForTime !== undefined) {
        outputForTime.etat_proc_collective = event.etat
        outputForTime.date_proc_collective = event.date_proc_col
        if (event.etat !== "in_bonis") outputForTime.tag_failure = true
      }
      output_indexed.set(time, outputForTime)
    })
  })
}
