const Ajv = require("ajv").default // eslint-disable-line @typescript-eslint/no-var-requires
const useMongoDbPlugin = require("ajv-bsontype") // eslint-disable-line @typescript-eslint/no-var-requires

// Initialize JSON Schema validator with MongoDB extensions (e.g. "bsonType" property)
const ajv = new Ajv({ strict: false })
useMongoDbPlugin(ajv)

const readStdin = () =>
  new Promise((resolve) => {
    const chunks = []
    process
      .openStdin()
      .on("data", (chunk) => chunks.push(chunk))
      .on("end", () => resolve(chunks.join("")))
  })

readStdin()
  .then((stdin) =>
    stdin
      .split("\n")
      .filter((line) => line.trim() !== "")
      .forEach((line) => {
        console.log(line)
        const { dataType, dataObject } = JSON.parse(line)
        const schema = require(`../validation/${dataType}.schema.json`) // eslint-disable-line @typescript-eslint/no-var-requires
        ajv.validate(schema, dataObject)
        ajv.errors.forEach(({ message, params }) =>
          console.log("- Validation error:", `"${message}"`, params)
        )
      })
  )
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })
