#!/bin/bash

# Test de bout en bout de la commande "import".
# Référence: https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#vue-densemble-des-canaux-de-transformation-des-donn%C3%A9es
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-import.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-import.golden.txt"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start
export MONGODB_PORT="27016" # for tests/helpers/sfdata-wrapper.sh

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
      "paydex":       [ "/../lib/paydex/testData/paydexTestData.csv" ],
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
      "effectif":     [ "/../lib/urssaf/testData/effectifTestData.csv" ],
      "effectif_ent": [ "/../lib/urssaf/testData/effectifEntTestData.csv" ],
      "procol":       [ "/../lib/urssaf/testData/procolTestData.csv" ],
    },
    "param" : {
        "date_debut" : ISODate("2019-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2019-02-01T00:00:00.000+0000")
    }
  })
CONTENTS

echo ""
echo "💎 Parsing and importing data..."
echo "- sfdata import 👉 $(tests/helpers/sfdata-wrapper.sh import --batch=1910 --no-filter)"

VALIDATION_REPORT=$(tests/helpers/sfdata-wrapper.sh validate --collection=ImportedData)
echo "- sfdata validate"

(tests/helpers/mongodb-container.sh run \
  | perl -p -e 's/"[0-9a-z]{24}"/"________ObjectId________"/' \
  | perl -p -e 's/"periode" : ISODate\("....-..-..T..:..:..Z"\)/"periode" : ISODate\("_______ Date _______"\)/' \
  > "${OUTPUT_FILE}" \
) << CONTENT
print("// Documents from db.ImportedData, after call to sfdata import:");
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
CONTENT

echo "- sfdata purgeNotCompacted 👉 $(tests/helpers/sfdata-wrapper.sh purgeNotCompacted --i-understand-what-im-doing)"

(tests/helpers/mongodb-container.sh run \
  >> "${OUTPUT_FILE}" \
) << CONTENT
print("// Reports from db.Journal:");
// on classe les données par type, de manière à ce que l'ordre soit stable
printjson(db.Journal.find().sort({ reportType: -1, parserCode: 1 }).toArray().map(doc => (doc.event ? {
  event: {
    headRejected: doc.event.headRejected,
    headFatal: doc.event.headFatal,
    linesSkipped: doc.event.linesSkipped,
    summary: doc.event.summary,
    batchKey: doc.event.batchKey
  },
  reportType: doc.reportType,
  parserCode: doc.parserCode,
  hasCommitHash: !!doc.commitHash,
  hasDate: !!doc.date,
  hasStartDate: !!doc.startDate,
} : {
  reportType: doc.reportType,
  hasCommitHash: !!doc.commitHash,
  hasDate: !!doc.date,
  hasStartDate: !!doc.startDate,
})));

print("// Results of call to sfdata validate:");
CONTENT

echo "${VALIDATION_REPORT}" \
  | perl -p -e 's/"[0-9a-z]{24}"/"________ObjectId________"/' \
  | perl -p -e 's/"periode" : ISODate\("....-..-..T..:..:..Z"\)/"periode" : ISODate\("_______ Date _______"\)/' \
  | sort \
  | while read -r line; do ( echo "- listing validation errors..." 1>&2; echo "$line"; cd "js"; npx -p ajv-cli -p ajv-bsontype ajv -c "ajv-bsontype" -s "../validation/bdf.schema.json" -d <(echo "$line") 2>&1 || true ); done \
  >> "${OUTPUT_FILE}"

# Print test results from stdin. Fails on any "false" result.
# Expected format for each line: "<test label> : <true|false>"
function reportFailedTests {
  while IFS='$\n' read -r line; do
    echo "  - $line" | (grep --color=always " : false") || true # display failed test
    echo "  - $line" | grep " : true" # display passing test, and make the test function fail otherwise
  done
}

(tests/helpers/mongodb-container.sh run \
  | reportFailedTests \
) << CONTENT
  Object.entries({
    "ImportedData was emptied by purgeNotCompacted": db.ImportedData.count() === 0,
  }).forEach(([ testName, testRes ]) => print(testName, ':', testRes));
CONTENT

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
