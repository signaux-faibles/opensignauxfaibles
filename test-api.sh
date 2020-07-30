#!/bin/bash

# Test de bout en bout des APIs "reduce" et "public"
# Source: https://github.com/signaux-faibles/documentation/blob/master/prise-en-main.md#%C3%A9tape-de-calculs-pour-populer-features

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
sudo docker stop sf-mongodb &>/dev/null

set -e # will stop the script if any command fails with a non-zero exit code

# Clean up on exit
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"
trap "{ killall dbmongo >/dev/null; [ -f config.toml ] && rm config.toml; [ -f config.backup.toml ] && mv config.backup.toml config.toml; sudo docker stop sf-mongodb >/dev/null; rm -rf ${DATA_DIR}; echo \"✨ Cleaned up temp directory\"; }" EXIT

echo ""
echo "🐳 Starting MongoDB container..."
sudo docker run \
    --name sf-mongodb \
    --publish 27016:27017 \
    --detach \
    --rm \
    mongo:4 \
    >/dev/null

echo ""
echo "🔧 Setting up dbmongo..."
cd dbmongo
go build
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
perl -pi'' -e "s,/foo/bar/data-raw,sample-data-raw," config.toml
perl -pi'' -e "s,27017,27016," config.toml

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
sudo docker exec -i sf-mongodb mongo signauxfaibles  > /dev/null << CONTENTS
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
sh -c "./dbmongo &>/dev/null &" # we run in a separate shell to hide the "terminated" message when the process is killed by trap
sleep 2 # give some time for dbmongo to start
echo "- POST /api/data/compact 👉 $(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=1910)"
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1910 key=012345678)"
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=1910 key=012345678)"

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
cd ..
sudo docker exec -i sf-mongodb mongo --quiet signauxfaibles > test-api.output.txt << CONTENTS
  print("// Documents from db.RawData, after call to /api/data/compact:");
  db.RawData.find().toArray();
  print("// Documents from db.Features_debug, after call to /api/data/reduce:");
  db.Features_debug.find().toArray();
  print("// Documents from db.Public_debug, after call to /api/data/public:");
  db.Public_debug.find().toArray();
CONTENTS

# Display JS errors logged by MongoDB, if any
sudo docker logs sf-mongodb | grep --color=always "uncaught exception" || true

# exclude random values
grep -v '"random_order" :' test-api.output.txt > test-api.output-documents.txt

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp "test-api.output-documents.txt" "test-api.golden-master.txt"
else
    # Diff between expected and actual output
    diff --brief test-api.golden-master.txt test-api.output-documents.txt
    echo "✅ No diff. The export worked as expected."
fi
echo ""
rm test-api.output.txt test-api.output-documents.txt
# Now, the "trap" commands will run, to clean up.
