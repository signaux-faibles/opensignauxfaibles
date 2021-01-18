import * as fs from "fs"
import { resolve } from "path"

import * as TJS from "typescript-json-schema"

const settings: TJS.PartialArgs = {
  ignoreErrors: true,
  ref: false,
}

const compilerOptions: TJS.CompilerOptions = {
  // strictNullChecks: true,
}

// optionally pass a base path
const basePath = "./my-dir"

const program = TJS.getProgramFromFiles(
  [resolve("../entr_diane.ts")],
  compilerOptions,
  basePath
)

const generator = TJS.buildGenerator(program, settings)
generator?.setSchemaOverride
const thru = generator?.getSchemaForSymbol("TransmittedVariables")
const computed1 = generator?.getSchemaForSymbol("ComputedVariables")?.allOf![0]
const computed2 = generator?.getSchemaForSymbol("ComputedVariables")?.allOf![1]
const computed3 = generator?.getSchemaForSymbol("ComputedVariables")?.allOf![2]

// Addition to the JSON Schema standard
type Attributes = {
  computed: boolean
}

const documentProp = (
  props: Record<string, TJS.DefinitionOrBoolean> = {},
  attributes: Attributes
) =>
  Object.entries(props).map(([key, value]) => ({
    name: key,
    ...(typeof value === "boolean" ? undefined : value),
    ...attributes,
  }))

const schema = {
  description:
    "Variables Diane générées par reduce.algo2 (opensignauxfaibles/sfdata)",
  type: "object",
  properties: [
    ...documentProp(thru?.properties, { computed: false }),
    ...documentProp(
      typeof computed1 === "boolean" ? {} : computed1?.properties,
      { computed: true }
    ),
    ...documentProp(
      typeof computed2 === "boolean" ? {} : computed2?.properties,
      { computed: true }
    ),
    ...documentProp(
      typeof computed3 === "boolean" ? {} : computed3?.properties,
      { computed: true }
    ),
  ],
}

fs.writeFileSync("entr_diane.out.json", JSON.stringify(schema, null, 4))
