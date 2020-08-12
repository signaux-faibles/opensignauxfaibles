#!/bin/bash

# Test de bout en bout de l'API "reduce" à l'aide de données réalistes.
# Inspiré de test-api.sh et finalize_test.js.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.
#
# To update golden files: `$ ./test-api-reduce-2.sh --update`
# 
# These tests require the presence of private files => Make sure to:
# - run `$ git secret reveal` before running these tests;
# - run `$ git secret hide` (to encrypt changes) after updating.

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
sudo docker stop sf-mongodb &>/dev/null

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
GOLDEN_FILE="tests/output-snapshots/reduce-Features.golden.json"
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"

# Clean up on exit
trap "{ killall dbmongo >/dev/null; [ -f config.toml ] && rm config.toml; [ -f config.backup.toml ] && mv config.backup.toml config.toml; sudo docker stop sf-mongodb >/dev/null; rm -rf ${DATA_DIR}; echo \"✨ Cleaned up temp directory\"; }" EXIT

echo ""
echo "🐳 Starting MongoDB container..."
sudo docker run \
    --name sf-mongodb \
    --publish 27016:27017 \
    --detach \
    --rm \
    mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a \
    >/dev/null

echo ""
echo "🔧 Setting up dbmongo..."
cd ./dbmongo
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
perl -pi'' -e "s,/foo/bar/data-raw,sample-data-raw," config.toml
perl -pi'' -e "s,27017,27016," config.toml

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
cat > "${DATA_DIR}/db_popul.js" << CONTENTS
  db.Admin.remove({})
  db.Admin.insertOne({
    "_id" : {
        "key" : "2002_1",
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
cat >> "${DATA_DIR}/db_popul.js" < ../tests/input-data/RawData.sample.json
echo ")" >> "${DATA_DIR}/db_popul.js"

sudo docker exec -i sf-mongodb mongo signauxfaibles > /dev/null < "${DATA_DIR}/db_popul.js"

echo ""
echo "💎 Computing Features and Public collections thru dbmongo API..."
sh -c "./dbmongo &>/dev/null &" # we run in a separate shell to hide the "terminated" message when the process is killed by trap
sleep 2 # give some time for dbmongo to start
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=2002_1)"

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
cd ..
(sudo docker exec -i sf-mongodb mongo --quiet signauxfaibles \
  | tests/helpers/remove-random_order.sh \
  > test-api-2.output.json \
) << CONTENT
  db.Features_TestData.find().toArray();
CONTENT

# Display JS errors logged by MongoDB, if any
sudo docker logs sf-mongodb | grep --color=always "uncaught exception" || true

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp test-api-2.output.json "${GOLDEN_FILE}"
    echo "ℹ️  Updated ${GOLDEN_FILE} => run: $ git secret hide" # to re-encrypt the golden master file, after having updated it
else
    # Diff between expected and actual output
    diff --brief "${GOLDEN_FILE}" test-api-2.output.json
    echo "✅ No diff. The reduce API works as usual."
fi
echo ""
rm test-api-2.output.json
# Now, the "trap" commands will run, to clean up.