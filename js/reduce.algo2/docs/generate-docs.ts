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

// Recursively extract properties from the schema, including in `allOf` and `anyOf` arrays.
const getAllProps = (schema: TJS.Definition = {}): TJS.Definition => {
  const combination = schema.allOf || schema.anyOf
  return combination
    ? combination.reduce((allProps: TJS.Definition, subSchema) => {
        return typeof subSchema === "object"
          ? { ...allProps, ...getAllProps(subSchema) }
          : allProps
      }, {})
    : schema.properties ?? {}
}

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

const getTypeDefsFromFiles = (
  typeName: string,
  filePaths: string[]
): { typeDefs: TJS.Definition[]; errors: Error[] } => {
  const program = TJS.getProgramFromFiles(filePaths)
  const generator = TJS.buildGenerator(program, settings)
  if (!generator) {
    throw new Error("failed to create generator")
  }
  const symbols = generator
    ?.getSymbols()
    .filter((smb) => smb.typeName === typeName)

  const errors: Error[] = []
  const typeDefs: TJS.Definition[] = []
  for (const symbol of symbols) {
    try {
      const fileName = path.basename(
        symbol?.fullyQualifiedName?.split('"')[1] ?? ""
      )
      console.warn(`- found ${typeName} type in file: ${fileName}.ts`)
      const typeDef = generator?.getSchemaForSymbol(symbol?.name ?? "") // can throw
      if (!typeDef) {
        throw new Error(`no schema for ${symbol?.fullyQualifiedName}`)
      }
      typeDefs.push(typeDef)
    } catch (err) {
      errors.push(err)
    }
  }
  return { typeDefs, errors }
}

// Generate the documentation of variables by importing types from a TypeScript file.
function documentPropertiesFromTypeDef(
  vars: TJS.Definition
): VarDocumentation[] {
  const transmittedVars = vars?.properties?.transmitted as TJS.Definition
  const computedVars = vars?.properties?.computed as TJS.Definition
  const source = ((vars?.properties?.source as TJS.Definition)?.enum ||
    [])[0]?.toString()

  return [
    ...documentProps(getAllProps(transmittedVars), { computed: false, source }),
    ...documentProps(getAllProps(computedVars), { computed: true, source }),
  ]
}

const { typeDefs, errors } = getTypeDefsFromFiles("Variables", INPUT_FILES)
errors.forEach(console.error)

const varsByFile = typeDefs.map((typeDef) => {
  const documentedVars = documentPropertiesFromTypeDef(typeDef)
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
