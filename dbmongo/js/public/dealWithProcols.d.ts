export function dealWithProcols(
  data_source: {
    [h: string]: {
      code_evenement: any
      action_procol: any
      stade_procol: any
      date_effet: any
    }
  },
  altar_or_procol: "altares" | "procol",
  output_indexed: any
): { etat: any; date_procol: Date }[]
