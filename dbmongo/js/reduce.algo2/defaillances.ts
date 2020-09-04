import * as f from "./dealWithProcols"
import { SortieProcols } from "./dealWithProcols"
import { EntréeDefaillances, ParPériode } from "../RawDataTypes"

export type SortieDefaillances = SortieProcols

export function defaillances(
  altares: Record<string, EntréeDefaillances>,
  procol: Record<string, EntréeDefaillances>,
  output_indexed: ParPériode<Partial<SortieDefaillances>>
): void {
  "use strict"
  f.dealWithProcols(altares, "altares", output_indexed)
  f.dealWithProcols(procol, "procol", output_indexed)
}
