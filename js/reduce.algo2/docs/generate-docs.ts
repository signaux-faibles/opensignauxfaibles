import * as path from "path"
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
  uniqueNames: true,
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

const getTypeDefFromFile = (
  generator: TJS.JsonSchemaGenerator,
  typeName: string,
  filePath: string
): TJS.Definition | undefined => {
  const filename = path.basename(filePath).replace(/\.ts$/, "")
  const symbolName = generator
    ?.getSymbols()
    .find(
      (smb) =>
        smb.typeName === typeName && smb.fullyQualifiedName.includes(filename)
    )?.name
  try {
    return generator?.getSchemaForSymbol(symbolName ?? "") // can throw
  } catch (err) {
    return undefined
  }
}

// Generate the documentation of variables by importing types from a TypeScript file.
function documentPropertiesFromTypeDef(filePath: string): VarDocumentation[] {
  const program = TJS.getProgramFromFiles([filePath])
  const generator = TJS.buildGenerator(program, settings)
  if (!generator) {
    throw new Error("failed to create generator")
  }

  const vars = getTypeDefFromFile(generator, "Variables", filePath)
  const transmittedVars = vars?.properties?.transmitted as TJS.Definition
  const computedVars = vars?.properties?.computed as TJS.Definition
  const source = ((vars?.properties?.source as TJS.Definition)?.enum ||
    [])[0]?.toString()

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
