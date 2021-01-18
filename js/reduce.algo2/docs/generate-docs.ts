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

const program = TJS.getProgramFromFiles(
  [resolve("../entr_diane.ts")],
  compilerOptions
)

const getAllProps = (schema: TJS.Definition): TJS.Definition =>
  schema.allOf
    ? schema.allOf.reduce((allProps: TJS.Definition, subSchema) => {
        return typeof subSchema === "object"
          ? { ...allProps, ...getAllProps(subSchema) }
          : allProps
      }, {})
    : schema.properties ?? {}

const generator = TJS.buildGenerator(program, settings)
generator?.setSchemaOverride
const transmittedProps = generator?.getSchemaForSymbol("TransmittedVariables")
  ?.properties
const computedProps = getAllProps(
  generator?.getSchemaForSymbol("ComputedVariables") || {}
)

// Addition to the JSON Schema standard
type Attributes = {
  computed: boolean
}

const documentProps = (props: TJS.Definition = {}, attributes: Attributes) =>
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
    ...documentProps(transmittedProps, { computed: false }),
    ...documentProps(computedProps, { computed: true }),
  ],
}

fs.writeFileSync("entr_diane.out.json", JSON.stringify(schema, null, 4))
