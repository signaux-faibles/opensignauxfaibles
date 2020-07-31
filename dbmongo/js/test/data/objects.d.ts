import "../../globals"
export type TestDataItem = { _id: string; value: CompanyDataValuesWithFlags }

export const objects: TestDataItem[]
export const makeObjects: (ISODate: (string) => Date) => TestDataItem[]
