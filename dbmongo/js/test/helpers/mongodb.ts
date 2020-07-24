declare let emit: (_id: unknown, value: unknown) => void

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

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
    const key = JSON.stringify(_id) // as { siren: string; batch: string; periode: Date }
    acc[key] = acc[key] || []
    acc[key].push({ key: _id, value: value })
    return acc
  }, {} as Indexed<K, V>)
