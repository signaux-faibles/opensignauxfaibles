#!/bin/bash

# Test de bout en bout de l'API "reduce" à l'aide de données publiques.
# Inspiré de test-api-reduce-2.sh et algo2_tests.ts.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
sudo docker stop sf-mongodb &>/dev/null

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
GOLDEN_FILE="tests/output-snapshots/test-api-reduce.golden.json"
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
go build
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
node -e "console.log(require('./js/test/data/objects.js').makeObjects.toString().replace('ISODate => ([', '[').replace('])', ']'))" \
  >> "${DATA_DIR}/db_popul.js"
echo ")" >> "${DATA_DIR}/db_popul.js"

sudo docker exec -i sf-mongodb mongo signauxfaibles > /dev/null < "${DATA_DIR}/db_popul.js"

echo ""
echo "💎 Computing Features and Public collections thru dbmongo API..."
sh -c "./dbmongo &>/dev/null &" # we run in a separate shell to hide the "terminated" message when the process is killed by trap
sleep 2 # give some time for dbmongo to start
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1905)"

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
cd ..
echo "db.Features_TestData.find().toArray();" \
  | sudo docker exec -i sf-mongodb mongo --quiet signauxfaibles \
  > "test-api-reduce.output-documents.json"

# Display JS errors logged by MongoDB, if any
sudo docker logs sf-mongodb | grep --color=always "uncaught exception" || true

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
