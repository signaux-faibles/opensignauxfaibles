import { forEachPopulatedProp } from "../common/forEachPopulatedProp"
import { applyPatchesToBatch } from "./applyPatchesToBatch"
import { applyPatchesToMemory } from "./applyPatchesToMemory"
import { compactBatch } from "./compactBatch"
import { complete_reporder } from "./complete_reporder"
import { currentState } from "./currentState"
import { fixRedundantPatches } from "./fixRedundantPatches"
import { listHashesToAddAndDelete } from "./listHashesToAddAndDelete"

export const f = {
  forEachPopulatedProp,
  listHashesToAddAndDelete,
  applyPatchesToBatch,
  applyPatchesToMemory,
  fixRedundantPatches,
  compactBatch,
  currentState,
  complete_reporder,
}
