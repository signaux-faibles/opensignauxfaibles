// Ce script sert à générer GeneratedTypes.d.ts à partir de fichiers *.schema.json
// dans lesquels les types sont exprimés à l'aide de bsonType.
// (cf https://docs.mongodb.com/manual/reference/operator/query/type/#document-type-available-types)
// (cf la table de correspondance goTypeToBsonType de validation.go.)

import * as fs from "fs"
import { compile, JSONSchema } from "json-schema-to-typescript"

const path = process.argv[2]

type JSONProps = JSONSchema["properties"]

type SchemaMutation = (node: JSONSchema) => JSONSchema

const mutateEach = (properties: JSONProps, mutate: SchemaMutation): JSONProps =>
  Object.entries(properties ?? {}).reduce(
    (acc, [propName, propDef]) => ({
      ...acc,
      [propName]: mutate(propDef),
    }),
    {}
  )

const tsTypes = new Map<string, JSONSchema["tsType"]>([["date", "Date"]])

const jsTypes = new Map<string, JSONSchema["type"]>([
  ["bool", "boolean"],
  ["long", "number"],
  ["double", "number"],
])

// Conversion recursive de bsonType en types reconnus par json-schema-to-typescript
const normalizeType = (node: JSONSchema): JSONSchema =>
  node.bsonType === "object"
    ? {
        ...node,
        properties: mutateEach(node.properties, normalizeType),
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

const DEFAULT_OPTIONS = {
  bannerComment: "",
}

const convertFile = async (filePath: string, options = DEFAULT_OPTIONS) => {
  const rawSchema = await fs.promises.readFile(filePath, "utf-8")
  const schema: JSONSchema = JSON.parse(rawSchema)
  const typeDef = await compile(normalizeType(schema), "", options)
  return typeDef.replace(
    /export interface ([^ ]+) \{/,
    `export interface ${schema.title} {` // ré-injection du nom, pour que les accents soient conservés
  )
}

fs.promises
  .readdir(path ?? ".")
  .then((files) =>
    Promise.all(
      files
        .filter((filename) => filename.endsWith(".schema.json"))
        .map((filename) => convertFile(`${path}/${filename}`))
    ).then((tsDefs) => tsDefs.map((tsDef) => process.stdout.write(tsDef)))
  )
