#!/bin/bash

# Test de bout en bout de l'API "import".
# Référence: https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#vue-densemble-des-canaux-de-transformation-des-donn%C3%A9es
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-api-import.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-api-import.golden.txt"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/dbmongo-server.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

MONGODB_PORT="27016" tests/helpers/dbmongo-server.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start

tests/helpers/mongodb-container.sh run > /dev/null << CONTENTS
  db.Admin.insertOne({
    "_id" : {
        "key" : "1910",
        "type" : "batch"
    },
    "files": {
      "apconso":      [ "/../lib/apconso/testData/apconsoTestData.csv" ],
      "apdemande":    [ "/../lib/apdemande/testData/apdemandeTestData.csv" ],
      "bdf":          [ "/../lib/bdf/testData/bdfTestData.csv" ],
      "diane":        [ "/../lib/diane/testData/dianeTestData.txt" ],
      "ellisphere":   [ "/../lib/ellisphere/testData/ellisphereTestData.excel" ],
      "sirene":       [ "/../lib/sirene/testData/sireneTestData.csv" ],
      "sirene_ul":    [ "/../lib/sirene_ul/testData/sireneULTestData.csv" ],
      "admin_urssaf": [ "/../lib/urssaf/testData/comptesTestData.csv" ],
      "debit":        [ "/../lib/urssaf/testData/debitTestData.csv" ],
      "ccsf":         [ "/../lib/urssaf/testData/ccsfTestData.csv" ],
      "cotisation":   [ "/../lib/urssaf/testData/cotisationTestData.csv" ],
      "delai":        [ "/../lib/urssaf/testData/delaiTestData.csv" ],
//      "effectif":     [ "/../lib/urssaf/testData/effectifTestData.csv" ],
//      "effectif_ent": [ "/../lib/urssaf/testData/effectifEntTestData.csv" ],
//      "procol":       [ "/../lib/urssaf/testData/procolTestData.csv" ],
    },
    "param" : {
        "date_debut" : ISODate("2019-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2019-02-01T00:00:00.000+0000")
    }
  })
CONTENTS

echo ""
echo "💎 Parsing and importing data thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/import 👉 $(http --print=b --ignore-stdin :5000/api/data/import batch=1910 noFilter:=true)"

OUTPUT_GZ_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=ImportedData | tr -d '"')
echo "- POST /api/data/validate 👉 ${OUTPUT_GZ_FILE}"

(tests/helpers/mongodb-container.sh run \
  | perl -p -e 's/"[0-9a-z]{32}"/"______________Hash______________"/' \
  | perl -p -e 's/"[0-9a-z]{24}"/"________ObjectId________"/' \
  | perl -p -e 's/"periode" : ISODate\("....-..-..T..:..:..Z"\)/"periode" : ISODate\("_______ Date _______"\)/' \
  > "${OUTPUT_FILE}" \
) << CONTENT
print("// Documents from db.ImportedData, after call to /api/data/import:");
printjson(db.ImportedData.find().sort({"value.key":1}).toArray().map(doc => ({
  ...doc,
  value: {
    ...doc.value,
    // on classe les données par type, de manière à ce que l'ordre soit stable
    batch: Object.keys(doc.value.batch).reduce((batch, batchKey) => ({
      ...batch,
      [ batchKey ]: Object.keys(doc.value.batch[batchKey]).sort().reduce((batchData, dataType) => ({
        ...batchData,
        [ dataType ]: Object.keys(doc.value.batch[batchKey][dataType]).sort().reduce((hashedData, hash) => ({
          ...hashedData,
          [ hash ]: doc.value.batch[batchKey][dataType][hash]
        }), {})
      }), {})
    }), {})
  }
})));

print("// Reports from db.Journal:");
// on classe les données par type, de manière à ce que l'ordre soit stable
printjson(db.Journal.find({ "event.report": { "\$exists": true } }).sort({ code: 1 }).toArray().map(doc => ({
  event: {
    headFilters: doc.event.headFilters,
    headErrors: doc.event.headErrors,
    headFatal: doc.event.headFatal,
    report: doc.event.report,
    batchKey: doc.event.batchKey
  },
  code: doc.code
})));

print("// Critical errors from db.Journal:");
printjson(db.Journal.find({ priority: "critical" }).sort({ code: 1 }).toArray());

print("// Results of call to /api/data/validate:");
CONTENT

zcat < "${OUTPUT_GZ_FILE}" \
  | perl -p -e 's/"[0-9a-z]{32}"/"______________Hash______________"/' \
  | perl -p -e 's/"[0-9a-z]{24}"/"________ObjectId________"/' \
  | perl -p -e 's/"periode" : ISODate\("....-..-..T..:..:..Z"\)/"periode" : ISODate\("_______ Date _______"\)/' \
  | sort \
  >> "${OUTPUT_FILE}"

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm "${OUTPUT_GZ_FILE}"
rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
