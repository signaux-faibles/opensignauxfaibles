type Event = { code_evenement; action_procol; stade_procol; date_effet }

export function dealWithProcols(
  data_source: { [hash: string]: Event },
  altar_or_procol: "altares" | "procol",
  output_indexed: {
    [time: string]: {
      etat_proc_collective: any
      date_proc_collective: any
    }
  }
): void
