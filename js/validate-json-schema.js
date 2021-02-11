var Ajv = require("ajv").default
var ajv = new Ajv({ strict: false })
require("ajv-bsontype")(ajv)

const readStdin = () =>
  new Promise((resolve) => {
    const chunks = []
    process
      .openStdin()
      .on("data", (chunk) => chunks.push(chunk))
      .on("end", () => resolve(chunks.join("")))
  })

;(async () => {
  ;(await readStdin())
    .split("\n")
    .filter((line) => line.trim() != "")
    .forEach((line) => {
      console.log(line)
      const { dataType, dataObject } = JSON.parse(line)
      const schema = require(`../validation/${dataType}.schema.json`) // eslint-disable-line @typescript-eslint/no-var-requires
      const valid = ajv.validate(schema, dataObject)
      ajv.errors.forEach(({ message, params }) =>
        console.log("- Validation error:", `"${message}"`, params)
      )
    })
})().catch((err) => {
  console.error(err)
  process.exit(1)
})
