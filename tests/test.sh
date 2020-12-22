#!/bin/bash

# Test de bout en bout des commandes "reduce" et "public"
# Source: https://github.com/signaux-faibles/documentation/blob/master/prise-en-main.md#%C3%A9tape-de-calculs-pour-populer-features
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test.output.txt"
GOLDEN_FILE="tests/output-snapshots/test.golden.txt"
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
echo "💎 Computing Features and Public collections..."
echo "- sfdata compact 👉 $(tests/helpers/sfdata-wrapper.sh compact --since-batch=1910)"
echo "- sfdata reduce 👉 $(tests/helpers/sfdata-wrapper.sh reduce --until-batch=1910 --key=012345678)"
echo "- sfdata public 👉 $(tests/helpers/sfdata-wrapper.sh public --until-batch=1910 --key=012345678)"

(tests/helpers/mongodb-container.sh run \
  | tests/helpers/remove-random_order.sh \
  > "${OUTPUT_FILE}" \
) << CONTENTS
  print("// Documents from db.RawData, after call to sfdata compact:");
  printjson(db.RawData.find().toArray());
  print("// Documents from db.Features_debug, after call to sfdata reduce:");
  printjson(db.Features_debug.find().toArray());
  print("// Documents from db.Public_debug, after call to sfdata public:");
  printjson(db.Public_debug.find().toArray());
CONTENTS

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
