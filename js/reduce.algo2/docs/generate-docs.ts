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

// We can either get the schema for one file and one type...
const schema = TJS.generateSchema(program, "SortieDiane", settings)

// ... or a generator that lets us incrementally get more schemas
// const generator = TJS.buildGenerator(program, settings)
// all symbols
// const symbols = generator.getUserSymbols()

// Get symbols for different types from generator.
// generator.getSchemaForSymbol("MyType")
// generator.getSchemaForSymbol("AnotherType")

// console.log(schema)
fs.writeFileSync("entr_diane.out.json", JSON.stringify(schema, null, 4))
