#!/bin/bash

# Test de bout en bout des APIs "reduce" et "public"
# Source: https://github.com/signaux-faibles/documentation/blob/master/prise-en-main.md#%C3%A9tape-de-calculs-pour-populer-features
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
GOLDEN_FILE="tests/output-snapshots/test-api.golden.txt"
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/dbmongo-server.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
    rm -rf "${DATA_DIR}"
    echo "✨ Cleaned up temp directory"
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

MONGODB_PORT="27016" tests/helpers/dbmongo-server.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start

tests/helpers/mongodb-container.sh run > /dev/null << CONTENTS
  db.Admin.remove({})
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

  db.ImportedData.remove({})
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

  db.RawData.remove({})
  db.Features_debug.remove({})
  db.Public_debug.remove({})

CONTENTS

echo ""
echo "💎 Computing Features and Public collections thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/compact 👉 $(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=1910)"
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1910 key=012345678)"
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=1910 key=012345678)"

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
(tests/helpers/mongodb-container.sh run \
  | tests/helpers/remove-random_order.sh \
  > test-api.output.txt \
) << CONTENTS
  print("// Documents from db.RawData, after call to /api/data/compact:");
  db.RawData.find().toArray();
  print("// Documents from db.Features_debug, after call to /api/data/reduce:");
  db.Features_debug.find().toArray();
  print("// Documents from db.Public_debug, after call to /api/data/public:");
  db.Public_debug.find().toArray();
CONTENTS

# Display JS errors logged by MongoDB, if any
sudo docker logs sf-mongodb | grep --color=always "uncaught exception" || true
# TODO: extract to tests/helpers/mongodb-container.sh

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp "test-api.output.txt" "${GOLDEN_FILE}"
else
    # Diff between expected and actual output
    diff --brief "${GOLDEN_FILE}" test-api.output.txt
    echo "✅ No diff. The export worked as expected."
fi
echo ""
rm test-api.output.txt
# Now, the "trap" commands will run, to clean up.
