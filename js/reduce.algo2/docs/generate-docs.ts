import * as fs from "fs"
import * as path from "path"
import * as TJS from "typescript-json-schema"

// Documentation schema for each variable.
type VarDocumentation = {
  name: string
  type: string
  computed: boolean
  description: string
}

const INPUT_FILES = process.argv.slice(2)

// Settings for typescript-json-schema
const settings: TJS.PartialArgs = {
  ignoreErrors: true,
  ref: false,
}

// Recursively extract properties from the schema, including in `allOf` arrays.
const getAllProps = (schema: TJS.Definition): TJS.Definition =>
  schema.allOf
    ? schema.allOf.reduce((allProps: TJS.Definition, subSchema) => {
        return typeof subSchema === "object"
          ? { ...allProps, ...getAllProps(subSchema) }
          : allProps
      }, {})
    : schema.properties ?? {}

// Generate the documentation of variables from properties of a JSON Schema.
const documentProps = (
  props: TJS.Definition = {},
  attributes: Partial<VarDocumentation>
): VarDocumentation[] =>
  Object.entries(props).map(([key, value]) => ({
    name: key,
    ...value,
    ...attributes,
  }))

// Generate the documentation of variables by importing types from a TypeScript file.
function documentPropertiesFromTypeDef(filePath: string): VarDocumentation[] {
  const program = TJS.getProgramFromFiles([filePath])
  const generator = TJS.buildGenerator(program, settings)

  const transmittedVars = generator?.getSchemaForSymbol("TransmittedVariables")
  const computedVars = generator?.getSchemaForSymbol("ComputedVariables")

  return [
    ...documentProps(transmittedVars?.properties, { computed: false }),
    ...documentProps(getAllProps(computedVars || {}), { computed: true }),
  ]
}

INPUT_FILES.forEach((filePath) => {
  const outFile = path.basename(filePath).replace(".ts", ".out.json")
  console.warn(`Generating ${outFile}...`)
  const documentedVars = documentPropertiesFromTypeDef(filePath)
  documentedVars.forEach((varDoc) => {
    if (!varDoc.description) {
      console.error(`variable ${varDoc.name} has no description`)
    }
  })
  fs.writeFileSync(outFile, JSON.stringify(documentedVars, null, 4))
})
