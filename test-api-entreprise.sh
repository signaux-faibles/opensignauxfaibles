#!/bin/bash

# Test de bout en bout de l'API /datapi/exportEntreprise.
#
# Inspiré de test-api.sh.
#
# Ce test requiert l'accès à un serveur privé, et n'est donc pas inclus dans la
# suite de tests exécutée en Integration Continue.

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
cd dbmongo
go build
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
perl -pi'' -e "s,/foo/bar/data-raw,sample-data-raw," config.toml

echo ""
echo "📄 Inserting test data..."
sleep 1 # give some time for MongoDB to start
cat > "${DATA_DIR}/db_popul.js" << CONTENTS
  db.Scores.remove({})
  db.Scores.insertMany([
    {
        "siret" : "01234567891011",
        "periode" : "2014-01-01",
        "score" : 0.97,
        "batch" : "2002_1",
        "timestamp" : ISODate("2014-01-01T00:00:00.000+0000"),
        "algo" : "algo_avec_urssaf",
        "alert" : "Alerte seuil F1"
    },
    {
        "siret" : "01234567891011",
        "periode" : "2015-01-01",
        "score" : 0.98,
        "batch" : "2002_1",
        "timestamp" : ISODate("2015-01-01T00:00:00.000+0000"),
        "algo" : "algo_avec_urssaf",
        "alert" : "Alerte seuil F1"
    },
  ])

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
        "date_fin_effectif" : ISODate("2016-03-01T00:00:00.000+0000")
    },
    "name" : "TestData"
  })

  db.RawData.remove({})
  db.RawData.insertMany(
CONTENTS
cat >> "${DATA_DIR}/db_popul.js" < "${DATA_DIR}/reduce_test_data.json"
echo ")" >> "${DATA_DIR}/db_popul.js"

docker exec -i sf-mongodb mongo signauxfaibles < "${DATA_DIR}/db_popul.js"

echo ""
echo "⚙️ Computing Features and Public collections thru dbmongo API..."
./dbmongo &
sleep 2 # give some time for dbmongo to start
http --ignore-stdin :5000/api/data/public batch=2002_1
EXPORT_FILE=$(http --ignore-stdin :5000/datapi/exportEntreprise batch=2002_1 | tr -d '"')

echo ""
echo "🆎 Diff between expected and actual output:"
# diff "${DATA_DIR}/test-api-entreprise_golden.json" "${EXPORT_FILE}"
diff "entreprise_golden.json" "${EXPORT_FILE}" # (diff provisoire)
echo "✅ No diff. The reduce API works as usual."
echo ""
rm "${EXPORT_FILE}"
# Now, the "trap" commands will run, to clean up.
