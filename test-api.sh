#!/bin/bash

# Test de bout en bout des APIs "reduce" et "public"
# Source: https://github.com/signaux-faibles/documentation/blob/master/prise-en-main.md#%C3%A9tape-de-calculs-pour-populer-features

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
docker stop sf-mongodb

set -e # will stop the script if any command fails with a non-zero exit code

# Clean up on exit
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
trap "{ killall dbmongo; [ -f config.toml ] && rm config.toml; [ -f config.backup.toml ] && mv config.backup.toml config.toml; docker stop sf-mongodb; rm -rf ${DATA_DIR}; echo \"✨ Cleaned up temp directory\"; }" EXIT

echo ""
echo "🐳 Starting MongoDB container..."
docker run \
    --name sf-mongodb \
    --publish 27017:27017 \
    --detach \
    --rm \
    mongo:4

echo ""
echo "🔧 Setting up dbmongo..."
mkdir -p "${DATA_DIR}"
touch "${DATA_DIR}/dummy.csv"
cd dbmongo
go build
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
perl -pi'' -e "s,/foo/bar/data-raw,${DATA_DIR}," config.toml
perl -pi'' -e 's,naf/.*\.csv,dummy.csv,' config.toml

echo ""
echo "📄 Inserting test data..."
sleep 1 # give some time for MongoDB to start
docker exec -i sf-mongodb mongo signauxfaibles << CONTENTS
  db.Admin.remove({})

  db.Admin.insertOne({
    "_id" : {
        "key" : "1910",
        "type" : "batch"
    },
    "files" : {
        "bdf" : [
            "/1910/bdf_1910.csv"
        ]
    },
    "complete_types" : [
    ],
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2019-10-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2019-07-01T00:00:00.000+0000")
    },
    "name" : "Octobre"
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
echo "⚙️ Computing Features and Public collections thru dbmongo API..."
./dbmongo &
sleep 2 # give some time for dbmongo to start
http --ignore-stdin :5000/api/data/compact batch=1910
http --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1910 key=012345678
http --ignore-stdin :5000/api/data/public batch=1910 key=012345678

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
cd ..
docker exec -i sf-mongodb mongo --quiet signauxfaibles > test-api.output.txt << CONTENTS
  print("// Documents from db.RawData, after call to /api/data/compact:");
  db.RawData.find().toArray();
  print("// Documents from db.Features_debug, after call to /api/data/reduce:");
  db.Features_debug.find().toArray();
  print("// Documents from db.Public_debug, after call to /api/data/public:");
  db.Public_debug.find().toArray();
CONTENTS

# exclude random values
grep -v '"random_order" :' test-api.output.txt > test-api.output-documents.txt

echo ""
echo "🆎 Diff between expected and actual output:"
diff test-api.golden-master.txt test-api.output-documents.txt
echo "✅ No diff. The reduce API works as usual."
echo ""
rm test-api.output.txt test-api.output-documents.txt
# Now, the "trap" commands will run, to clean up.
