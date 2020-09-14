// Paramètres globaux utilisés par "reduce.algo2"
let offset_effectif: number
let actual_batch: string // import { BatchKey } from "../RawDataTypes"
let date_fin: Date
let serie_periode: Date[]
let naf: any // import { NAF } from "./populateNafAndApe"
let includes: Record<"all" | "apart", boolean>
