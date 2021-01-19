import * as TJS from "typescript-json-schema"

// Documentation schema for each variable.
type VarDocumentation = {
  source: string
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
const getAllProps = (schema: TJS.Definition = {}): TJS.Definition =>
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
    source: attributes.source,
    name: key,
    ...value,
    ...attributes,
  }))

// Generate the documentation of variables by importing types from a TypeScript file.
function documentPropertiesFromTypeDef(filePath: string): VarDocumentation[] {
  const program = TJS.getProgramFromFiles([filePath])
  const generator = TJS.buildGenerator(program, settings)

  let transmittedVars: TJS.Definition | undefined
  try {
    transmittedVars = generator?.getSchemaForSymbol("TransmittedVariables")
  } catch (err) {
    console.error(err.message)
  }

  let computedVars: TJS.Definition | undefined
  try {
    computedVars = generator?.getSchemaForSymbol("ComputedVariables")
  } catch (err) {
    console.error(err.message)
  }

  let source: string | undefined
  try {
    source = (generator?.getSchemaForSymbol("VariablesSource")?.enum ||
      [])[0]?.toString()
  } catch (err) {
    console.error(err.message)
  }

  return [
    ...documentProps(transmittedVars?.properties, { computed: false, source }),
    ...documentProps(getAllProps(computedVars), { computed: true, source }),
  ]
}

const varsByFile = INPUT_FILES.map((filePath) => {
  const documentedVars = documentPropertiesFromTypeDef(filePath)
  documentedVars.forEach((varDoc) => {
    if (!varDoc.description && !varDoc.name.includes("_past_")) {
      console.error(`variable ${varDoc.name} has no description`)
    }
  })
  return documentedVars
})

console.log(
  JSON.stringify(
    varsByFile.reduce((acc, vars) => acc.concat(...vars), []),
    null,
    4
  )
)
