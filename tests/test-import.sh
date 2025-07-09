#!/bin/bash

# Test de bout en bout de la commande "import".
# Référence: https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#vue-densemble-des-canaux-de-transformation-des-donn%C3%A9es
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
export TMP_DIR="tests/tmp-test-execution-files" # for tests/helpers/sfdata-wrapper.sh
OUTPUT_FILE="${TMP_DIR}/test-import.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-import.golden.txt"
mkdir -p "${TMP_DIR}"

# Clear output
> "$OUTPUT_FILE"

# Clean up on exit
function teardown {
    tests/helpers/mongodb-container.sh stop
    tests/helpers/postgres-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start
PORT="5432" tests/helpers/postgres-container.sh start
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
      "dummy":        [],
      "filter":       [],
      "apconso":      [ "/../lib/apconso/testData/apconsoTestData.csv" ],
      "apdemande":    [ "/../lib/apdemande/testData/apdemandeTestData.csv" ],
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
CONTENT

echo -e "\n// Data exported to csv files:\n\n" >> "${OUTPUT_FILE}"

# Loop through all files in the export directory, always in same order
find "./${TMP_DIR}/1910/" -maxdepth 1 -type f | sort | while IFS= read -r file; do
    if [[ -f "$file" ]]; then
        echo "==== $(basename "$file") ====" >> "$OUTPUT_FILE"
        echo -e "\n" >> "$OUTPUT_FILE"
        cat "$file" >> "$OUTPUT_FILE"
        echo -e "\n" >> "$OUTPUT_FILE"
    fi
done

echo -e "\n// Data imported to database:\n\n" >> "${OUTPUT_FILE}"

# List postgresql tables
(tests/helpers/postgres-container.sh run \
  >> "$OUTPUT_FILE" \
) << CONTENT
\echo '==== stg_apconso ===='
SELECT * FROM stg_apconso;
\echo '==== stg_apdemande ===='
SELECT * FROM stg_apdemande;
CONTENT


tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
