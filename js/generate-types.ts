import * as fs from "fs"
import { compile, JSONSchema } from "json-schema-to-typescript"

const path = process.argv[2]

type JSONProps = JSONSchema["properties"]

const normalizePropTypes = (properties: JSONProps): JSONProps =>
  Object.entries(properties ?? {}).reduce(
    (acc, [propName, propDef]) => ({
      ...acc,
      [propName]: normalizeType(propDef),
    }),
    {}
  )

const tsTypes = new Map<string, JSONSchema["tsType"]>([["date", "Date"]])

const jsTypes = new Map<string, JSONSchema["type"]>([
  ["bool", "boolean"],
  ["long", "number"],
  ["double", "number"],
])

const normalizeType = (node: JSONSchema): JSONSchema =>
  node.bsonType === "object"
    ? {
        ...node,
        properties: normalizePropTypes(node.properties),
      }
    : tsTypes.has(node.bsonType)
    ? {
        ...node,
        tsType: tsTypes.get(node.bsonType),
      }
    : {
        ...node,
        type: jsTypes.get(node.bsonType) ?? node.bsonType,
      }

const options = {
  bannerComment: "",
}

const convertFile = (filePath: string) => {
  const schema: JSONSchema = require(filePath) // eslint-disable-line @typescript-eslint/no-var-requires
  return compile(normalizeType(schema), "", options).then((ts) =>
    ts
      .replace(
        /export interface ([^ ]+) \{/,
        `export interface ${schema.title} {`
      )
      .trim()
  )
}

fs.promises
  .readdir(path ?? ".")
  .then((files) =>
    Promise.all(
      files
        .filter((filename) => filename.endsWith(".schema.json"))
        .map((filename) => convertFile(`${path}/${filename}`))
    ).then((tsDefs) => tsDefs.map((tsDef) => console.log(tsDef)))
  )
