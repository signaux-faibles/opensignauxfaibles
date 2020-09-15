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
      "admin_urssaf": [
        "/../lib/urssaf/testData/comptesTestData.csv"
      ],
      "delai": [
        "/../lib/urssaf/testData/delaiTestData.csv"
      ]
    },
    "param" : {
        "date_debut" : ISODate("2001-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2019-02-01T00:00:00.000+0000")
    }
  })
CONTENTS

echo ""
echo "💎 Parsing and importing data thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/import 👉 $(http --print=b --ignore-stdin :5000/api/data/import batch=1910 parsers:='["delai"]')"

sleep 1 # give some time for dbmongo to parse and import data

(tests/helpers/mongodb-container.sh run \
  | tests/helpers/remove-random_order.sh \
  > "${OUTPUT_FILE}" \
) <<< 'printjson(db.ImportedData.find().toArray());'

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
