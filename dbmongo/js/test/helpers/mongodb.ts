/*global globalThis*/

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

;(Object as any).bsonsize = (obj: unknown): number => JSON.stringify(obj).length // eslint-disable-line @typescript-eslint/no-explicit-any

type Document = Record<string, unknown>
type MapResult<K, V> = { _id: K; value: V }

// Run a map() function designed for MongoDB, i.e. that calls emit() an
// inderminate number of times, instead of returning one value per iteration.
export const runMongoMap = <DocumentId, Doc extends Document>(
  mapFct: (this: Doc) => void, // will call global emit()
  documents: Doc[]
): MapResult<DocumentId, Document>[] => {
  const results: MapResult<DocumentId, Document>[] = [] // holds all the { _id, value } objects emitted from mapFct()
  // define a emit() function that mapFct() can call
  global.emit = (_id: DocumentId, value: Document): void => {
    results.push({ _id, value })
  }
  documents.forEach((doc) => mapFct.call(doc))
  return results
}

type Indexed<K, V> = Record<string, { key: K; value: V }[]>

export const indexMapResultsByKey = <K, V>(
  mapResults: MapResult<K, V>[]
): Indexed<K, V> =>
  mapResults.reduce((acc, { _id, value }) => {
    const key = JSON.stringify(_id) // e.g. _id: { siren; batch; periode }
    acc[key] = acc[key] || []
    acc[key].push({ key: _id, value })
    return acc
  }, {} as Indexed<K, V>)

type SerializedDate = { _ISODate: string }

// Converts a "JSON" object returned by MongoDB (including mentions to ISODate
// and NumberInt) into a valid JavaScript object with Date instances.
export const parseMongoObject = (serializedObj: string): unknown =>
  JSON.parse(
    serializedObj
      .replace(/ISODate\("([^"]+)"\)/g, `{ "_ISODate": "$1" }`)
      .replace(/NumberInt\(([^)]+)\)/g, "$1"),
    (_key, value: SerializedDate | unknown) =>
      value && typeof value === "object" && (value as SerializedDate)._ISODate
        ? new Date((value as SerializedDate)._ISODate)
        : value
  )

const isStringifiedDate = (date: string | unknown) =>
  typeof date === "string" && /\d{4}-\d{2}-\d{2}T\d{2}:\d{2}.*Z/.test(date)

// Converts an object into the same format as the one returned by the
// `find().toArray()` command when executed from  MongoDB's mongo shell.
// E.g. Dates are serialized as ISODate() instances.
// Thanks to this function, algo2_golden_tests.ts and test-api-reduce-2.sh
// can produce the exact same content, when updating.
export const serializeAsMongoObject = (obj: unknown): string => {
  const INDENT_SPACES = 2
  return (
    JSON.stringify(
      obj,
      (_, val) =>
        isStringifiedDate(val) ? `ISODate_${val.replace(/\.000Z/, "Z")}` : val,
      INDENT_SPACES
    )
      .replace(/"ISODate_([^"]+)"/g, `ISODate("$1")`) // replace ISODate strings by function calls
      .replace(/":/g, `" :`) // formatting: add a space before property assignments
      .replace(new RegExp(` {${INDENT_SPACES}}`, "g"), "\t") + "\n" // formatting: tabs and trailing line break
  )
}
