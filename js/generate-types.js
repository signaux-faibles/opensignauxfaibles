const fs = require("fs")
const { compile } = require("json-schema-to-typescript")

const path = process.argv[2]

const normalizeType = (node) =>
  node.bsonType === "object"
    ? {
        ...node,
        properties: Object.entries(node.properties).reduce(
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
          .replace("int64", "number")
          .replace("double", "number"),
      }

const options = {
  bannerComment: "",
}

const convertFile = (filePath) => {
  const schema = require(filePath)
  return compile(normalizeType(schema), "", options).then((ts) =>
    ts
      .replace(
        /export interface ([^ ]+) \{/,
        `export interface ${schema.title} {`
      )
      .trim()
  )
}

fs.promises.readdir(path).then(async (files) => {
  const sortedFiles = [
    "apconso.schema.json",
    "apdemande.schema.json",
    "bdf.schema.json",
    "ccsf.schema.json",
    "compte.schema.json",
    "cotisation.schema.json",
    "debit.schema.json",
    "delai.schema.json",
    "diane.schema.json",
    "effectif_ent.schema.json",
    "effectif.schema.json",
    "ellisphere.schema.json",
    "paydex.schema.json",
    "procol.schema.json",
    "sirene_ul.schema.json",
    "sirene.schema.json",
  ] // TODO: remove the hard-coded list above, in favor to the default sorting of files:
  // const sortedFiles = files
  //   .filter((filename) => filename.endsWith(".schema.json"))
  for (const filename of sortedFiles) {
    console.log(await convertFile(`${path}/${filename}`))
  }
})
