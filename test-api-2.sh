#!/bin/bash

# Test de bout en bout de l'API "reduce" à l'aide de données réalistes

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
docker stop sf-mongodb

set -e # will stop the script if any command fails with a non-zero exit code

# Clean up on exit
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"
trap "{ killall dbmongo; [ -f config.toml ] && rm config.toml; [ -f config.backup.toml ] && mv config.backup.toml config.toml; docker stop sf-mongodb; rm -rf ${DATA_DIR}; echo \"✨ Cleaned up temp directory\"; }" EXIT

echo ""
echo "🚚 Downloading realistic data set..."
scp "stockage:/home/centos/opensignauxfaibles_tests/*" "${DATA_DIR}/"

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
touch "${DATA_DIR}/dummy.csv"
cd dbmongo
go build
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' "s,/foo/bar/data-raw,${DATA_DIR}," config.toml
  sed -i '' 's,naf/.*\.csv,dummy.csv,' config.toml
else
  sed -i "s,/foo/bar/data-raw,${DATA_DIR}," config.toml
  sed -i 's,naf/.*\.csv,dummy.csv,' config.toml
fi

echo ""
echo "📄 Inserting test data..."
sleep 1 # give some time for MongoDB to start
cat > "${DATA_DIR}/db_popul.js" << CONTENTS
  db.Admin.remove({})
  db.Admin.insertOne({
    "_id" : {
        "key" : "2002_1",
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
        "date_fin" : ISODate("2016-01-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2016-01-01T00:00:00.000+0000")
    },
    "name" : "TestData"
  })

  db.RawData.remove({})
  db.RawData.insertMany(
CONTENTS
cat >> "${DATA_DIR}/db_popul.js" < "${DATA_DIR}/reduce_test_data.json"
cat >> "${DATA_DIR}/db_popul.js" << CONTENTS
  )
CONTENTS

docker exec -i sf-mongodb mongo signauxfaibles < "${DATA_DIR}/db_popul.js"

echo ""
echo "⚙️ Computing Features and Public collections thru dbmongo API..."
./dbmongo &
sleep 2 # give some time for dbmongo to start
http --ignore-stdin :5000/api/data/reduce algo=algo2 batch=2002_1

echo ""
echo "🕵️‍♀️ Checking resulting Features..."
cd ..
echo "db.Features_TestData.find().toArray();" \
  | docker exec -i sf-mongodb mongo --quiet signauxfaibles \
  | grep -v '"random_order" :' \
  | npx prettier --stdin-filepath test-api-2.output.json \
  > test-api-2.output.json

echo ""
echo "🆎 Diff between expected and actual output:"
diff "${DATA_DIR}/finalize_golden.log" test-api-2.output.json
echo "✅ No diff. The reduce API works as usual."
echo ""
rm test-api-2.output.json
# Now, the "trap" commands will run, to clean up.
