#!/bin/bash

# Test de bout en bout de l'API "reduce"
# Source: https://github.com/signaux-faibles/documentation/blob/master/prise-en-main.md#%C3%A9tape-de-calculs-pour-populer-features

# Interrompre le conteneur Docker d'une exécution précédente de ce test, si besoin
docker stop sf-mongodb

set -e # will stop the script if any command fails with a non-zero exit code

# Clean up on exit
trap "{ mv config.backup.toml config.toml; docker stop sf-mongodb; rm -rf ${DATA_DIR}; echo \"Cleaned up temp directory\"; }" EXIT

# 1. Lancement de mongodb avec Docker
docker run \
    --name sf-mongodb \
    --publish 27017:27017 \
    --detach \
    --rm \
    mongo:4

# 2. Préparation du répertoire de données
DATA_DIR=$(pwd)/tmp-opensignauxfaibles-data-raw
mkdir -p ${DATA_DIR}
touch ${DATA_DIR}/dummy.csv

# 3. Installation et configuration de dbmongo
cd dbmongo
go build
mv config.toml config.backup.toml
cp config-sample.toml config.toml
sed -i '' "s,/foo/bar/data-raw,${DATA_DIR}," config.toml
sed -i '' 's,naf/.*\.csv,dummy.csv,' config.toml

# 4. Ajout de données de test
docker exec -i sf-mongo mongo signauxfaibles << CONTENTS
  db.createCollection('RawData')

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

  db.RawData.remove({})

  db.RawData.insertOne({
    "_id": "01234567891011",
    "value": {
        "scope": "etablissement",
        "index": {
        "algo2": true
        }
    }
  })
CONTENTS

# 5. Exécution des calculs pour populer la collection "Features"
./dbmongo &
DBMONGO_PID=$!
sleep 5 
http :5000/api/data/reduce algo=algo2 batch=1910 key=012345678
kill ${DBMONGO_PID}
echo "db.Features_debug.find()" \
  | docker exec -i sf-mongo mongo signauxfaibles
