#!/bin/bash

# Test de bout en bout de l'API "reduce" à l'aide de données publiques.
# Inspiré de test-api-reduce-2.sh et algo2_tests.ts.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
GOLDEN_FILE="tests/output-snapshots/test-api-reduce.golden.json"
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

echo ""
echo "🐳 Starting MongoDB container..."
PORT="27016" tests/helpers/mongodb-container.sh start

echo ""
echo "🔧 Setting up dbmongo..."
MONGODB_PORT="27016" tests/helpers/dbmongo-server.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
cat > "${DATA_DIR}/db_popul.js" << CONTENTS
  db.Admin.remove({})
  db.Admin.insertOne({
    "_id" : {
        "key" : "1905",
        "type" : "batch"
    },
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2016-01-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2016-03-01T00:00:00.000+0000")
    },
    "name" : "TestData"
  })

  db.Features_TestData.remove({})

  db.RawData.remove({})
  db.RawData.insertMany(
CONTENTS
node -e "console.log(require('./dbmongo/js/test/data/objects.js').makeObjects.toString().replace('ISODate => ([', '[').replace('])', ']'))" \
  >> "${DATA_DIR}/db_popul.js"
echo ")" >> "${DATA_DIR}/db_popul.js"

tests/helpers/mongodb-container.sh run < "${DATA_DIR}/db_popul.js" >/dev/null

echo ""
echo "💎 Computing the Features collection thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1905)"

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
(tests/helpers/mongodb-container.sh run \
  > "test-api-reduce.output-documents.json" \
) << CONTENT
  db.Features_TestData.find().toArray();
CONTENT

# Display JS errors logged by MongoDB, if any
sudo docker logs sf-mongodb | grep --color=always "uncaught exception" || true
# TODO: extract to tests/helpers/mongodb-container.sh

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp "test-api-reduce.output-documents.json" "${GOLDEN_FILE}"
else
    # Diff between expected and actual output
    diff --brief "${GOLDEN_FILE}" "test-api-reduce.output-documents.json"
    echo "✅ No diff. The reduce API works as usual."
fi
echo ""
rm "test-api-reduce.output-documents.json"
# Now, the "trap" commands will run, to clean up.
