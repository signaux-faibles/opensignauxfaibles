#!/bin/bash

# Test de bout en bout de GET /api/data/entreprise et /api/data/etablissement.
# Inspiré de test-api.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
sudo docker stop sf-mongodb &>/dev/null

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
COLOR_YELLOW='\033[1;33m'
COLOR_DEFAULT='\033[0m'
ETAB_GOLDEN_FILE="tests/output-snapshots/test-api-export-etablissements.golden.json"
ENTR_GOLDEN_FILE="tests/output-snapshots/test-api-export-entreprises.golden.json"
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"

# Clean up on exit
trap "{ echo -e \"${COLOR_DEFAULT}\"; killall dbmongo >/dev/null; [ -f config.toml ] && rm config.toml; [ -f config.backup.toml ] && mv config.backup.toml config.toml; sudo docker stop sf-mongodb >/dev/null; rm -rf ${DATA_DIR}; echo \"✨ Cleaned up temp directory\"; }" EXIT

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
cd dbmongo
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
        "date_fin" : ISODate("2014-03-01T00:00:00.000+0000")
    }
  })

  db.ImportedData.remove({})
  // The random order of documents is intentional, to make sure that the output is correctly sorted no matter what
  db.ImportedData.insertMany([
    {
        "_id": "etab2",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "etablissement",
            "key": "01234567891012"
        }
    },
    {
        "_id": "etab21",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "etablissement",
            "key": "21234567891011"
        }
    },
    {
        "_id": "entr2",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "entreprise",
            "key": "212345678"
        }
    },
    {
        "_id": "entr1",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "entreprise",
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
            "key": "01234567891011"
        }
    },
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
  db.Public.remove({})
  db.Public_debug.remove({})
CONTENTS

sudo docker exec -i sf-mongodb mongo signauxfaibles < "${DATA_DIR}/db_popul.js" >/dev/null

echo ""
echo "💎 Computing Features and Public collections thru dbmongo API..."
sh -c "./dbmongo &>/dev/null &" # we run in a separate shell to hide the "terminated" message when the process is killed by trap
sleep 2 # give some time for dbmongo to start
echo "- POST /api/data/compact 👉 $(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2002_1)"
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=2002_1 key=.........)"

echo ""
echo "🚚 Asking API to export enterprise data..."
# This step is required only if key was provided when calling POST /api/data/public
RENAME_RESULT=$(echo 'db.Public_debug.renameCollection("Public");' | sudo docker exec -i sf-mongodb mongo --quiet signauxfaibles)
echo "- rename 'Public_debug' collection to 'Public' 👉 ${RENAME_RESULT}"
# Make sure that the export only relies on Score and Public collections => clear collections that were populated for/by other endpoints
CLEAN_RESULT=$(echo 'db.Admin.drop(); db.ImportedData.drop(); db.RawData.drop();' | sudo docker exec -i sf-mongodb mongo --quiet signauxfaibles)
echo "- drop other db collections 👉 ${CLEAN_RESULT}"
# Export enterprise data
ETABLISSEMENTS_FILE=$(http --print=b --ignore-stdin GET :5000/api/data/etablissements | tr -d '"')
echo "- GET /api/data/etablissements 👉 ${ETABLISSEMENTS_FILE}"
ENTREPRISES_FILE=$(http --print=b --ignore-stdin GET :5000/api/data/entreprises | tr -d '"')
echo "- GET /api/data/entreprises 👉 ${ENTREPRISES_FILE}"

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file using ${ETABLISSEMENTS_FILE}..."
    cp "${ETABLISSEMENTS_FILE}" "../${ETAB_GOLDEN_FILE}"
    echo "🖼  Updating golden master file using ${ENTREPRISES_FILE}..."
    cp "${ENTREPRISES_FILE}" "../${ENTR_GOLDEN_FILE}"
else
    # Diff between expected and actual output
    echo -e "${COLOR_YELLOW}"
    diff --brief "../${ETAB_GOLDEN_FILE}" "${ETABLISSEMENTS_FILE}" # will stop the script if files are different
    diff --brief "../${ENTR_GOLDEN_FILE}" "${ENTREPRISES_FILE}" # will stop the script if files are different
    echo -e "${COLOR_DEFAULT}"
    echo "✅ No diff. The export worked as expected."
fi
echo ""
rm "${ETABLISSEMENTS_FILE}" "${ENTREPRISES_FILE}"
# Now, the "trap" commands will run, to clean up.
