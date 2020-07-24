declare let emit: (_id: unknown, value: unknown) => void

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

type DocumentId = unknown // can be an ObjectID, or other
type Document = Record<string, unknown>
type MapResult = { _id: DocumentId; value: Document }

// Run a map() function designed for MongoDB, i.e. that calls emit() an
// inderminate number of times, instead of returning one value per iteration.
export const runMongoMap = <Doc extends Document>(
  mapFct: (this: Doc) => void, // will call global emit()
  documents: Doc[]
): MapResult[] => {
  const results: MapResult[] = [] // holds all the { _id, value } objects emitted from mapFct()
  // define a emit() function that mapFct() can call
  global.emit = (_id: DocumentId, value: Document): void => {
    results.push({ _id, value })
  }
  documents.forEach((doc) => mapFct.call(doc))
  return results
}
