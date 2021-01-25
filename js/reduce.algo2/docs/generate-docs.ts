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

// Generator that yields a function to get the definition of a type for each file.
function* extractTypeDefsFromFiles(
  typeName: string,
  filePaths: string[]
): Generator<{ fileName: string; getTypeDef: () => TJS.Definition }> {
  const program = TJS.getProgramFromFiles(filePaths)
  const generator = TJS.buildGenerator(program, settings)
  if (!generator) {
    throw new Error("failed to create generator")
  }
  const symbols = generator
    .getSymbols()
    .filter((smb) => smb.typeName === typeName)
  for (const symbol of symbols) {
    // We let the caller call getTypeDef(), for each fileName
    yield {
      fileName: path.basename(symbol.fullyQualifiedName?.split('"')[1] ?? ""),
      getTypeDef: () => generator.getSchemaForSymbol(symbol.name), // can throw
    }
  }
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

const missingDesc = (varDoc: VarDocumentation) =>
  !varDoc.description && !varDoc.name.includes("_past_")

// Main script: print the documented variables in JSON format + log errors.

const allDocumentedVars = []
const typeDefGeneretor = extractTypeDefsFromFiles("Variables", INPUT_FILES)
for (const { fileName, getTypeDef } of typeDefGeneretor) {
  console.warn(`- extracting variables from file: ${fileName}.ts`)
  const typeDef = getTypeDef() // can throw
  const documentedVars = documentPropertiesFromTypeDef(typeDef)
  documentedVars
    .filter(missingDesc)
    .forEach(({ name }) => console.error(`  ⚠ no description for: ${name}`))
  allDocumentedVars.push(...documentedVars)
}

console.log(JSON.stringify(allDocumentedVars, null, 4))
