import * as fs from "fs"
import * as path from "path"
import * as TJS from "typescript-json-schema"

const INPUT_FILES = process.argv.slice(2)

const settings: TJS.PartialArgs = {
  ignoreErrors: true,
  ref: false,
}

const compilerOptions: TJS.CompilerOptions = {
  // strictNullChecks: true,
}

const getAllProps = (schema: TJS.Definition): TJS.Definition =>
  schema.allOf
    ? schema.allOf.reduce((allProps: TJS.Definition, subSchema) => {
        return typeof subSchema === "object"
          ? { ...allProps, ...getAllProps(subSchema) }
          : allProps
      }, {})
    : schema.properties ?? {}

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

function documentPropertiesFromTypeDef(filePath: string) {
  const outFile = path.basename(filePath).replace(".ts", ".out.json")
  console.warn(`Generating ${outFile}...`)

  const program = TJS.getProgramFromFiles([filePath], compilerOptions)

  const generator = TJS.buildGenerator(program, settings)
  generator?.setSchemaOverride

  const transmittedProps = generator?.getSchemaForSymbol("TransmittedVariables")
    ?.properties

  const computedProps = getAllProps(
    generator?.getSchemaForSymbol("ComputedVariables") || {}
  )

  const schema = {
    description:
      "Variables Diane générées par reduce.algo2 (opensignauxfaibles/sfdata)",
    type: "object",
    properties: [
      ...documentProps(transmittedProps, { computed: false }),
      ...documentProps(computedProps, { computed: true }),
    ],
  }

  fs.writeFileSync(outFile, JSON.stringify(schema, null, 4))
}

INPUT_FILES.forEach(documentPropertiesFromTypeDef)
