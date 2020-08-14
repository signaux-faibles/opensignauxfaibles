#!/bin/bash

# Test de bout en bout des APIs "reduce" et "public"
# Source: https://github.com/signaux-faibles/documentation/blob/master/prise-en-main.md#%C3%A9tape-de-calculs-pour-populer-features
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-api.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-api.golden.txt"
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
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2019-10-01T00:00:00.000+0000")
    }
  })

  db.ImportedData.insertOne({
    "_id": "random123abc",
    "value": {
      "batch": {
        "1910": {}
      },
      "scope": "etablissement",
      "index": {
        "algo2": true
      },
      "key": "01234567891011"
    }
  })
CONTENTS

echo ""
echo "💎 Computing Features and Public collections thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/compact 👉 $(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=1910)"
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1910 key=012345678)"
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=1910 key=012345678)"

(tests/helpers/mongodb-container.sh run \
  | tests/helpers/remove-random_order.sh \
  > "${OUTPUT_FILE}" \
) << CONTENTS
  print("// Documents from db.RawData, after call to /api/data/compact:");
  db.RawData.find().toArray();
  print("// Documents from db.Features_debug, after call to /api/data/reduce:");
  db.Features_debug.find().toArray();
  print("// Documents from db.Public_debug, after call to /api/data/public:");
  db.Public_debug.find().toArray();
CONTENTS

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp "${OUTPUT_FILE}" "${GOLDEN_FILE}"
else
    # Diff between expected and actual output
    diff --brief "${GOLDEN_FILE}" "${OUTPUT_FILE}"
    echo "✅ No diff. The export worked as expected."
fi

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
