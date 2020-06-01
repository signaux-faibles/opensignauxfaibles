// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import * as f from "../common/generatePeriodSerie.js"

const ISODate = (date: string): Date => new Date(date)

const renderedDate = (d: string): string => new Date(d).toString()

const removeRandomOrder = (obj: object): object => {
  Object.keys(obj).forEach(
    (key) =>
      (key === "random_order" && delete obj[key]) ||
      (typeof obj[key] === "object" && removeRandomOrder(obj[key]))
  )
  return obj
}

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
            [renderedDate("Wed Jan 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Feb 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Mar 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Apr 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu May 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Jun 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Jul 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Aug 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Sep 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Oct 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Nov 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Dec 01 2014 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2014-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Jan 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Feb 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Mar 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Apr 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri May 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Jun 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Jul 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Aug 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Sep 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Oct 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Nov 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Dec 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Jan 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Feb 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Mar 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Apr 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun May 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Jun 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Jul 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Aug 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Sep 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Oct 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Nov 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Dec 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Jan 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Feb 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Mar 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Apr 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon May 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Jun 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Jul 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Aug 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Sep 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Oct 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Nov 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Dec 01 2017 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2017-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Jan 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Feb 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Mar 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Apr 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue May 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Jun 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Jul 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed Aug 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Sep 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-09-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Oct 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-10-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Nov 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-11-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Dec 01 2018 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2018-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Tue Jan 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Feb 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-02-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Mar 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-03-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Apr 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-04-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Wed May 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-05-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sat Jun 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-06-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Mon Jul 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-07-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Thu Aug 01 2019 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2019-08-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Sun Sep 01 2019 00:00:00 GMT+0000 (UTC)")]: {
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
  const potentialMapResults = {
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
  }
  t.deepEqual(mapResults, potentialMapResults)

  // 2. reduce
  const reduceKey = importedData.value.key
  const reduceValues = [mapResults[reduceKey]]
  const reduceResults = reduce(reduceKey, reduceValues)
  const potentialReduceResults = {
    batch: {
      1910: {},
    },
    key: "01234567891011",
    scope: "etablissement",
  }
  t.deepEqual(
    reduceResults,
    /*expected[0].value*/ (potentialReduceResults as unknown) as CompanyDataValues // TODO: update types to match data
  )

  // 3. finalize
  const global = globalThis as any
  global.serie_periode = f.generatePeriodSerie(
    ISODate("2014-01-01T00:00:00.000+0000"),
    ISODate("2019-10-01T00:00:00.000+0000")
  )
  const index: ReduceIndexFlags = { algo1: true, algo2: true }
  const finalizeKey = reduceKey
  const finalizeValues = { ...reduceResults, index }
  const finalizeResultValue = finalize(finalizeKey, finalizeValues)
  const finalizeResults = [
    { _id: finalizeKey, value: removeRandomOrder(finalizeResultValue) },
  ]
  t.deepEqual(finalizeResults, expected as unknown)
  // => sample of `actual` VS `expected`:
  //   -             'Tue Oct 01 2019 02:00:00 GMT+0200 (GMT+02:00)': {
  //   -               periode: Date 2019-10-01 00:00:00 UTC {},
  //   -               random_order: 0.19479352943685613,
  //   -               siret: '01234567891011',
  //   -             },
  //   -             'Wed Jan 01 2014 01:00:00 GMT+0100 (GMT+01:00)': {
  //   -               periode: Date 2014-01-01 00:00:00 UTC {},
  //   -               random_order: 0.6133162030905268,
  //   -               siret: '01234567891011',
  //   -             },
  //   +             'Fri Apr 01 2016 00:00:00 GMT+0000 (UTC)': {
  //   +               periode: Date 2016-04-01 00:00:00 UTC {},
  //   +               siret: '01234567891011',
  //   +             },
  //   +             'Fri Aug 01 2014 00:00:00 GMT+0000 (UTC)': {
  //   +               periode: Date 2014-08-01 00:00:00 UTC {},
  //   +               siret: '01234567891011',
  //   +             },
})
