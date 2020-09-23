#!/bin/bash

# Test de validité des définitions API / Swagger.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

set -e # will stop the script if any command fails with a non-zero exit code

# Convert swagger.yaml to swagger.json
(cd dbmongo/docs && ./convert-yaml-to-json.sh)

# Check if swagger.json has changes that need to be committed
PENDING_CHANGES=$(git diff dbmongo/docs/swagger/swagger.json | wc -l)
if [[ "${PENDING_CHANGES}" -ne "0" ]]
then
    echo "⚠️  swagger.json has ${PENDING_CHANGES} pending changes. please commit them."
    exit 1
fi

echo "✅  swagger.yaml and swagger.json are both valid and in sync."

# Note: If you want to re-generate swagger.yaml based on Go annotations:
# $ go run github.com/swaggo/swag/cmd/swag init
# (see https://github.com/swaggo/swag#how-to-use-it-with-gin)
