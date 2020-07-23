// This module will be overwritten by test_algo2.sh, to feed private test data into test_algo2.ts

export const makeTestData = ({
  ISODate,
  NumberInt,
}: {
  ISODate: (date: string) => Date
  NumberInt: (i: number) => number
}): unknown[] => [
  {
    dummy1: ISODate(""),
    dummy2: NumberInt(1),
  },
]
