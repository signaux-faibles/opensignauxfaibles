import { BatchKey } from "../RawDataTypes"
import { NAF } from "./populateNafAndApe"

// Paramètres globaux utilisés par "reduce.algo2"
export declare const offset_effectif: number
export declare const actual_batch: BatchKey
export declare const date_fin: Date
export declare const serie_periode: Date[]
export declare const naf: NAF
export declare const includes: Record<"all" | "apart", boolean>
