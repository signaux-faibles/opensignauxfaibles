/*global globalThis*/

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

;(Object as any).bsonsize = (obj: unknown): number => JSON.stringify(obj).length // eslint-disable-line @typescript-eslint/no-explicit-any

type Document<K> = { _id: K } & Record<string, unknown>
type MapResult<K, V> = { _id: K; value: V }

// Run a map() function designed for MongoDB, i.e. that calls emit() an
// inderminate number of times, instead of returning one value per iteration.
export const runMongoMap = <
  Key,
  InDoc extends Document<Key>,
  OutDoc extends Document<Key>
>(
  mapFct: (this: InDoc) => void, // will call global emit()
  documents: InDoc[]
): MapResult<Key, OutDoc>[] => {
  const results: MapResult<Key, OutDoc>[] = [] // holds all the { _id, value } objects emitted from mapFct()
  // define a emit() function that mapFct() can call
  global.emit = (_id: Key, value: OutDoc): void => {
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
export const serializeAsMongoObject = (obj: unknown): string =>
  JSON.stringify(
    obj,
    (_, val) =>
      isStringifiedDate(val) ? `ISODate_${val.replace(/\.000Z/, "Z")}` : val,
    "\t"
  )
    .replace(/"ISODate_([^"]+)"/g, `ISODate("$1")`) // replace ISODate strings by function calls
    .replace(/":/g, `" :`) + "\n" // formatting: add a space before property assignments + trailing line break

// Run a reduce() function designed for MongoDB, based on the values returned
// by runMongoMap().
export const runMongoReduce = <Key, Doc extends Document<Key>>(
  reduceFct: (_key: Key, values: Doc[]) => Doc,
  mapResults: MapResult<Key, Doc>[]
): MapResult<Key, Doc>[] => {
  const valuesPerKey: Record<string, MapResult<Key, Doc[]>> = {}
  mapResults.forEach(({ _id, value }) => {
    const idString = JSON.stringify(_id)
    valuesPerKey[idString] = valuesPerKey[idString] || { _id, value: [] } // TODO: renommer `value` --> `values`
    valuesPerKey[idString].value.push(value as Doc)
  })
  return Object.values(valuesPerKey).map(({ _id, value }) => ({
    _id,
    value: reduceFct(_id, value),
  }))
}
