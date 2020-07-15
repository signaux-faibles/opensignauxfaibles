#!/bin/bash

# Test de bout en bout de l'API /datapi/exportEntreprise.
#
# Inspiré de test-api.sh.

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
docker stop sf-mongodb

set -e # will stop the script if any command fails with a non-zero exit code

# Clean up on exit
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"
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
cd dbmongo
go build
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
perl -pi'' -e "s,/foo/bar/data-raw,sample-data-raw," config.toml

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
    "files" : {
        "bdf" : [
            "/1910/bdf_1910.csv"
        ]
    },
    "complete_types" : [
    ],
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2014-03-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2014-03-01T00:00:00.000+0000")
    },
    "name" : "TestData"
  })

  db.ImportedData.remove({})
  db.ImportedData.insertMany([
    {
        "_id": "entr1",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "entreprise",
            "index": {
                "algo2": true
            },
            "key": "012345678"
        }
    },
    {
        "_id": "etab1",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "etablissement",
            "index": {
                "algo2": true
            },
            "key": "01234567891011"
        }
    }
  ])

  db.Scores.remove({})
  db.Scores.insertMany([
    {
        "_id": "score1",
        "siret" : "01234567891011",
        "periode" : "2014-01-01",
        "score" : 0.97,
        "batch" : "2002_1",
        "timestamp" : ISODate("2014-01-01T00:00:00.000+0000"),
        "algo" : "algo_avec_urssaf",
        "alert" : "Alerte seuil F1"
    },
    {
        "_id": "score2",
        "siret" : "01234567891011",
        "periode" : "2014-02-01",
        "score" : 0.98,
        "batch" : "2002_1",
        "timestamp" : ISODate("2014-02-01T00:00:00.000+0000"),
        "algo" : "algo_avec_urssaf",
        "alert" : "Alerte seuil F1"
    },
  ])

  db.RawData.remove({})
  db.Features.remove({})
  db.Public.remove({})
  db.Features_debug.remove({})
  db.Public_debug.remove({})
CONTENTS

docker exec -i sf-mongodb mongo signauxfaibles < "${DATA_DIR}/db_popul.js" >/dev/null

echo ""
echo "⚙️ Computing Features and Public collections thru dbmongo API..."
./dbmongo &
sleep 2 # give some time for dbmongo to start
echo "- POST /api/data/compact 👉 $(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2002_1)"
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=2002_1 key=012345678)"

docker exec -i sf-mongodb mongo --quiet signauxfaibles > test-api.output.txt << CONTENTS
  print("// Documents from db.RawData, after call to /api/data/compact:");
  db.RawData.find().toArray();
  db.Public_debug.renameCollection("Public")
  print("// Documents from db.Public, after call to /api/data/public:");
  db.Public.find().toArray();
CONTENTS

echo ""
echo "🚚 Asking API to export enterprise data..."
EXPORT_FILE=$(http --ignore-stdin :5000/datapi/exportEntreprise batch=2002_1 | tr -d '"')
echo "- POST /datapi/exportEntreprise 👉 ${EXPORT_FILE}"

echo ""
# Diff between expected and actual output
cd ..
diff --brief "test-api-entreprise.golden-master.json" "dbmongo/${EXPORT_FILE}"
echo "✅ No diff. The export worked as expected."
echo ""
rm "dbmongo/${EXPORT_FILE}"
# Now, the "trap" commands will run, to clean up.
