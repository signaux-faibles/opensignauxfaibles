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
const thru = generator?.getSchemaForSymbol("DonnéesDianeTransmises")
const computed1 = generator?.getSchemaForSymbol("SortieDiane")?.allOf![1]
const computed2 = generator?.getSchemaForSymbol("SortieDiane")?.allOf![2]
const computed3 = generator?.getSchemaForSymbol("SortieDiane")?.allOf![3]

// Addition to the JSON Schema standard
type Additions = {
  computed: boolean
}

const appendToValues = (
  props: Record<string, TJS.DefinitionOrBoolean> = {},
  additions: Additions
) =>
  Object.entries(props).reduce(
    (res, [key, value]) => ({
      ...res,
      [key]: typeof value === "boolean" ? value : { ...additions, ...value },
    }),
    {}
  )

const schema = {
  description:
    "Variables Diane générées par reduce.algo2 (opensignauxfaibles/sfdata)",
  type: "object",
  properties: {
    ...appendToValues(thru?.properties, { computed: false }),
    ...appendToValues(
      typeof computed1 === "boolean" ? {} : computed1?.properties,
      { computed: true }
    ),
    ...appendToValues(
      typeof computed2 === "boolean" ? {} : computed2?.properties,
      { computed: true }
    ),
    ...appendToValues(
      typeof computed3 === "boolean" ? {} : computed3?.properties,
      { computed: true }
    ),
  },
}

fs.writeFileSync("entr_diane.out.json", JSON.stringify(schema, null, 4))
