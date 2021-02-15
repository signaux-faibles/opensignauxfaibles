import * as fs from "fs"
import { compile, JSONSchema } from "json-schema-to-typescript"

const path = process.argv[2]

const normalizeType = (node: JSONSchema): JSONSchema =>
  node.bsonType === "object"
    ? {
        ...node,
        properties: Object.entries(node.properties ?? {}).reduce(
          (acc, [propName, propDef]) => ({
            ...acc,
            [propName]: normalizeType(propDef),
          }),
          {}
        ),
      }
    : node.bsonType === "date"
    ? {
        ...node,
        tsType: "Date",
      }
    : {
        ...node,
        type: node.bsonType
          .replace("date", "Date")
          .replace("bool", "boolean")
          .replace("long", "number")
          .replace("double", "number"),
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
