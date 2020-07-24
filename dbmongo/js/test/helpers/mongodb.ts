const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

type DocumentId = unknown // can be an ObjectID, or other
type Document = Record<string, unknown>
type MapResult = { _id: DocumentId; value: Document }

export const runMongoMap = <Doc extends Document>(
  mapFct: (this: Doc) => void, // will call global emit()
  document: Doc
): MapResult[] => {
  const results: MapResult[] = []
  global.emit = (_id: DocumentId, value: Document): void => {
    results.push({ _id, value })
  }
  mapFct.call(document)
  return results
}
