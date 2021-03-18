#!/bin/bash

# Test de bout en bout de la commande "check".
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-check.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-check.golden.txt"
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
      "admin_urssaf": [
        "/../lib/urssaf/testData/comptesTestData.csv"
      ],
      "debit": [
        "/../lib/urssaf/testData/debitCorrompuTestData.csv"
      ]
    },
    "param" : {
        "date_debut" : ISODate("2001-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2019-02-01T00:00:00.000+0000")
    }
  })
CONTENTS

echo ""
echo "💎 Parsing data..."
RESULT1=$(tests/helpers/sfdata-wrapper.sh check --batch=1910 --parsers='debit')
echo "- sfdata check (with 1 file to parse) 👉 ${RESULT1}"
RESULT2=$(tests/helpers/sfdata-wrapper.sh check --batch=1910)
echo "- sfdata check (with 2 files to parse) 👉 ${RESULT2}"

(tests/helpers/mongodb-container.sh run \
  > "${OUTPUT_FILE}" \
) << CONTENT
print("// Documents from db.Journal:");
printjson(db.Journal.find().sort({ reportType: -1, parserCode: 1 }).toArray().map(doc => ({
  // note: we use map() to force the order of properties at every run of this test
  event: {
    headRejected: doc.event.headRejected,
    headFatal: doc.event.headFatal,
    linesSkipped: doc.event.linesSkipped,
    summary: doc.event.summary,
    batchKey: doc.event.batchKey
  },
  reportType: doc.reportType,
  parserCode: doc.parserCode,
  hasDate: !!doc.date,
  hasStartDate: !!doc.startDate,
})));

print("// Response body from sfdata check, after parsing 1 file:");
CONTENT

echo "${RESULT1}" >> "${OUTPUT_FILE}"

echo "// Response body from sfdata check, after parsing 2 files:" >> "${OUTPUT_FILE}"
echo "${RESULT2}" >> "${OUTPUT_FILE}"

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
