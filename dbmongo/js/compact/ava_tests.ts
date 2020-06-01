import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { reduce } from "./reduce"

const ISODate = (date: string): Date => new Date(date)

// input data from test-api.sh
const importedData = {
  _id: "random123abc",
  value: {
    batch: {
      "1910": {},
    },
    scope: "etablissement",
    index: {
      algo2: true,
    },
    key: "01234567891011",
  },
}

// output data from test-api.sh
const expected = [
  {
    _id: "01234567891011",
    value: {
      batch: {
        "1910": {
          reporder: {
            "Wed Jan 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Feb 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Mar 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Apr 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu May 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Jun 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Jul 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Aug 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Sep 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Oct 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Nov 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Dec 01 2014 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2014-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Jan 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Feb 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Mar 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Apr 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri May 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Jun 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Jul 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Aug 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Sep 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Oct 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Nov 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Dec 01 2015 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2015-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Jan 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Feb 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Mar 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Apr 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun May 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Jun 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Jul 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Aug 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Sep 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Oct 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Nov 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Dec 01 2016 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2016-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Jan 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Feb 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Mar 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Apr 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon May 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Jun 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Jul 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Aug 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Sep 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Oct 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Nov 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Dec 01 2017 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2017-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Jan 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Feb 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Mar 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Apr 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue May 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Jun 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Jul 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed Aug 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Sep 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Oct 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Nov 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Dec 01 2018 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2018-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Tue Jan 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Feb 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Fri Mar 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Apr 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Wed May 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sat Jun 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Mon Jul 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Thu Aug 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            "Sun Sep 01 2019 00:00:00 GMT+0000 (UTC)": {
              periode: ISODate("2019-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
          },
        },
      },
      scope: "etablissement",
      index: {
        algo1: false,
        algo2: false,
      },
      key: "01234567891011",
    },
  },
]

const runMongoMap = (mapFct: () => void, keyVal: object): object => {
  const results = {}
  globalThis.emit = (key: string, value: any): void => {
    results[key] = value
  }
  mapFct.call(keyVal)
  return results
}

test(`exécution complète de la chaine "compact"`, (t: ExecutionContext) => {
  // 1. map
  const mapResults = runMongoMap(map, importedData)
  t.deepEqual(mapResults, {
    "01234567891011": {
      batch: {
        1910: {},
      },
      index: {
        algo2: true,
      },
      key: "01234567891011",
      scope: "etablissement",
    },
  })

  // 2. reduce
  const key = "01234567891011"
  const values = [mapResults[key]]
  const reduceResults = reduce(key, values)
  t.deepEqual(
    reduceResults,
    (expected[0].value as unknown) as CompanyDataValues // TODO: update types to match data
  )
  // => some fields of final expected results are still missing:
  //
  //   {
  //     batch: {
  //       1910: {
  // +       reporder: Object { … },
  //       },
  //     },
  // +   index: {
  // +     algo1: false,
  // +     algo2: false,
  // +   },
  //     key: '01234567891011',
  //     scope: 'etablissement',
  //   }
})
