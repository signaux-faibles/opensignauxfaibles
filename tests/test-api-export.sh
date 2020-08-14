#!/bin/bash

# Test de bout en bout de GET /api/data/entreprise et /api/data/etablissement.
# Inspiré de test-api.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
COLOR_YELLOW='\033[1;33m'
COLOR_DEFAULT='\033[0m'
ETAB_GOLDEN_FILE="tests/output-snapshots/test-api-export-etablissements.golden.json"
ENTR_GOLDEN_FILE="tests/output-snapshots/test-api-export-entreprises.golden.json"
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p "${DATA_DIR}"

# Clean up on exit
function teardown {
    echo -e "${COLOR_DEFAULT}"
    tests/helpers/dbmongo-server.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
    rm -rf ${DATA_DIR}
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

tests/helpers/mongodb-container.sh run < "${DATA_DIR}/db_popul.js" >/dev/null

echo ""
echo "💎 Computing the Public collection thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/compact 👉 $(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2002_1)"
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=2002_1 key=.........)" # we specify a placeholder value as key, so that PublicOne() is run instead of Public(), so the data is generated for etablissements that don't have effectif values, and therefore are outside of the "algo2" scope.

echo ""
echo "🚚 Asking API to export enterprise data..."
# This step is required only if key was provided when calling POST /api/data/public
RENAME_RESULT=$(tests/helpers/mongodb-container.sh run <<< 'db.Public_debug.renameCollection("Public");')
echo "- rename 'Public_debug' collection to 'Public' 👉 ${RENAME_RESULT}"
# Make sure that the export only relies on Score and Public collections => clear collections that were populated for/by other endpoints
CLEAN_RESULT=$(tests/helpers/mongodb-container.sh run <<< 'db.Admin.drop(); db.ImportedData.drop(); db.RawData.drop();')
echo "- drop other db collections 👉 ${CLEAN_RESULT}"

function stopIfFailed {
    if [[ "$1" == *failed* ]]
    then
        exit 1
    fi
}

# Parameter validation
RESULT=$(http --print=b --ignore-stdin GET :5000/api/data/etablissements key=="invalid" | (grep "key doit être un numéro SIREN" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}"))
echo "- GET /api/data/etablissements with invalid key 👉 ${RESULT}"
stopIfFailed "${RESULT}"
RESULT=$(http --print=b --ignore-stdin GET :5000/api/data/entreprises key=="invalid" | (grep "key doit être un numéro SIREN" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}"))
echo "- GET /api/data/entreprises with invalid key 👉 ${RESULT}"
stopIfFailed "${RESULT}"

# GET /api/data/etablissements with key=212345678 should return just one match
FILE=dbmongo/$(http --print=b --ignore-stdin GET :5000/api/data/etablissements key=="212345678" | tr -d '"')
MATCH=$(grep --quiet "etablissement_21234567891011" "${FILE}" && echo "found etablissement_21234567891011" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}")
COUNT=$(wc -l <"${FILE}")
rm "${FILE}"
echo "- GET /api/data/etablissements with key=212345678 👉 ${MATCH}, ${COUNT} result(s)"
stopIfFailed "${MATCH}"
if [[ "${COUNT}" -ne "1" ]]
then
    exit 1
fi

# GET /api/data/entreprises with key=212345678 should return just one match
FILE=dbmongo/$(http --print=b --ignore-stdin GET :5000/api/data/entreprises key=="212345678" | tr -d '"')
MATCH=$(grep --quiet "entreprise_212345678" "${FILE}" && echo "found entreprise_212345678" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}")
COUNT=$(wc -l <"${FILE}")
rm "${FILE}"
echo "- GET /api/data/entreprises with key=212345678 👉 ${MATCH}, ${COUNT} result(s)"
stopIfFailed "${MATCH}"
if [[ "${COUNT}" -ne "1" ]]
then
    exit 1
fi

# Export enterprise data
ETABLISSEMENTS_FILE=dbmongo/$(http --print=b --ignore-stdin GET :5000/api/data/etablissements | tr -d '"')
echo "- GET /api/data/etablissements 👉 ${ETABLISSEMENTS_FILE}"
ENTREPRISES_FILE=dbmongo/$(http --print=b --ignore-stdin GET :5000/api/data/entreprises | tr -d '"')
echo "- GET /api/data/entreprises 👉 ${ENTREPRISES_FILE}"

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file using ${ETABLISSEMENTS_FILE}..."
    cp "${ETABLISSEMENTS_FILE}" "${ETAB_GOLDEN_FILE}"
    echo "🖼  Updating golden master file using ${ENTREPRISES_FILE}..."
    cp "${ENTREPRISES_FILE}" "${ENTR_GOLDEN_FILE}"
else
    # Diff between expected and actual output
    echo -e "${COLOR_YELLOW}"
    diff --brief "${ETAB_GOLDEN_FILE}" "${ETABLISSEMENTS_FILE}" # will stop the script if files are different
    diff --brief "${ENTR_GOLDEN_FILE}" "${ENTREPRISES_FILE}" # will stop the script if files are different
    echo -e "${COLOR_DEFAULT}"
    echo "✅ No diff. The export worked as expected."
fi
echo ""
rm "${ETABLISSEMENTS_FILE}" "${ENTREPRISES_FILE}"
# Now, the "trap" commands will run, to clean up.
